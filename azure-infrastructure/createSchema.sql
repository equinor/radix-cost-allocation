IF NOT EXISTS(SELECT 1 FROM sys.schemas WHERE schema_id = schema_id('cost'))
BEGIN
	EXEC sys.sp_executesql N'CREATE SCHEMA cost'
END

IF NOT EXISTS(SELECT 1 FROM sys.database_principals WHERE principal_id = DATABASE_PRINCIPAL_ID('datawriter'))
BEGIN
	CREATE ROLE datawriter
END

GRANT EXEC ON SCHEMA::cost TO datawriter
GRANT SELECT, INSERT, UPDATE, DELETE ON SCHEMA::cost TO datawriter

IF NOT EXISTS(SELECT 1 FROM sys.database_principals WHERE principal_id = DATABASE_PRINCIPAL_ID('datareader'))
BEGIN
	CREATE ROLE datareader
END

GRANT SELECT ON SCHEMA::cost TO datareader


IF NOT EXISTS(SELECT 1 FROM sys.database_principals WHERE name = 'radix-id-vulnerability-scan-writer-$(RADIX_ZONE)')
BEGIN
    CREATE USER [radix-id-vulnerability-scan-writer-$(RADIX_ZONE)] FROM EXTERNAL PROVIDER;
END
ALTER ROLE datawriter ADD MEMBER [radix-id-vulnerability-scan-writer-$(RADIX_ZONE)]

IF NOT EXISTS(SELECT 1 FROM sys.database_principals WHERE name = 'radix-id-vulnerability-scan-reader-$(RADIX_ZONE)')
BEGIN
    CREATE USER [radix-id-vulnerability-scan-reader-$(RADIX_ZONE)] FROM EXTERNAL PROVIDER;
END
ALTER ROLE datareader ADD MEMBER [radix-id-vulnerability-scan-reader-$(RADIX_ZONE)]
