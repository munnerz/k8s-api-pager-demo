# Kubernetes custom API types Pager demo

This repository contains a demo application showing how you can use the
Kubernetes generator applications to automatically create supporting code for
your own apiserver or controller component.

It implements a very simple controller that will watch for 'Alert' resources,
and fire Pushbullet alerts for any unsent alerts in the API.

You can see the automation for generating the required code in the
[Makefile](Makefile), and in [hack/update-client-gen.sh](hack/update-client-gen.sh)
script.

## Generating code

If you change any API types, you must regenerate all of the supporting code.
This can be done by using the `generate` make target, e.g.

```bash
$ make generate
```

The generators tend not to **error** even if something you may consider to be a
problem occurs, thus it's important to ensure you can still build your
application after running the generators.

## Building

The actual golang app can be built with a simple `go build`. The generated
files required are committed to the repo to ensure it stays in sync. We should
also use a `verify` step to ensure that the generated files are in sync with
their respective types.go files, but for brevity have omitted that here.

```bash
$ go build
```

## Running

To run this demo, you will need a [Pushbullet](https://www.pushbullet.com)
account. You can sign up for free on their website. Once done, you will need to
create yourself an API key before running the application as follows:

```bash
$ ./k8s-api-pager-demo -pushbullet-token 'token-goes-here'
```

Optionally, you can set the `-apiserver` flag too. This particular demo will
**not** automatically detect credentials for the API server if it is running
within a cluster, once again for brevity.

Once started, you will also need to create the CustomResourceDefinition in the
target API server, as we have not implemented our own API server in this repo.

```bash
$ kubectl create -f docs/crd.yaml
```

Then we can go ahead and create an Alert!

```bash
$ kubectl create -f docs/test-alert.yaml
```
