package main

import (
	"flag"
	"fmt"

	"file_syncer/client"
)

var (
	serverAddr = flag.String("addr", "", "server addr")
	rootDir    = flag.String("root", "", "root directory")
)

func main() {
	flag.Parse()
	if *serverAddr == "" {
		panic(fmt.Errorf("addr is empty"))
	}

	if *rootDir == "" {
		panic(fmt.Errorf("root directory is empty"))
	}

	fmt.Printf("addr: %s\n", *serverAddr)
	fmt.Printf("rootDir: %s\n", *rootDir)

	cl, err := client.NewSyncClient(*serverAddr, *rootDir)
	if err != nil {
		panic(err)
	}

	err = cl.Sync()
	if err != nil {
		panic(err)
	}
}
