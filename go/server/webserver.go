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
	"time"

	"github.com/jaeyeom/sstable/go/sstable"
)

func indexHandler(tbl *sstable.SSTable) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		var from []byte
		if v, ok := r.Form["from"]; ok {
			from = []byte(v[0])
		}

		var to []byte
		if v, ok := r.Form["to"]; ok {
			to = []byte(v[0])
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, "<!DOCTYPE html>\n<html><body><ol>\n")

		for c := tbl.ScanFrom(from); !c.Done(); c.Next() {
			e := c.Entry()
			if to != nil && bytes.Compare(to, e.Key) < 0 {
				break
			}

			href := "/lookup?key=" + url.QueryEscape(string(e.Key))
			fmt.Fprintf(w, "<li><a href=\"%s\">%s</a></li>", href, string(e.Key))
		}
		fmt.Fprint(w, "</ol></body></html>\n")
	}
}

func lookupHandler(tbl *sstable.SSTable) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
		if _, err := w.Write(e.Value); err != nil {
			log.Printf("Error writing the value to response writer: %+v", err)
		}
	}
}

func main() {
	var (
		path = flag.String("path", "", "path of SSTable filename")
		addr = flag.String("addr", ":9001", "address of the web server")
	)

	flag.Parse()

	if *path == "" {
		log.Fatal("Please specify --path flag.")
	}

	f, err := os.Open(*path)
	if err != nil {
		log.Fatal("Open file failed:", err)
	}

	tbl, err := sstable.NewSSTable(f)
	if err != nil {
		log.Fatal("SSTable creation failed:", err)
	}

	http.HandleFunc("/", indexHandler(tbl))
	http.HandleFunc("/lookup", lookupHandler(tbl))

	server := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 10 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe failed: ", err)
	}
}
