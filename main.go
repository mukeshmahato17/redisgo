package main

import (
	"flag"

	"github.com/mukeshmahato17/goredis/server"
)

func main() {
	listenAddr := flag.String("listenAddr", ":7379", "Listen Address")
	flag.Parse()

	server := server.NewServer(*listenAddr)
	server.Start()

}
