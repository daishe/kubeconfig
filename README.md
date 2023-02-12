# Kubeconfig

*simple Kubernetes config manager*

[![Latest stable version](https://img.shields.io/github/v/tag/daishe/kubeconfig?label=latest%20stable%20version&sort=semver)](https://github.com/daishe/kubeconfig/releases)
[![Latest release status](https://img.shields.io/github/actions/workflow/status/daishe/kubeconfig/release.yaml?label=release%20build&logo=github&logoColor=fff)](https://github.com/daishe/kubeconfig/actions/workflows/release.yaml)

[![Go reference](https://pkg.go.dev/badge/github.com/daishe/kubeconfig.svg)](https://pkg.go.dev/github.com/daishe/kubeconfig)
[![Go version](https://img.shields.io/github/go-mod/go-version/daishe/kubeconfig?label=version&logo=go&logoColor=fff)](https://golang.org/dl/)
[![Go report card](https://goreportcard.com/badge/github.com/daishe/kubeconfig)](https://goreportcard.com/report/github.com/daishe/kubeconfig)
[![License](https://img.shields.io/github/license/daishe/kubeconfig)](https://github.com/daishe/kubeconfig/blob/master/LICENSE)

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
