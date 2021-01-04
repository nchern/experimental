package mock

import (
	"fmt"
	"testing"
)

var (
	// Equal is a settable comparator to compare two values for tests
	// Hint: you can use a function with the same signature from any of a test library to get
	// nicer output regarding non equal values, e.g. from stretchr/testify.Equal
	Equal func(t *testing.T, expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool = equal

	selfTestMode = false
)

// Method is a helper mock to mock a method
type Method interface {
	CallCount() int
	CalledOnce() bool
	Call(int) Call

	Record(args ...interface{})
}

type method struct {
	calls []*call
}

// NewMethod creates a new instance of Method moc
func NewMethod() Method { return &method{} }

// CallCount returns how many times this mock was called
func (m *method) CallCount() int { return len(m.calls) }

// CalledOnce returns true if this method was called exactly once
func (m *method) CalledOnce() bool { return m.CallCount() == 1 }

// Call returns the i-th Call instance. If this method was called less than a given i,
// it returns a mock that indicates no call
func (m *method) Call(i int) Call {
	if i < 0 {
		panic(fmt.Sprintf("i must not be less than zero; actual=%d", i))
	}
	if i >= len(m.calls) {
		return &call{}
	}
	return m.calls[i]
}

// Record records a method call with a given set of arguments
func (m *method) Record(args ...interface{}) {
	m.calls = append(m.calls, &call{called: true, args: args})
}

// Call represents a concrete mocked call
type Call interface {
	Called() bool
	CalledWith(t *testing.T, args ...interface{}) bool
}

type call struct {
	called bool
	args   []interface{}
}

// Called returns true if this call happened
func (c *call) Called() bool { return c.called }

// CalledWith returns true if this call was called with the same set of arguments
func (c *call) CalledWith(t *testing.T, args ...interface{}) bool {
	if !c.Called() {
		return false
	}
	if len(c.args) != len(args) {
		if !selfTestMode {
			t.Errorf("len(c.args) = %d; len(args) = %d", len(c.args), len(args))
		}
		return false
	}
	for i, arg := range c.args {
		if !Equal(t, args[i], arg) {
			return false
		}
	}
	return true
}
