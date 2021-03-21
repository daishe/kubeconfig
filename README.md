# Kubeconfig

*simple Kubernetes config manager*

[![Latest release status](https://github.com/daishe/kubeconfig/actions/workflows/release.yaml/badge.svg)](https://github.com/daishe/kubeconfig/actions/workflows/release.yaml)
[![Latest release candidate status](https://github.com/daishe/kubeconfig/actions/workflows/release-candidate.yaml/badge.svg)](https://github.com/daishe/kubeconfig/actions/workflows/release-candidate.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/daishe/kubeconfig.svg)](https://pkg.go.dev/github.com/daishe/kubeconfig)
[![Go Report Card](https://goreportcard.com/badge/github.com/daishe/kubeconfig)](https://goreportcard.com/report/github.com/daishe/kubeconfig)

Instead of keeping one large kubectl config file with lots of entries, kubeconfig allows you to simply switch the entire file! This approach is much easier to use, especially when working with lots of temporary Kubernetes clusters.

## Usage

Just place all your kubectl config files in the `.kubeconfig` directory (under your home directory).

Then to list all configs, use

```sh
kubeconfig list
```

and to actually switch current kubectl config (located under `.kube/config` in your home directory), use

```sh
kubeconfig switch <config file name>
```

That's it!

## Help

To get the complete list of all commands, use

```sh
kubeconfig --help
```

or just simply type

```sh
kubeconfig
```

To get information about particular command including list of aliases and flags, use

```sh
kubeconfig <command> --help
```

or alternatively

```sh
kubeconfig help <command>
```

## License

Kubeconfig is open-sourced software licensed under the [Apache License 2.0](http://www.apache.org/licenses/).
