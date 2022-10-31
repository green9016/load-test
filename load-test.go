package main

import (
	"fmt"
	"runtime"
)

func main() {
	cpus := runtime.NumCPU()
	fmt.Println("cpus:", cpus)
	w := newWorkers(cpus)
	w.start()
}
