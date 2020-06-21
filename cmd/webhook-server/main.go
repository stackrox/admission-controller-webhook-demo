package main

import (
	"errors"
	"fmt"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"net/http"
	"path/filepath"
)

const (
	tlsDir      = `/run/secrets/tls`
	tlsCertFile = `tls.crt`
	tlsKeyFile  = `tls.key`
)

var (
	connectionResource = metav1.GroupVersionResource{
		Version: "v1", 
		Resource: "connections",
	}
)

func validateAccountConnection(req *v1beta1.AdmissionRequest) (bool, error) {
	if req.Resource != connectionResource {
		log.Printf("expect resource to be %s", connectionResource)
		return false, nil
	}

	// Parse the Connection object.
	raw := req.Object.Raw
	connection := corev1.Pod{}
	if _, _, err := universalDeserializer.Decode(raw, nil, &connection); err != nil {
		return false, fmt.Errorf("could not deserialize connection object: %v", err)
	}

	// Handle account connection
	// Validate the user has access

	return false, errors.New("Not implemented yet")
}

func main() {
	certPath := filepath.Join(tlsDir, tlsCertFile)
	keyPath := filepath.Join(tlsDir, tlsKeyFile)

	mux := http.NewServeMux()
	mux.Handle("/validate", admitFuncHandler(validateAccountConnection))
	server := &http.Server{
		// We listen on port 8443 such that we do not need root privileges or extra capabilities for this server.
		// The Service object will take care of mapping this port to the HTTPS port 443.
		Addr:    ":8443",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServeTLS(certPath, keyPath))
}
