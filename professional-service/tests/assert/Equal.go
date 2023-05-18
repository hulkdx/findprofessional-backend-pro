package assert

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func EqualJSON(t *testing.T, a string, e any) {
	j, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	actual := strings.TrimSpace(a)
	expected := strings.TrimSpace(string(j))
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nExpected %v\nActual   %v", expected, actual)
	}
}

func Equal(t *testing.T, actual, expected any) {
	if objectsAreEqual(actual, expected) {
		return
	}
	t.Helper()
	t.Fatalf("\nExpected %v\nActual   %v", expected, actual)
}

func objectsAreEqual(actual, expected any) bool {
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
