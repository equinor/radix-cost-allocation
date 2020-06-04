CREATE SCHEMA radix_cost;
GO
CREATE TABLE radix_cost.runs
(
    id INT IDENTITY(1, 1) NOT NULL PRIMARY KEY,
    measured_time_utc DATETIME2,
);
GO
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
GO
/* based on style guide: https://www.sqlstyle.guide/#do */