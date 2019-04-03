# Kompass

_Kompass_ is a CLI tool using [Starlark](https://github.com/google/starlark-go) to 
package kubernetes deployments and relevent tasks.

## Problems kompass is trying to solve

_Kompass_ was created to solve some very specific problems experienced at work:
* The kubernetes ecosystem is fairly fragmented, so are technology choices of teams (e.g. helm, ksonnet, kustomize, ...)
* Different application have different needs (e.g. scaling down etcd or cockroachdb requires additional effort)
* Some deployments bring along chores or tasks (e.g. migration of databases or an application not yet being kubernetes-ready and needing an ungraceful restart)
* Tasks, scripts, helm-chart etc are not always next to the code, from DevOps-perspective we would like to bundle the application including how it is deployed and operated and make it available to dependent bundles

