# ksort

[![Build Status](https://travis-ci.org/superbrothers/ksort.svg?branch=master)](https://travis-ci.org/superbrothers/ksort)

This is a tool which does in-place sort of Kubernetes manifests by Kind.

## Description

When installing manifests, they should be sorted in a proper order by Kind. For example, Namespace object must be in the first place when installing them.

ksort sorts manfest files in a proper order by Kind.

## Usage

Sort manifest files in the `deploy` directory in the proper order, and output the result to the stdout.

```
$ ls ./manifests
deployment.yaml  ingress.yaml  namespace.yaml  service.yaml
$ ksort -f ./manifests
```

To pass the result into the stdin of `kubectl apply` command is also convenient.

```
$ ksort -f ./manifests | kubectl apply -f -
```

Sort manifests contained the manifest file that is specified.
```
$ ksort -f ./app.yaml
```

Sort manifests passed into stdin.
```
$ cat app.yaml | ksort -f-
```

## Installation

You can download an archive file from [GitHub Releases](https://github.com/superbrothers/ksort/releases), then extract it and install a binary.

Or use `go get` as follows:

```
$ go get github.com/superbrothers/ksort/cmd
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
