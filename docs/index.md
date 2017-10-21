
# Kubernetes Test Controller

- Run end to end tests on your k8s application
- Run targeted component tests after a rollout or change

## Motivation

### Why Test?

- Rollout your CI/CD pipeline with confidence
- Validate that your configuration is correct
  - Make sure your username and password work correctly
  - Make sure an incorrect username and password does not work
- Assert that your services are up and correctly load balancing

## Cluster Installation

To use this test controller run:

```sh
kubectl create -f https://srossross.github.io/k8s-test-controller/controller.yaml
kubectl --namespace kube-system rollout status deploy test-controller-deployment --watch
```

That's it! Now you can get started running tests. This controller adds two
custom resources to the Kubernetes cluster - A `Test` and a `TestRun`.


## Resources

### Kind: Test

A `Test` resource will look something like this:

```yaml
# File test-success.yaml
apiVersion: srossross.github.io/v1alpha1
kind: Test
metadata:
  name: test-success
  labels:   #  This can be used to filter tests in a testrun
    app: mytest
spec:
  template: # This is just like a Kubernetes Job or Deployment template
    spec:
      containers:
      - name: alpine
        image: alpine
        command: [echo, hello]
      restartPolicy: Never
```

It will contain `Pod` definition inthe  `spec.template` field. The `Test` will
will instantiate this pod when a new  testrun is created.

To add this test to your cluster, first create the file `./test-success.yaml` then run:

```
kubectl create -f ./test-success.yaml
```

A `Test` is comparable to a Kubernetes `[Job](https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/)`.
Unlike a `Job`, a `Test` will not run by itself. For that you will need to create
a `TestRun`.

### Kind: TestRun

A creating a `TestRun` will instantiate all `Tests` that matched
by its optional `selector`. If omitted, all tests in the namespace will be run.

```yaml
# File: test-run.yaml
apiVersion: srossross.github.io/v1alpha1
kind: TestRun
metadata:
  name: test-run-1
spec:
  max-jobs: 1     # The maximum number of test jobs to run simultaneously
  selector:       # Optional -- This will filter the tests to run.
    matchLabels:
      app: mytest
```

To add this to your cluster, first create the file `test-run.yaml` and trigger
test runs with:

```
kubectl create -f ./test-run.yaml
```

The controller will now start running your tests.

## Inspecting a TestRun

A TestRun will emit Kubernetes events. You can inspect these events by either running
`kubectl describe`. eg:

```
kubectl describe testrun test-run-1
```

Or inspecting the events  eg:

```
kubectl get events -l test-run=test-run-1 --watch
```

### Monitoring Tests with runner.sh:


For now you can use our bash script. This script will create and watch a TestRun and wait until it is finished:

```
curl --fail https://srossross.github.io/k8s-test-controller/runner.sh > ./runner.sh
bash ./runner.sh test-run-2
```


Waiting on [This CRD proposal](https://github.com/kubernetes/kubernetes/issues/38113) to be able to run:

```sh
# I wish! But this will not work ... yet
kubectl rollout status testrun test-run-1 --watch
```


## Motivation Pt2

### Comparison with helm tests

This test controller works great with helm.
You can create **Test** and **TestRun** resources from your charts, the only
difference is how a test is launched. you can now use `kubectl create -f testrun.yaml`
instead  of `helm test`.

This test controller can be used in any k8s cluster without requiring helm.

### Tests are kubernetes resources.

This means you can update your tests when you update your deployments.

For example if you have a CI/CD pipeline that pushes a new image to your
micro-service e.g.:

```
kubectl set image deploy mydeploy *=srossross/mynewimage
```

You may also want to update the tests that get run:

```
kubectl patch test test-mydeploy -p '{"spec":{"template":{"spec":{"containers":[{"name":"server","image":"srossross/mynewimage.test"}]}}}}'
```

### TestRuns can filter out tests.

In our CI/CD scenario deploying to our cluster, you may only want to run tests
that are pertinent to the updated components, rather than all of the tests
in your cluster.

You can do this with `TestRun` selectors. Lets say we created our resources with
helm, and we labeled all resources in a chart with `chart=mychart`

e.g:

```yaml
# File: test-run.yaml
apiVersion: srossross.github.io/v1alpha1
kind: TestRun
metadata:
  name: test-run-1
spec:
  selector:
    matchLabels:
      chart: mychart
```



---


## Example Test

In `wordpress/tests/test-mariadb-connection.yaml`:

```
apiVersion: srossross.github.io/v1alpha1
kind: Test
metadata:
  name: "credentials-test"
  labels:
    app: test
spec:
  template:
    spec:
      containers:
      - name: credentials-test
        image: mariadb
        env:
          - name: MARIADB_HOST
            value: mariadb
          - name: MARIADB_PORT
            value: "3306"
          - name: WORDPRESS_DATABASE_NAME
            value: wordpress
          - name: WORDPRESS_DATABASE_USER
            value: root
          - name: WORDPRESS_DATABASE_PASSWORD
            valueFrom:
              secretKeyRef:
                name: mariadb-secrets
                key: mariadb-password
        command: ["sh", "-c", "mysql --host=$MARIADB_HOST --port=$MARIADB_PORT --user=$WORDPRESS_DATABASE_USER --password=$WORDPRESS_DATABASE_PASSWORD"]
      restartPolicy: Never
```

### Steps to Run a Test Suite on this Resource

1. ```sh
$ cat <<EOF | kubectl create -f -
apiVersion: srossross.github.io/v1alpha1
kind: TestRun
metadata:
  name: test-run-1
EOF
```

## runner.sh

```
Create a Kubernetes TestRun and wait until it is finished.

Usage: bash runner.sh [options] [name] [command]

Options:
  --namespace            : The kubernetes namespace
  -s|--selectors         : map of test selectors in JSON format eg '{\"key\": \"value\"}'

Commands:
  watch  : Only watch a test, do not create it.
  run    : (default)
```


## TODO

 - Parameterized TestRuns
 - Implement `maxfail` spec, to stop running tests after a number of builds
 - Aggregate logging from all pods during test runs
