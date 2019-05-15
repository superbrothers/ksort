# ksort

[![Build Status](https://travis-ci.org/superbrothers/ksort.svg?branch=master)](https://travis-ci.org/superbrothers/ksort)

This is a tool which does in-place sort of Kubernetes manifests by Kind.

## Description

When installing manifests, they should be sorted in a proper order by Kind. For example, Namespace object must be in the first place when installing them.

ksort sorts manfest files in a proper order by Kind, which is implementd by using `tiller.SortByKind()` in Kubernetes Helm.

## Usage

Sort manifest files in the `deploy` directory in the proper order, and output the result to the stdout.

```
$ ls ./manifests
deployment.yaml  ingress.yaml  namespace.yaml  service.yaml
$ ksort ./manifests
```

To pass the result into the stdin of `kubectl apply` command is also convenient.

```
$ ksort ./manifests | kubectl apply -f -
```

Sort manifests contained the manifest file that is specified.
```
$ ksort ./app.yaml
```

## Installation

You can download an archive file from [GitHub Releases](https://github.com/superbrothers/ksort/releases), then extract it and install a binary.

Or use `go get` as follows:

```
$ GO111MODULE=on go get github.com/superbrothers/ksort/cmd
```

## Installation as kubectl plugin

You can also use ksort as kubectl plugin. The name as kubectl plugin is `sort-manifests`.

1. Install [krew](https://github.com/GoogleContainerTools/krew) that is a plugin manager for kubectl
2. Run:

        kubectl krew install sort-manifests

3. Try it out

        kubectl sort-manifests -h

## License

This software is released under the MIT License and includes the work that is distributed in the Apache License 2.0.
