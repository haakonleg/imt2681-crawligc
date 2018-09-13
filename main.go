package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/haakonleg/imt2681-crawligc/crawligc"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [URL]\n", args[0])
		os.Exit(1)
	}

	// Get base URL and url
	index := strings.LastIndex(args[1], "/")
	baseURL := args[1][0 : index+1]
	url := args[1][index+1 : len(args[1])]

	// Run crawligc
	crawligc.CrawlIGC(baseURL, url)
}
