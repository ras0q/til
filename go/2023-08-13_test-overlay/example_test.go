package example

import (
	"testing"

	// NOTE: must be imported as fixed name "impl"
	impl "github.com/ras0q/go-playground-test-overlay/caller/impl"
)

func Test_f(t *testing.T) {
	c := impl.New()

	// NOTE: IDE will show error here:
	// > c.SetMsg undefined (type caller.Caller has no field or method SetMsg)
	c.SetMsg("called from test")

	actual := f(c)
	expected := "called from test"

	if actual != expected {
		t.Errorf("f() = %v, want %v", actual, expected)
	}
}
