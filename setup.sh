#!/bin/sh

provider_name=$1
provider_domain=$2

sed -e 's@$PROVIDER_NAME@'"$provider_name"'@g' deployment/deployment.yaml > deployment/deployment.yaml
sed -e 's@$PROVIDER_DOMAIN@'"$provider_domain"'@g' deployment/deployment.yaml > deployment/deployment.yaml

sed -e 's@$PROVIDER_NAME@'"$provider_name"'@g' deployment/custom-resource-definition.yaml > deployment/custom-resource-definition.yaml
sed -e 's@$PROVIDER_DOMAIN@'"$provider_domain"'@g' deployment/custom-resource-definition.yaml > deployment/custom-resource-definition.yaml

sed -e 's@$PROVIDER_NAME@'"$provider_name"'@g' deployment/webhook-config.yaml > deployment/webhook-config.yaml
sed -e 's@$PROVIDER_DOMAIN@'"$provider_domain"'@g' deployment/webhook-config.yaml > deployment/webhook-config.yaml

rm setup.sh
