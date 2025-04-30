package example

import "github.com/ras0q/go-playground-test-overlay/caller"

func f(c caller.Caller) string {
	return c.Call()
}
