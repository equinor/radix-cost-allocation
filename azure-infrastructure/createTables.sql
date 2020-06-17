IF (NOT EXISTS (SELECT *
FROM INFORMATION_SCHEMA.TABLES
WHERE TABLE_SCHEMA = 'radix_cost'
    AND TABLE_NAME = 'runs'))
                 BEGIN
    CREATE TABLE radix_cost.runs
    (
        id INT IDENTITY(1, 1) NOT NULL PRIMARY KEY,
        measured_time_utc DATETIME2,
    );
END

IF (NOT EXISTS (SELECT *
FROM INFORMATION_SCHEMA.TABLES
WHERE TABLE_SCHEMA = 'radix_cost'
    AND TABLE_NAME = 'required_resources'))
                 BEGIN
    CREATE TABLE radix_cost.required_resources
    (
        id INT IDENTITY(1, 1) NOT NULL PRIMARY KEY,
        run_id INT FOREIGN KEY REFERENCES radix_cost.runs(id),
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
/* based on style guide: https://www.sqlstyle.guide/#do */