# go-benchrun

Convenience wrapper around "go test" + [benchstat](https://godoc.org/golang.org/x/perf/cmd/benchstat).

Run benchmarking in 1 simple command.

## Installation & Quick start

This install `go-benchrun` binary under your `$GOPATH/bin`:

```bash
go get github.com/Quasilyte/go-benchrun
```

If `$GOPATH/bin` is under your system `$PATH`, `go-benchrun` command should be available after that.<br>
This should print the help message:

```bash
$ go-benchrun --help
Usage: go-benchrun [flags...] oldBench newBench [go test args...]
* oldBench is a pattern for `old` benchmark (w/o `Benchmark` prefix)
* newBench is a pattern for `new` benchmark (w/o `Benchmark` prefix)

Example:
	# compare BenchmarkOld and BenchmarkNew from foopkg package with -count=10
	$ go-benchrun Old New -v -count=10 foopkg

Flags and defaults:
  -newFile string
    	new benchmark results destination file (default "./new.txt")
  -oldFile string
    	old benchmark results destination file (default "./old.txt")
```

See "Workflow" section for more usage info".

## Workflow

Without `go-benchrun`, your workflow is either of these two:

1. Rely on VSC.
	* Store old benchmark results (run go test).
	* Apply optimizations.
	* Run benchmarks again with optimized code.
	* Compare results with `benchstat`.
	* If you need to switch between implementations, you use stash and/or branches.
2. Rely on renaming.
	* Use one branch, two different benchmarks.
	* Collect results from both benchmarks.
	* Before running `benchstat`, rename benchmarks, so their name matches.
	
`go-benchrun` automates (2) scheme for you.

1. First, it runs `-bench=oldBench` and saves results to `oldFile`.
2. Then it runs `-bench=newBench` and saves results to `newFile`.
3. After that, it renames `newBench` from `newFile` to `oldBench`.
4. Finally, it runs `benchstat -geomean oldFile newFile`.

For example, lets say that you have this test file with benchmarks:

```go
package benchmark

import (
	"testing"
)

//go:noinline
func emptySliceLit() []int {
	return []int{}
}

//go:noinline
func makeEmptySlice() []int {
	return make([]int, 0)
}

func BenchmarkEmptySliceLit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = emptySliceLit()
	}
}

func BenchmarkMakeEmptySlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = makeEmptySlice()
	}
}
```

In order to compare `BenchmarkEmptySliceLit` and `BenchmarkMakeEmptySlice` you do:

```bash
$ benchrun EmptySliceLit MakeEmptySlice -v -count=5 .
	Running old benchmarks:
goos: linux
goarch: amd64
BenchmarkEmptySliceLit-8   	300000000	         6.08 ns/op
BenchmarkEmptySliceLit-8   	200000000	         5.83 ns/op
BenchmarkEmptySliceLit-8   	300000000	         5.85 ns/op
BenchmarkEmptySliceLit-8   	300000000	         5.89 ns/op
BenchmarkEmptySliceLit-8   	300000000	         5.71 ns/op
PASS
ok  	_/home/quasilyte/CODE/go/bench	11.184s
	Running new benchmarks:
goos: linux
goarch: amd64
BenchmarkMakeEmptySlice-8   	200000000	         7.68 ns/op
BenchmarkMakeEmptySlice-8   	200000000	         8.20 ns/op
BenchmarkMakeEmptySlice-8   	200000000	         8.31 ns/op
BenchmarkMakeEmptySlice-8   	200000000	         7.98 ns/op
BenchmarkMakeEmptySlice-8   	200000000	         8.69 ns/op
PASS
ok  	_/home/quasilyte/CODE/go/bench	12.128s
	Benchstat results:
name             old time/op  new time/op  delta
EmptySliceLit-8  5.87ns ± 4%  8.17ns ± 6%  +39.17%  (p=0.008 n=5+5)
```

If there are unit tests (non-benchmarks), you can specify `-run` flag for `go test`, as usual.
