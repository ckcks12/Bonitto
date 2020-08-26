#!/bin/bash

if ! which yq > /dev/null; then
  echo "install yq first"
  echo "pip install yq"
  exit 1
fi

secret_name=$(kubectl get sa/bonitto -o yaml | yq -r '.secrets[0].name')
kubectl get secret/${secret_name} -o yaml > tmp
ca_crt=$(cat tmp | yq -r '.data["ca.crt"]')
token=$(cat tmp | yq -r '.data.token')

sudo mkdir -p /var/run/secrets/kubernetes.io/serviceaccount
echo "${ca_crt}" | base64 -d > tmp
sudo mv tmp /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
echo "${token}" | base64 -d > tmp
sudo mv tmp /var/run/secrets/kubernetes.io/serviceaccount/token

echo "done"
