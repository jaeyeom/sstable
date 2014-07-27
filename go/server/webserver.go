// Binary webserver implements an HTTP server for reading an SSTable.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/jaeyeom/sstable/go/sstable"
)

var (
	path = flag.String("path", "", "path of SSTable filename")
	addr = flag.String("addr", ":9001", "address of the web server")
	tbl  *sstable.SSTable
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// This check prevents the "/" handler from handling all
	// requests by default
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.NotFound(w, r)
		return
	}
	var from, to []byte
	toLast := true
	if v, ok := r.Form["from"]; ok {
		from = []byte(v[0])
	}
	if v, ok := r.Form["to"]; ok {
		to = []byte(v[0])
		toLast = false
	}
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<!DOCTYPE html>\n<html><body><ol>\n")
	for c := tbl.ScanFrom(from); !c.Done(); c.Next() {
		e := c.Entry()
		if !toLast && bytes.Compare(to, e.Key) < 0 {
			break
		}
		href := "/lookup?key=" + url.QueryEscape(string(e.Key))
		fmt.Fprintf(w, "<li><a href=\"%s\">%s</a></li>", href, string(e.Key))
	}
	fmt.Fprint(w, "</ol></body></html>\n")
}

func lookupHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/lookup" {
		http.NotFound(w, r)
		return
	}
	key := []byte(r.FormValue("key"))
	c := tbl.ScanFrom(key)
	if c.Done() {
		http.NotFound(w, r)
		return
	}
	e := c.Entry()
	if !bytes.Equal(key, e.Key) {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", http.DetectContentType(e.Value))
	w.Write(e.Value)
}

func main() {
	flag.Parse()
	if *path == "" {
		log.Fatal("Please specify --path flag.")
	}
	f, err := os.Open(*path)
	if err != nil {
		log.Fatal("Open file failed:", err)
	}
	tbl, err = sstable.NewSSTable(f)
	if err != nil {
		log.Fatal("SSTable creation failed:", err)
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/lookup", lookupHandler)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe failed: ", err)
	}
}
