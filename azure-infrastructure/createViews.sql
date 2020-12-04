
-- Set the options to support indexed views.
SET NUMERIC_ROUNDABORT OFF;
SET ANSI_PADDING, ANSI_WARNINGS, CONCAT_NULL_YIELDS_NULL, ARITHABORT,
   QUOTED_IDENTIFIER, ANSI_NULLS ON;

GO

CREATE OR ALTER VIEW cost.application_resource_run_aggregation
WITH SCHEMABINDING
AS
SELECT
	rr.run_id,
	r.measured_time_utc,
	r.cluster_cpu_millicores,
	r.cluster_memory_mega_bytes,
	rr.application,
	rr.wbs,
	SUM(ISNULL(rr.cpu_millicores,0) * ISNULL(rr.replicas, 0)) as cpu_millicores, 
	SUM(ISNULL(rr.memory_mega_bytes,0) * ISNULL(rr.replicas, 0)) as memory_mega_bytes,
	COUNT_BIG(*) as row_count
FROM
	cost.runs r
	INNER JOIN cost.required_resources rr on r.id=rr.run_id
GROUP BY
	rr.run_id,
	r.measured_time_utc,
	r.cluster_cpu_millicores,
	r.cluster_memory_mega_bytes,
	rr.application,
	rr.wbs

GO

if EXISTS(select 1 from sys.views where object_id=OBJECT_ID('cost.application_resource_run_aggregation'))
	AND NOT EXISTS(select 1 from sys.indexes WHERE object_id=OBJECT_ID('cost.application_resource_run_aggregation') AND index_id=1)
BEGIN
	CREATE UNIQUE CLUSTERED INDEX ux_application_resource_run_aggregation 
		ON cost.application_resource_run_aggregation(measured_time_utc, application, wbs, run_id)
END

