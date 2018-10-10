package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"os"
	"net/http"
	"strings"
	"github.com/hashicorp/consul/api"
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
	fmt.Printf("Health check performed from %s (%s)\n", r.RemoteAddr, r.URL.Path)
	w.WriteHeader(http.StatusOK)
}

func registerWithConsul() {
    addr := os.Getenv("CONSUL_SERVICE_SERVICE_HOST") + ":" + os.Getenv("CONSUL_SERVICE_SERVICE_PORT")
	fmt.Println("Registering with Consul...")
	client, err := api.NewClient(&api.Config{
		Address: addr,
	})
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	agent := client.Agent()
    addr = os.Getenv("HELLO_WORLD_SERVICE_SERVICE_HOST") + ":" + os.Getenv("HELLO_WORLD_SERVICE_SERVICE_PORT")
	err = agent.ServiceRegister(&api.AgentServiceRegistration{
		Name: "hello-world",
		Port: 9000,
		Address: "localhost",
		Tags: []string{"http", "urlprefix=/hello strip=/hello"},
		Check: &api.AgentServiceCheck{
			Name: "hello-world-check",
			HTTP: "http://" + addr + "/health",
			Timeout: "2s",
			Interval: "30s",
		},
	})
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Printf("Registered!\n")
}

func main() {
	go registerWithConsul()

	r := mux.NewRouter()
	r.HandleFunc("/health", healthHandler)
	r.HandleFunc("/", mainHandler)

	fmt.Println("Starting app on port 9000")
	err := http.ListenAndServe("0.0.0.0:9000", r)
	if err != nil {
		fmt.Printf("Cannot start app: %s\n", err)
	}
}
