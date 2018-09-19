package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"strings"
)

func write(w http.ResponseWriter, fmt_ string, args ...interface{}) {
	w.Write([]byte(fmt.Sprintf(fmt_, args...)))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Request from %s\n", r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
	write(w, "Hello World from %s\n", r.Host)
	write(w, "Path: %s\n", r.URL.Path)
	write(w, "\nHeaders:\n")
	r.Header.WriteSubset(w, nil)

	addr, err := net.InterfaceAddrs()
	if err != nil {
		write(w, "Cannot read interfaces\n")
	} else {
		write(w, "\nNetwork interfaces:\n")
		for _, a := range addr {
			if !strings.Contains(a.String(), "::") {
				write(w, "%s\n", a.String())
			}
		}
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Health check performed (%s)\n", r.URL.Path)
	w.WriteHeader(http.StatusOK)
}

func main() {
	fmt.Println("Starting app on port 9000")

	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(mainHandler)
	r.HandleFunc("/health", healthHandler)
	r.HandleFunc("/index", mainHandler)

	err := http.ListenAndServe("0.0.0.0:9000", r)
	if err != nil {
		fmt.Printf("Cannot start app: %s\n", err)
	}
}
