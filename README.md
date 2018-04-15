# ksort

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

Sort the manifests contained the manifest file that is specified.
```
$ ksort ./app.yaml
```

## Installation

```
$ go get github.com/superbrothers/ksort/cmd
```

## License

This software is released under the MIT License and includes the work that is distributed in the Apache License 2.0.
