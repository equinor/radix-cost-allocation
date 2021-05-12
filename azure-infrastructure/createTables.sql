IF (NOT EXISTS (SELECT *
FROM INFORMATION_SCHEMA.TABLES
WHERE TABLE_SCHEMA = 'cost'
    AND TABLE_NAME = 'runs'))
                 BEGIN
    CREATE TABLE cost.runs
    (
        id INT IDENTITY(1, 1) NOT NULL PRIMARY KEY,
        measured_time_utc DATETIME2,
        cluster_cpu_millicores INTEGER,
        cluster_memory_mega_bytes INTEGER,
    );
END

IF (NOT EXISTS (SELECT *
FROM INFORMATION_SCHEMA.TABLES
WHERE TABLE_SCHEMA = 'cost'
    AND TABLE_NAME = 'required_resources'))
                 BEGIN
    CREATE TABLE cost.required_resources
    (
        id INT IDENTITY(1, 1) NOT NULL PRIMARY KEY,
        run_id INT FOREIGN KEY REFERENCES cost.runs(id),
        wbs VARCHAR(256),
        application VARCHAR(256),
        environment VARCHAR(256),
        component VARCHAR(256),
        cpu_millicores INTEGER,
        memory_mega_bytes INTEGER,
        replicas INTEGER,
    );
END
IF (NOT EXISTS (SELECT *
FROM sys.indexes
WHERE name='req_resource_to_run' AND object_id = OBJECT_ID('cost.required_resources')))
        BEGIN
    CREATE INDEX req_resource_to_run ON cost.required_resources (run_id);
END


IF NOT EXISTS(SELECT 1 FROM sys.tables WHERE object_id=OBJECT_ID('cost.node_pool'))
BEGIN
	CREATE TABLE cost.node_pool(
		id INT IDENTITY CONSTRAINT pk_node_pool PRIMARY KEY,
		[name] VARCHAR(253) NOT NULL
	)

	CREATE UNIQUE INDEX ux_node_pool_name ON cost.node_pool([name])
END

IF NOT EXISTS(SELECT 1 FROM sys.tables WHERE object_id=OBJECT_ID('cost.node_pool_cost'))
BEGIN
	CREATE TABLE cost.node_pool_cost(
		id INT IDENTITY CONSTRAINT pk_node_pool_cost PRIMARY KEY,
		pool_id INT NOT NULL,
		cost INT NOT NULL,
		cost_currency char(3) NOT NULL,
		from_date DATETIMEOFFSET(0) NOT NULL,
		to_date DATETIMEOFFSET(0) NOT NULL,
		CONSTRAINT ck_from_to_date CHECK (to_date>from_date),
		CONSTRAINT fk_node_pool_cost_pool_id FOREIGN KEY(pool_id) REFERENCES cost.node_pool(id)
	)

	CREATE INDEX ix_node_pool_cost_pool_id on cost.node_pool_cost(pool_id)
END

IF NOT EXISTS(SELECT 1 FROM sys.tables WHERE object_id=OBJECT_ID('cost.nodes'))
BEGIN
	CREATE TABLE cost.nodes(
		id INT IDENTITY CONSTRAINT pk_nodes PRIMARY KEY,
		[name] VARCHAR(253) NOT NULL,
		pool_id INT NULL,
		CONSTRAINT fk_nodes_pool_id FOREIGN KEY(pool_id) REFERENCES cost.node_pool(id)
	)

	CREATE UNIQUE INDEX ux_nodes_name on cost.nodes([name])
	CREATE INDEX ix_nodes_pool_id on cost.nodes(pool_id)
END

IF NOT EXISTS(SELECT 1 FROM sys.tables WHERE object_id=OBJECT_ID('cost.containers'))
BEGIN
	CREATE TABLE cost.containers(
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
		node_id INT NOT NULL,
		CONSTRAINT pk_container_resources PRIMARY KEY NONCLUSTERED(container_id) ,
		CONSTRAINT fk_containers_node_id FOREIGN KEY(node_id) REFERENCES cost.nodes(id)
	)

	CREATE CLUSTERED INDEX cx_container_started_at on cost.containers(started_at)
	CREATE INDEX ix_container_resources_node_id on cost.containers(node_id)
END

