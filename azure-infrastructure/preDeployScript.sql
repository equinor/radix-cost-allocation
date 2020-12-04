/*
    This is the first script that is executed by the push workflows (to master and release branch)
*/

-- View cost.application_resource_run_aggregation is created in a later script, and uses WITH SCHEMABINDING
-- This will prevent the columns in source tables from being modified, e.g. changing length on nullability
-- We therefore drop the view to prevent any ALTERs in source tables from failing
IF EXISTS(select 1 from sys.views where object_id=OBJECT_ID('cost.application_resource_run_aggregation'))
BEGIN
	DROP VIEW cost.application_resource_run_aggregation
END
