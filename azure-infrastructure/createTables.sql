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
		CONSTRAINT pk_container_resources PRIMARY KEY(container_id) ,
		CONSTRAINT fk_containers_node_id FOREIGN KEY(node_id) REFERENCES cost.nodes(id)
	)

	CREATE INDEX ix_container_last_known_running_at on cost.containers(last_known_running_at) include(started_at)
	CREATE INDEX ix_container_resources_node_id on cost.containers(node_id)
END

