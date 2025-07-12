package main

//
// start the coordinator process, which is implemented
// in ../mr/coordinator.go
//
// go run mrcoordinator.go pg*.txt
//
// Please do not change this file.
//

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"6.5840/mr"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: mrcoordinator inputfiles...\n")
		os.Exit(1)
	}

	    var absFiles []string
    for _, f := range os.Args[1:] {
        abs, err := filepath.Abs(path.Join("Distributed_Systems/6.5840/src/main", f))
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error getting absolute path for %s: %v\n", f, err)
            os.Exit(1)
        }
        absFiles = append(absFiles, abs)
    }

	m := mr.MakeCoordinator[int](absFiles, 10)
	for m.Done() == false {
		time.Sleep(time.Second)
	}

	time.Sleep(time.Second)
}
