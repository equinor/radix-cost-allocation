![build workflow](https://github.com/equinor/radix-cost-allocation/actions/workflows/build-push.yml/badge.svg) [![SCM Compliance](https://scm-compliance-api.radix.equinor.com/repos/equinor/radix-cost-allocation/badge)](https://developer.equinor.com/governance/scm-policy/)

# RADIX-COST-ALLOCATION

Pulls and stores container and node information into a SQL Server database.

We use helm charts to install on cluster

We use arm template and github action to create azure resources

# Deploy SQL scripts
The SQL Server database and objects are deployed on push to master and release branch.
All SQL scripts on azure-infrastructure must be idempotent.

## Deploy to cluster

Installation on cluster is handled by flux through [flux repo](https://github.com/equinor/radix-flux). 

tag in git repository (in master branch) - matching to the version of Version in charts/Chart.yaml

## Development Process

The `radix-cost-allocation` project follows a **trunk-based development** approach.

### 🔁 Workflow

- **External contributors** should:
  - Fork the repository
  - Create a feature branch in their fork

- **Maintainers** may create feature branches directly in the main repository.

### ✅ Merging Changes

All changes must be merged into the `master` branch using **pull requests** with **squash commits**.

The squash commit message must follow the [Conventional Commits](https://www.conventionalcommits.org/en/about/) specification.


## Release Process

Merging a pull request into `master` triggers the **Prepare release pull request** workflow.  
This workflow analyzes the commit messages to determine whether the version number should be bumped — and if so, whether it's a major, minor, or patch change.  

It then creates two pull requests:

- one for the new stable version (e.g. `1.2.3`), and  
- one for a pre-release version where `-rc.[number]` is appended (e.g. `1.2.3-rc.1`).

---

Merging either of these pull requests triggers the **Create releases and tags** workflow.  
This workflow reads the version stored in `version.txt`, creates a GitHub release, and tags it accordingly.

The new tag triggers the **Build and deploy Docker and Helm** workflow, which:

- builds and pushes a new container image and Helm chart to `ghcr.io`, and  
- uploads the Helm chart as an artifact to the corresponding GitHub release.

### Contributing

Want to contribute? Read our [contributing guidelines](./CONTRIBUTING.md)

### Generating mocks
We use gomock to generate mocks used in unit test.
You need to regenerate mocks if you make changes to any of the interface types used by the application; **Repository**

```
make mocks
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
