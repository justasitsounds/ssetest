package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/justasitsounds/ssetest/sse"
)

func main() {
	l := log.New(os.Stderr, "", log.LUTC)
	b := sse.NewBroker(*l)

	go func() {
		for {
			msg := time.Now().String()
			e := sse.NewEvent([]byte(msg), 0)
			b.Publish(e)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	http.Handle("/sse", b)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<html><head><style type=\"text/css\">%s</style></head><body><script type=\"text/javascript\">(%s)()</script></body></html>", clientSideStyles(), clientSideCode())
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func clientSideCode() string {
	return `function js() {
		var es = new EventSource('/sse')
		  , pre = document.createElement('p')
		  , closed = false
	  
		document.body.appendChild(pre)
	  
		es.onmessage = function(ev) {
		  if(closed) return
	  
		  pre.appendChild(document.createTextNode(ev.data))
	  
		  window.scrollTo(0, pre.clientHeight)
		}
	  
		es.addEventListener('end', function() {
		  es.close()
		  closed = true
		}, true)
	  
		es.onerror = function(e) {
		  closed = true
		}
	  }`
}

func clientSideStyles() string {
	return `
	`
}
