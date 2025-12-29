package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// unused
// FIXME: vettools don't currently catch this
func f() {}

// nilness
// FIXME: vettools don't currently catch this
type X struct{ f int }

func Fnilness(x *X) {
	if x == nil {
		fmt.Println(x.f) // nil dereference in field selection
		return
	}
	fmt.Println(x.f)
}

// shadow
func BadRead(f *os.File, buf []byte) error {
	var err error
	for {
		_, err := f.Read(buf) // declaration of "err" shadows declaration
		if err != nil {
			break
		}
	}
	return err
}

func main() {
	// unmarshal
	var a int
	json.Unmarshal([]byte(`{"key": "value"}`), a)
}
