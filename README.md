![build workflow](https://github.com/equinor/radix-cost-allocation/actions/workflows/build-push.yml/badge.svg) 

# RADIX-COST-ALLOCATION

Pulls and stores container and node information into a SQL Server database.

We use helm charts to install on cluster

We use arm template and github action to create azure resources

# Deploy SQL scripts
The SQL Server database and objects are deployed on push to master and release branch.
All SQL scripts on azure-infrastructure must be idempotent.

## Deploy to cluster

Installation on cluster is handled by flux through [flux repo](https://github.com/equinor/radix-flux). Before being installed, it requires that there exist a namespace called `radix-cost-allocation`. In that namespace there must be a secret called `cost-db-secret` that contains the database password. This is handled through the setup script in [radix-platform](https://github.com/equinor/radix-platform)

tag in git repository (in master branch) - matching to the version of Version in docs/docs.go

## Developing

You need Go installed. Make sure `GOPATH` and `GOROOT` are properly set up.

Also needed:

- [`gomock`](https://github.com/golang/mock) (GO111MODULE=on go get github.com/golang/mock/mockgen@v1.5.0)

Clone the repo into your `GOPATH` and run `go mod download`.

### Contributing

Want to contribute? Read our [contributing guidelines](./CONTRIBUTING.md)

### Generating mocks
We use gomock to generate mocks used in unit test.
You need to regenerate mocks if you make changes to any of the interface types used by the application; **Repository**

Repository:
```
$ mockgen -source ./pkg/repository/repository.go -destination ./pkg/repository/mock/repository.go -package mock
```
listers:
```
$ mockgen -source ./pkg/listers/limitrange.go -destination ./pkg/listers/mock/limitrange.go -package mock
$ mockgen -source ./pkg/listers/node.go -destination ./pkg/listers/mock/node.go -package mock
$ mockgen -source ./pkg/listers/pod.go -destination ./pkg/listers/mock/pod.go -package mock
$ mockgen -source ./pkg/listers/radixregistration.go -destination ./pkg/listers/mock/radixregistration.go -package mock
$ mockgen -source ./pkg/listers/containerbulkdto.go -destination ./pkg/listers/mock/containerbulkdto.go -package mock
$ mockgen -source ./pkg/listers/nodebulkdto.go -destination ./pkg/listers/mock/nodebulkdto.go -package mock
```

## Update version

`tag` in git repository (in master branch) - matching to the version of appVersion in `charts/Chart.yaml`

Run following command to set tag (with corresponding version)
```
git tag v1.0.0
git push origin v1.0.0
```

## Debugging locally

Create a copy of .env.template and name it .env. Set variables to allow local debugging. This file is ignored by git.


---------

[Security notification](./SECURITY.md)