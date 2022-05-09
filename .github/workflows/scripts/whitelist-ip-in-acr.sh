#!/bin/bash

runner_ip="$(curl --silent http://ifconfig.me/ip)"
if [ -z "$(az acr network-rule list --name ${ACR_NAME} --subscription ${AZURE_SUBSCRIPTION_ID} | grep ${runner_ip})"]; then
    echo "Adding runner IP '${runner_ip}' to Azure Container Registry '${ACR_NAME}' firewall whitelist"
    az acr network-rule add --name ${ACR_NAME} --subscription ${AZURE_SUBSCRIPTION_ID} --ip-address $runner_ip
else
    echo "Runner is already whitelisted. Skipping..."
fi