package assert

import (
	"bytes"
	"reflect"
	"testing"
)

func Equal(t *testing.T, actual, expected interface{}) {
	if objectsAreEqual(actual, expected) {
		return
	}
	t.Helper()
	t.Fatalf("Expected %s\nActual %s", expected, actual)
}

func objectsAreEqual(actual, expected interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	exp, eok := expected.([]byte)
	act, aok := actual.([]byte)

	if eok && aok {
		return bytes.Equal(exp, act)
	}

	return reflect.DeepEqual(expected, actual)
}
