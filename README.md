# Kale

_Kale_ is a scripting tool using [Starlark](https://github.com/google/starlark-go) to build, package, deploy and operate applications running on kubernetes.

## Features

* Literally not much as of yet

## Why?

_Kale_ was created from the urge to re-iterate on deployment practices experienced at work. 
It started with a variety of helm charts ranging from `nginx-ingress` to `prometheus` and the inability to use tiller, because of a strict security concept caused by regulatory requirments. The initial deployment process was implemented in bash, which was invoking `helm template` and the likes to first generate the manifests and then apply them. A lot of the charts were also customized with patches in the process.

__Wouldn't it be great to have a scripting system were kubernetes and its ecosystem are first-class citizens?__

The ecosystem is also something at work people would constantly fight about. Some teams prefered _ksonnet_, others _helm_ and another team something new and fancy.

__Why choose? Rather combine and use whatever makes sense.__

If we have a scripting tool stitching everything together, there is another thought that crossed my mind. Even though a lot of companies practice DevOps, the artifacts deployed (e.g. helm charts) usually do not bundle operational knowledge. But if we can script kubernetes workflows:

__Why limit it to deployments? We can automate every process around the application.__

Good examples for this are etcd and cockroachdb, which can scale up automatically, but not down! Another simple example is triggering backups for databases.

The last and final realization came once I drafted, what is necessary to automate this and that process. Which lead me to my final conclusion:

__I guess we need a plugin system...__

Linters, integration tests et cetera not everything is coverable by _Kale_, but hopefully plugins can help alleviate that.



