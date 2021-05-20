IF NOT EXISTS(SELECT 1 FROM sys.schemas WHERE schema_id = schema_id('cost'))
BEGIN
	EXEC sys.sp_executesql N'CREATE SCHEMA cost'
END

IF NOT EXISTS(SELECT 1 FROM sys.database_principals WHERE principal_id = DATABASE_PRINCIPAL_ID('datawriter'))
BEGIN
	CREATE ROLE datawriter
END

GRANT EXEC ON SCHEMA::cost TO datawriter