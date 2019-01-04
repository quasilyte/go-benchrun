package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func main() {
	flag.Usage = func() {
		lines := []string{
			"Usage: go-benchrun [flags...] oldBench newBench [go test args...]",
			"* oldBench is a pattern for `old` benchmark (w/o `Benchmark` prefix)",
			"* newBench is a pattern for `new` benchmark (w/o `Benchmark` prefix)",
			"",
			"Example:",
			"\t# compare BenchmarkOld and BenchmarkNew from foopkg package with -count=10",
			"\t$ go-benchrun Old New -v -count=10 foopkg",
			"",
			"Flags and defaults:",
		}
		for _, l := range lines {
			fmt.Fprintln(flag.CommandLine.Output(), l)
		}
		flag.PrintDefaults()
	}

	oldFile := flag.String("oldFile", "./old.txt",
		`old benchmark results destination file`)
	newFile := flag.String("newFile", "./new.txt",
		`new benchmark results destination file`)

	flag.Parse()

	oldBench := flag.Arg(0)
	newBench := flag.Arg(1)
	if oldBench == "" {
		log.Fatal("empty first positional arg (old bench name)")
	}
	if newBench == "" {
		log.Fatal("empty second positional arg (new bench name)")
	}

	// The "Benchmark" prefix is added implicitly.
	oldBench = "Benchmark" + oldBench
	newBench = "Benchmark" + newBench

	testArgs := flag.Args()[2:]
	fmt.Println("  Running old benchmarks:")
	runBenchmarks(*oldFile, oldBench, "", testArgs)
	fmt.Println("  Running new benchmarks:")
	runBenchmarks(*newFile, newBench, oldBench, testArgs)
	fmt.Println("  Benchstat results:")
	runBenchstat(*oldFile, *newFile)
}

func runBenchmarks(dstFile, selector, rename string, args []string) {
	testArgs := []string{
		"test",
		"-bench", selector,
	}
	testArgs = append(testArgs, args...)
	var output bytes.Buffer
	cmd := exec.Command("go", testArgs...)
	cmd.Stdout = io.MultiWriter(&output, os.Stdout)
	cmd.Stderr = io.MultiWriter(&output, os.Stderr)
	err := cmd.Run()
	out := output.Bytes()
	if err != nil {
		log.Fatalf("%q: run go test: %v: %s", selector, err, out)
	}
	if rename != "" {
		out = bytes.Replace(out, []byte(selector), []byte(rename), -1)
	}
	if err := ioutil.WriteFile(dstFile, out, 0666); err != nil {
		log.Fatalf("%q: write results: %v", selector, err)
	}
}

func runBenchstat(file1, file2 string) {
	out, err := exec.Command("benchstat", "-geomean", file1, file2).CombinedOutput()
	if err != nil {
		log.Fatalf("run benchstat: %v: %s", err, out)
	}
	fmt.Print(string(out))
}
