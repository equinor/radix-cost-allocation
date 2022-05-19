$targetSqlServerFQDN = "$(az sql server show -n ${env:SQL_SERVER_NAME} -g ${env:RESOURCE_GROUP} | jq -r .fullyQualifiedDomainName)"
Invoke-Sqlcmd -InputFile ${env:GITHUB_WORKSPACE}/azure-infrastructure/preDeployScript.sql -ServerInstance $targetSqlServerFQDN -Database ${env:DB_NAME} -Username ${env:SQL_ADMIN_USER_NAME} -password ${env:DB_ADMIN_PASSWORD}
Invoke-Sqlcmd -InputFile ${env:GITHUB_WORKSPACE}/azure-infrastructure/createSchema.sql -ServerInstance $targetSqlServerFQDN -Database ${env:DB_NAME} -Username ${env:SQL_ADMIN_USER_NAME} -password ${env:DB_ADMIN_PASSWORD}
Invoke-Sqlcmd -InputFile ${env:GITHUB_WORKSPACE}/azure-infrastructure/createTables.sql -ServerInstance $targetSqlServerFQDN -Database ${env:DB_NAME} -Username ${env:SQL_ADMIN_USER_NAME} -password ${env:DB_ADMIN_PASSWORD}
Invoke-Sqlcmd -InputFile ${env:GITHUB_WORKSPACE}/azure-infrastructure/createTypes.sql -ServerInstance $targetSqlServerFQDN -Database ${env:DB_NAME} -Username ${env:SQL_ADMIN_USER_NAME} -password ${env:DB_ADMIN_PASSWORD}
Invoke-Sqlcmd -InputFile ${env:GITHUB_WORKSPACE}/azure-infrastructure/createProcedures.sql -ServerInstance $targetSqlServerFQDN -Database ${env:DB_NAME} -Username ${env:SQL_ADMIN_USER_NAME} -password ${env:DB_ADMIN_PASSWORD}