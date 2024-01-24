package assert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
)

/**
 * Ignore the order of JSON arrays.
 * For example the below jsons would be equals:
 *
 * [{"id":1,"email":"test1@gmail.com"},{"id":2,"email":"test2@gmail.com"}]
 * [{"id":2,"email":"test2@gmail.com"},{"id":1,"email":"test1@gmail.com"}]
 */
func EqualAnyOrderJSON(t *testing.T, a string, e any) {
	if reflect.TypeOf(e).Kind() == reflect.Slice {
		// convert a into []map[string]any
		var a1 []map[string]any
		err := json.Unmarshal([]byte(a), &a1)
		if err != nil {
			panic(err)
		}
		// convert e into []map[string]any
		j, err := json.Marshal(e)
		if err != nil {
			panic(err)
		}
		var e1 []map[string]any
		err = json.Unmarshal([]byte(j), &e1)
		if err != nil {
			panic(err)
		}
		sort.Slice(a1, func(i, j int) bool {
			return fmt.Sprintf("%v", a1[i]) < fmt.Sprintf("%v", a1[j])
		})
		sort.Slice(e1, func(i, j int) bool {
			return fmt.Sprintf("%v", e1[i]) < fmt.Sprintf("%v", e1[j])
		})
		Equal(t, a1, e1)
	} else {
		EqualJSON(t, a, e)
	}
}

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

func NotEqualJSON(t *testing.T, a string, e any) {
	j, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	actual := strings.TrimSpace(a)
	expected := strings.TrimSpace(string(j))
	if reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nEqualJSON\nExpected %v\nActual   %v", expected, actual)
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

func EqualAnyOrder(t *testing.T, actual, expected interface{}) {
	t.Helper()
	extraA, extraB := diffLists(actual, expected)
	if len(extraA) == 0 && len(extraB) == 0 {
		return
	}
	t.Fatalf("\nExpected %v\nActual   %v", expected, actual)
}

func diffLists(listA, listB interface{}) (extraA, extraB []interface{}) {
	aValue := reflect.ValueOf(listA)
	bValue := reflect.ValueOf(listB)

	aLen := aValue.Len()
	bLen := bValue.Len()

	// Mark indexes in bValue that we already used
	visited := make([]bool, bLen)
	for i := 0; i < aLen; i++ {
		element := aValue.Index(i).Interface()
		found := false
		for j := 0; j < bLen; j++ {
			if visited[j] {
				continue
			}
			if objectsAreEqual(bValue.Index(j).Interface(), element) {
				visited[j] = true
				found = true
				break
			}
		}
		if !found {
			extraA = append(extraA, element)
		}
	}

	for j := 0; j < bLen; j++ {
		if visited[j] {
			continue
		}
		extraB = append(extraB, bValue.Index(j).Interface())
	}

	return
}
