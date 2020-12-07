# RADIX-COST-ALLOCATION

Pulls cost data from cluster prometheus instance for each application and push it into a sql database

We use helm charts to install on cluster

We use arm template and github action to create azure resources

# Deploy SQL scripts
The SQL Server database and objects are deployed on push to master and release branch.
All SQL scripts on azure-infrastructure must be idempotent.

## Deploy to cluster

Installation on cluster is handled by flux through [flux repo](https://github.com/equinor/radix-flux). Before being installed, it requires that there exist a namespace called `radix-cost-allocation`. In that namespace there must be a secret called `cost-db-secret` that contains the database password. This is handled through the setup script in [radix-platform](https://github.com/equinor/radix-platform)
