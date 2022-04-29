param (
    [Parameter(Mandatory=$true)][string]$FirewallRuleName
 )

$runnerIp = (Invoke-WebRequest -uri "http://ifconfig.me/ip").Content
az sql server firewall-rule create `
    --end-ip-address $runnerIp `
    --start-ip-address $runnerIp `
    --name ${FirewallRuleName} `
    --resource-group ${env:RESOURCE_GROUP} `
    --server ${env:SQL_SERVER_NAME}