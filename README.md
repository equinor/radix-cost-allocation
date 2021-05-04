# RADIX-COST-ALLOCATION

Pulls cost data from cluster prometheus instance for each application and push it into a sql database

We use helm charts to install on cluster

We use arm template and github action to create azure resources

# Deploy SQL scripts
The SQL Server database and objects are deployed on push to master and release branch.
All SQL scripts on azure-infrastructure must be idempotent.

## Deploy to cluster

Installation on cluster is handled by flux through [flux repo](https://github.com/equinor/radix-flux). Before being installed, it requires that there exist a namespace called `radix-cost-allocation`. In that namespace there must be a secret called `cost-db-secret` that contains the database password. This is handled through the setup script in [radix-platform](https://github.com/equinor/radix-platform)

tag in git repository (in master branch) - matching to the version of Version in docs/docs.go

## Update version

`tag` in git repository (in master branch) - matching to the version of appVersion in `charts/Chart.yaml`

Run following command to set tag (with corresponding version)
```
git tag v1.0.0
git push origin v1.0.0
```

## Debugging locally

Create a copy of .env.template and name it .env. Set variables to allow local debugging. This file is ignored by git.

Prometheus operator must be port forwarded to a local port:
```
k port-forward pod/prometheus-prometheus-operator-prometheus-0 9090:9090
```