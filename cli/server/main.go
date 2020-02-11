package main

import (
	"flag"
	"fmt"

	"file_syncer/server"
)

var (
	port = flag.Int("port", 3333, "port")
	root = flag.String("root", "", "")
)

func main() {
	flag.Parse()

	if *root == "" {
		panic(fmt.Errorf("root is empty"))
	}

	fmt.Printf("port: %d\n", *port)
	fmt.Printf("root: %s\n", *root)

	serv, err := server.NewSyncServer(*port, *root)
	if err != nil {
		panic(err)
	}

	err = serv.Run()
	if err != nil {
		panic(err)
	}
}
