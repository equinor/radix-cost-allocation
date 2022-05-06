param (
    [Parameter(Mandatory=$true)][string]$FirewallRuleName
 )

az sql server firewall-rule delete `
    --name ${FirewallRuleName} `
    --resource-group ${env:RESOURCE_GROUP} `
    --server ${env:SQL_SERVER_NAME}