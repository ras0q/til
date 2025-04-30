# go-playground-test-overlay

without `-overlay`

```sh
$ go test . -v
# github.com/ras0q/go-playground-test-overlay [github.com/ras0q/go-playground-test-overlay.test]
./example_test.go:11:4: c.SetMsg undefined (type caller.Caller has no field or method SetMsg)
FAIL    github.com/ras0q/go-playground-test-overlay [build failed]
FAIL
```

with `-overlay`

```sh
$ go test . -overlay overlay-mock.json -v
=== RUN   Test_f
2023/08/13 15:02:46 called from test
--- PASS: Test_f (0.00s)
PASS
ok      github.com/ras0q/go-playground-test-overlay
```
