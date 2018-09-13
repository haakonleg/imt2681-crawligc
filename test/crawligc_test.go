package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/haakonleg/imt2681-crawligc/crawligc"
)

const pageDir = "../pages"
const port = ":8080"

func srvStatic() {
	// Start web server
	fs := http.FileServer(http.Dir(pageDir))
	http.Handle("/", fs)
	fmt.Printf("Serving html pages on port %s...\n", port)
	http.ListenAndServe(port, nil)
}

func TestCrawligc(t *testing.T) {
	go srvStatic()
	crawligc.CrawlIGC("http://localhost"+port+"/", "0.html")
}
