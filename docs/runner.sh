#!/bin/bash

echo_stderr ()
{
    echo "$@" >&2
}

boom (){
  echo_stderr
  echo_stderr "    💥  $*"
  echo_stderr
  exit 1
}


POSITIONAL=()

NAMESPACE="default"
SELECTORS=""
while [[ $# -gt 0 ]]
do
  key="$1"

  case $key in
      -h|--help)
      DO_HELP="yes"
      shift # past argument
      ;;
      --namespace)
      NAMESPACE="$2"
      shift # past argument
      shift # past value
      ;;
      -s|--selectors)
      SELECTORS="$2"
      shift # past argument
      shift # past value
      ;;
      *)    # unknown option
      POSITIONAL+=("$1") # save it in an array for later
      shift # past argument
      ;;
  esac
done

set -- "${POSITIONAL[@]}" # restore positional parameters

if [ ! -z "${DO_HELP+x}" ]; then
  cat <<EOF
Create a Kubernetes TestRun and wait until it is finished.

Usage: bash $0 [options] [name] [command]

Options:
  --namespace            : The kubernetes namespace
  -s|--selectors         : map of test selectors in JSON format eg '{\"key\": \"value\"}'

Commands:
  watch  : Only watch a test, do not create it.
  run    : (default)
EOF
  exit 1
fi

TEST_NAME="$1"
COMMAND="$2"

if [ -z "${TEST_NAME+x}" ]; then
  boom "  Argument name is required. Run '$0 --help' for usage"
fi


KUBECTL="kubectl --namespace ${NAMESPACE}"

if [[ "$COMMAND" != "watch" ]]; then
  echo "Creating TestRun ${TEST_NAME}"
  cat <<EOF | ${KUBECTL} create -f -
apiVersion: srossross.github.io/v1alpha1
kind: TestRun
metadata:
  name:  ${TEST_NAME}
spec:
  selectors: ${SELECTORS}
EOF

fi

sleep 1

if ${KUBECTL} get testruns ${TEST_NAME} > /dev/null; then
  echo
else
  boom "Could not find test $TEST_NAME"
fi

${KUBECTL} get ev -l test-run=${TEST_NAME} --watch -o 'go-template={{.reason}} {{.message}}
' &
child=$!
sleep 1

until [[ $(${KUBECTL} get testruns ${TEST_NAME} -o jsonpath='{.status.status}') == "Complete" ]]; do
  sleep 5;
done

kill -TERM "$child" 2>/dev/null


MESSAGE=$(${KUBECTL} get testruns ${TEST_NAME} -o jsonpath='{.status.message}')
if [[ $(${KUBECTL} get testruns ${TEST_NAME} -o jsonpath='{.status.success}') == "true" ]]; then
  echo "Success: ${MESSAGE}"
else
  boom "Fail: ${MESSAGE}"
fi
