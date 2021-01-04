package mock

import (
	"bytes"
	"reflect"
	"testing"
)

func equal(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return objectsAreEqual(expected, actual)
}

// Taken from https://github.com/stretchr/testify/blob/92707c0b2d501c60de82176c4aa1cf880abac720/assert/assertions.go#L58
func objectsAreEqual(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
}
