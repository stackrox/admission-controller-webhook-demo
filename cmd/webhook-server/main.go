/*
Copyright (c) 2019 StackRox Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"k8s.io/api/admission/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	tlsDir      = `/run/secrets/tls`
	tlsCertFile = `tls.crt`
	tlsKeyFile  = `tls.key`
)

var (
	deployResource = metav1.GroupVersionResource{Version: "v1", Group: "apps", Resource: "deployments"}
)

// applySecurityDefaults implements the logic of our example admission controller webhook. For every pod that is created
// (outside of Kubernetes namespaces), it first checks if `runAsNonRoot` is set. If it is not, it is set to a default
// value of `false`. Furthermore, if `runAsUser` is not set (and `runAsNonRoot` was not initially set), it defaults
// `runAsUser` to a value of 1234.
//
// To demonstrate how requests can be rejected, this webhook further validates that the `runAsNonRoot` setting does
// not conflict with the `runAsUser` setting - i.e., if the former is set to `true`, the latter must not be `0`.
// Note that we combine both the setting of defaults and the check for potential conflicts in one webhook; ideally,
// the latter would be performed in a validating webhook admission controller.
func applySecurityDefaults(req *v1beta1.AdmissionRequest) ([]patchOperation, error) {
	log.Printf("applying changes")

	// This handler should only get called on deploy objects as per the MutatingWebhookConfiguration in the YAML file.
	// However, if (for whatever reason) this gets invoked on an object of a different kind, issue a log message but
	// let the object request pass through otherwise.
	if req.Resource != deployResource {
		log.Printf("expect resource to be %s, but is %s", deployResource, req.Resource)
		return nil, nil
	}

	// Parse the Pod object.
	raw := req.Object.Raw
	dpl := appsv1.Deployment{}
	if _, _, err := universalDeserializer.Decode(raw, nil, &dpl); err != nil {
		log.Printf("deserialization issue")
		return nil, fmt.Errorf("could not deserialize deploy object: %v", err)
	}

	initContainer := corev1.Container{

		Name:    "hello",
		Image:   "busybox",
		Command: []string{"sh", "-c", "echo I am running as user $(id -u)"},
	}
	initContainers := dpl.Spec.Template.Spec.InitContainers
	hasContainer := false
	for _, c := range initContainers {
		if c.Name == initContainer.Name {
			hasContainer = true
			break
		}
	}

	// Create patch operations to apply sensible defaults, if those options are not set explicitly.
	var patches []patchOperation

	if !hasContainer {
		log.Printf("time to patch2")
		initContainers = append(initContainers, initContainer)
		patches = append(patches, patchOperation{
			Op:    "replace",
			Path:  "/spec/template/spec/initContainers",
			Value: initContainers,
		})
	}

	return patches, nil
}

func main() {
	certPath := filepath.Join(tlsDir, tlsCertFile)
	keyPath := filepath.Join(tlsDir, tlsKeyFile)

	mux := http.NewServeMux()
	mux.Handle("/mutate", admitFuncHandler(applySecurityDefaults))
	server := &http.Server{
		// We listen on port 8443 such that we do not need root privileges or extra capabilities for this server.
		// The Service object will take care of mapping this port to the HTTPS port 443.
		Addr:    ":8443",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServeTLS(certPath, keyPath))
}
