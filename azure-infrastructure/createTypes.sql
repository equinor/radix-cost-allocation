
IF NOT EXISTS(SELECT 1 FROM sys.table_types WHERE name='node_upsert_type' AND schema_id=SCHEMA_ID('cost'))
BEGIN
	CREATE TYPE cost.node_upsert_type AS TABLE(
		node_name VARCHAR(253) NOT NULL PRIMARY KEY,
		node_pool_name varchar(253) NOT NULL
	)
END

IF NOT EXISTS(SELECT 1 FROM sys.table_types WHERE name='container_upsert_type' AND schema_id=SCHEMA_ID('cost'))
BEGIN
	CREATE TYPE cost.container_upsert_type AS TABLE(
		container_id VARCHAR(253) NOT NULL,
		container_name VARCHAR(253) NOT NULL,
		pod_name VARCHAR(253) NOT NULL,
		application_name VARCHAR(253) NOT NULL,
		environment_name VARCHAR(253) NOT NULL,
		component_name VARCHAR(253) NOT NULL,
		wbs varchar(253) NOT NULL,
		started_at DATETIMEOFFSET(0) NOT NULL,
		last_known_running_at DATETIMEOFFSET(0) NOT NULL,
		cpu_request_millicores BIGINT NOT NULL,
		memory_request_bytes BIGINT NOT NULL,
		node_name varchar(253),
		PRIMARY KEY(container_id)
	)
END