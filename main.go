package main

import (
	"runtime"

	"github.com/rustyoz/fatcontroller/lib"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fatcontroller.Run()
}
