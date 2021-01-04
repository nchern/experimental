package mock

import "testing"

func init() {
	selfTestMode = true
}

func TestMethodMockSingleCall(t *testing.T) {
	m := NewMethod()

	if count := m.CallCount(); count != 0 {
		t.Errorf("CallCount != 0: %d", count)
	}
	if once := m.CalledOnce(); once {
		t.Errorf("CalledOnce is true without calls")
	}

	m.Record(1, "foo")

	if count := m.CallCount(); count != 1 {
		t.Errorf("CallCount != 1: %d", count)
	}
	if c := m.CalledOnce(); !c {
		t.Errorf("CalledOnce is false with exactly one call")
	}

	notCalled := m.Call(1)
	if notCalled.Called() {
		t.Errorf("Must be false")
	}

	c := m.Call(0)
	if !c.Called() {
		t.Errorf("Must be true")
	}
	if !c.CalledWith(t, 1, "foo") {
		t.Errorf("Call arguments are not the same")
	}
}

func TestMethodMockMultipleCalls(t *testing.T) {
	m := NewMethod()

	m.Record("foo", "bar")
	m.Record(1, 2.3)
	m.Record(1, 10, []string{"foo", "bar"})

	if count := m.CallCount(); count != 3 {
		t.Errorf("CallCount != 1: %d", count)
	}
	if once := m.CalledOnce(); once {
		t.Errorf("CalledOnce is true on multiple calls")
	}

	var tests = []struct {
		name     string
		given    int
		expected []interface{}
	}{
		{"first call", 0, []interface{}{"foo", "bar"}},
		{"first call", 1, []interface{}{1, 2.3}},
		{"first call", 2, []interface{}{1, 10, []string{"foo", "bar"}}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if !m.Call(tt.given).CalledWith(t, tt.expected...) {
				t.Errorf("given(%d): expected %+v", tt.given, tt.expected)
			}
		})
	}
}
func TestMethodMockShouldReturnFalseIfCalledWithDifferentArguments(t *testing.T) {
	m := NewMethod()

	m.Record("foo", 1)

	c := m.Call(0)
	if c.CalledWith(t, "foo", "bar") {
		t.Errorf("Call arguments are not the same, but CalledOnce returned true")
	}
	if c.CalledWith(t, "foo", 1, "bar") {
		t.Errorf("More arguments were provided, but CalledOnce returned true")
	}
	if c.CalledWith(t) {
		t.Errorf("Was called with arguments, with empty tested")
	}
}
