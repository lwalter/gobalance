package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func rootHandler(name string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract headers for logging
		var headers []string
		for n, hs := range r.Header {
			n = strings.ToLower(n)
			for _, h := range hs {
				headers = append(headers, fmt.Sprintf("%v: %v", n, h))
			}
		}

		// Extract body for logging
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println(fmt.Sprintf("Could not parse request body: %s", err))
		}

		fmt.Println(fmt.Sprintf(
			"[%s] method: %s route: %s remote: %s referer: %s user-agent: %s\nheaders: %s\nbody: %s",
			name,
			r.Method,
			r.URL.String(),
			r.RemoteAddr,
			r.Referer(),
			r.UserAgent(),
			headers,
			body,
		))

		w.Header().Set("Content-Type", "application-json")
		w.Header().Set("Status", "200")
		json.NewEncoder(w).Encode(string(body))
	}
}

func initRouters(count int) []*mux.Router {
	var routers []*mux.Router
	// TODO(lnw) catch all instead of single
	route := "/"
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	for i := 0; i < count; i++ {
		router := mux.NewRouter()
		name := fmt.Sprintf("server-%d", i)
		catchAll := router.PathPrefix("/")
		catchAll.Handler(http.HandlerFunc(rootHandler(name)))
		catchAll.Methods(
			"GET",
			"POST",
			"PUT",
			"DELETE")
		fmt.Println(fmt.Sprintf("Registered path %s supprting methods %v for %s", route, methods, name))
		routers = append(routers, router)
	}

	return routers
}

func runServers(routers []*mux.Router, ports []int) error {
	rCount := len(routers)
	if len(ports) != rCount {
		return errors.New("Unequal routers to ports")
	}

	fmt.Println(fmt.Sprintf("Starting servers on ports %v", ports))

	for i := 0; i < rCount-1; i++ {
		router := routers[i]
		p := ports[i]
		go func(p int, r *mux.Router) {
			pStr := strconv.Itoa(p)
			if err := http.ListenAndServe(":"+pStr, r); err != nil {
				log.Fatal(err)
			}
		}(p, router)
	}

	lastRouterIdx := rCount - 1
	p := strconv.Itoa(ports[lastRouterIdx])
	if err := http.ListenAndServe(":"+p, routers[lastRouterIdx]); err != nil {
		log.Fatal(err)
	}

	return nil
}

func main() {
	// TODO(lnw) pass in via command line arg
	ports := []int{8080, 8081}
	count := len(ports)
	routers := initRouters(count)
	runServers(routers, ports)
}
