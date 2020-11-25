package main

import (
	"flag"
	"net/http"
	"time"
)

func main() {
	bind := flag.String("bind", "127.0.0.1:9898", "`addr:port` to listen on")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		delay := time.Second
		if v := r.FormValue("t"); v != "" {
			var err error
			if delay, err = time.ParseDuration(v); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		go func() {
			time.Sleep(delay)
			panic("panic! at the disco")
		}()

		w.WriteHeader(http.StatusAccepted)
	})

	http.ListenAndServe(*bind, nil)
}
