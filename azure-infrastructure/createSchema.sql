IF NOT EXISTS(SELECT 1 FROM sys.schemas WHERE schema_id = schema_id('cost'))
BEGIN
	EXEC sys.sp_executesql N'CREATE SCHEMA cost'
END