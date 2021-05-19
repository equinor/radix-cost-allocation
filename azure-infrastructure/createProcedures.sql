CREATE OR ALTER PROCEDURE cost.node_upsert_bulk
@nodes cost.node_upsert_type READONLY
AS

SET NOCOUNT ON;
SET XACT_ABORT ON;

INSERT INTO cost.node_pool(name)
SELECT DISTINCT
	n.node_pool_name
FROM
	@nodes n
WHERE
	n.node_pool_name <> ''
	AND NOT EXISTS(SELECT 1 FROM cost.node_pool t WHERE t.name=n.node_pool_name);

WITH nodes AS (
	SELECT
		t.node_name, p.id as pool_id
	FROM
		@nodes t
		LEFT JOIN cost.node_pool p ON t.node_pool_name = p.[name]
)
MERGE INTO cost.nodes as t
USING nodes as s
	ON s.node_name = t.name
WHEN NOT MATCHED BY TARGET THEN
	INSERT(name, pool_id)
	VALUES(s.node_name, s.pool_id)
WHEN MATCHED AND s.pool_id IS NOT NULL THEN
	UPDATE SET
		pool_id=s.pool_id;

GO


CREATE OR ALTER PROCEDURE cost.container_upsert_bulk
@containers cost.container_upsert_type READONLY
AS

SET NOCOUNT ON;
--SET XACT_ABORT ON;

DECLARE @nodes cost.node_upsert_type
INSERT INTO @nodes(node_name, node_pool_name)
SELECT DISTINCT node_name, '' FROM @containers

exec cost.node_upsert_bulk @nodes=@nodes;

WITH containers AS(
	select
		t.container_id,
		t.container_name,
		t.pod_name,
		t.application_name,
		t.environment_name,
		t.component_name,
		t.wbs,
		t.started_at,
		t.last_known_running_at,
		t.cpu_request_millicores,
		t.memory_request_bytes,
		n.id as node_id
	from
		@containers t
		inner join cost.nodes n on t.node_name=n.[name]
)
MERGE INTO cost.containers as t
USING containers as s
	ON t.container_id = s.container_id
WHEN NOT MATCHED BY TARGET THEN
	INSERT(container_id, container_name, pod_name, application_name, environment_name, component_name, wbs, started_at, last_known_running_at, cpu_request_millicores, memory_request_bytes, node_id)
	VALUES(s.container_id, s.container_name, s.pod_name, s.application_name, s.environment_name, s.component_name, s.wbs, s.started_at, s.last_known_running_at, s.cpu_request_millicores, s.memory_request_bytes, s.node_id)
WHEN MATCHED THEN
	UPDATE SET
		last_known_running_at=s.last_known_running_at,
		wbs=s.wbs;

go

