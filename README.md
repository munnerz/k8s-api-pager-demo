# k8s-test-controller

Add test resources to test your Kubernetes cluster setup.

[Check out the docs](https://srossross.github.io/k8s-test-controller/)
---

# Development

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

Running the application is as follows:

```bash
$ go run main.go -kubeconfig ~/.kube/config
```

Then we can go ahead and create a Tests and TestRuns!

```bash
$ kubectl create -f docs/tests.yaml
$ kubectl create -f docs/test-run.yaml
```
