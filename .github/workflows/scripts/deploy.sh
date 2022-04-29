#!/bin/sh

# Get Storage Account ID for sql logs
storageId=$(az storage account show --name "$SQL_LOG_STORAGE_NAME" --query "id" --output tsv)
az deployment group create -g "$RESOURCE_GROUP" -n "$DEPLOY_NAME" \
    --template-file ./azure-infrastructure/azuredeploy.json \
    --parameters sqlServerName="$SQL_SERVER_NAME" \
    --parameters databaseName="$DB_NAME" \
    --parameters sqlAdministratorLoginPassword="$DB_ADMIN_PASSWORD" \
    --parameters sqlAdministratorLoginUser="$SQL_ADMIN_USER_NAME" \
    --parameters adminGroupName="$ADMIN_GROUP_NAME" \
    --parameters adminGroupId="$ADMIN_GROUP_ID" \
    --parameters storageAccountId=$storageId