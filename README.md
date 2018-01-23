# ksort

This is a tool which does in-place sort of Kubernetes manifests by Kind.

## Description

When installing manifests, they should be sorted in a proper order by Kind. For example, Namespace object must be in the first place when installing them.

ksort sorts manfest files in a proper order by Kind, which is implementd by using `tiller.SortByKind()` in Kubernetes Helm.

## Usage

Sort manifest files in the `deploy` directory in the proper order, and output the result to the stdout.

```
$ ls ./deploy
deployment.yaml  ingress.yaml  namespace.yaml  service.yaml
$ ksort ./deploy
```

To pass the result into the stdin of `kubectl apply` command is also convenient.

```
$ ksort ./deploy | kubectl apply -f -
```

## Installation

```
$ go get github.com/superbrothers/ksort
```

## License

This software is released under the MIT License and includes the work that is distributed in the Apache License 2.0.
