# go-benchrun

Convenience wrapper around "go test" + "benchstat".

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

