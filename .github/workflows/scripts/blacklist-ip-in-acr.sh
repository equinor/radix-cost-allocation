#!/bin/bash
runner_ip="$(curl --silent http://ifconfig.me/ip)"
echo "Removing runner IP '${runner_ip}' from Azure Container Registry '${ACR_NAME}' firewall whitelist"
az acr network-rule remove --name ${ACR_NAME} --subscription ${AZURE_SUBSCRIPTION_ID} --ip-address $runner_ip --only-show-errors --output none