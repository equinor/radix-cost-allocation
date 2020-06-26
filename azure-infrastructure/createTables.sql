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
GO
/* 
    based on style guide: https://www.sqlstyle.guide/#do 
    add new column:
    ALTER TABLE cost.runs
    ADD cluster_memory_mega_bytes INTEGER;
*/