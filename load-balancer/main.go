package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type simpleServer struct{
	addr string
	proxy *httputil.ReverseProxy
}

func newSimpleServer(addr string) *simpleServer{
	serverUrl, err := url.Parse(addr)
	handleError(err)

	return &simpleServer{
		addr: addr,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

type LoadBalancer struct{
	port 		string
	roundRobin 	int
	servers 	[]Server
}

type Server interface{
	Address () string
	IsAlive () bool
	Serve(rw http.ResponseWriter, req *http.Request)
}

func NewLoadBalancer(port string, servers []Server) *LoadBalancer{
	return &LoadBalancer{
		port: port,
		roundRobin: 0,
		servers: servers,
	}
}

func handleError(err error){
	if err!=nil{
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func (s *simpleServer) Address() string{
	return s.addr
}

func (s *simpleServer) IsAlive() bool {
	return true
}

func (s *simpleServer) Serve(rw http.ResponseWriter, req *http.Request){
	s.proxy.ServeHTTP(rw, req)
}

func (lb *LoadBalancer) getNextAvailServer() Server{
	nextServer := lb.servers[lb.roundRobin%len(lb.servers)] //get the next server in the list
	for !nextServer.IsAlive(){
		lb.roundRobin++
		nextServer = lb.servers[lb.roundRobin%len(lb.servers)]
	} //increment the round robin counter
	lb.roundRobin++
	return nextServer
}

func (lb *LoadBalancer) serveProxy(rw http.ResponseWriter, req *http.Request){
	targetServer := lb.getNextAvailServer()
	fmt.Printf("Redirecting request to %s\n", targetServer.Address())
	targetServer.Serve(rw, req)
}

func main(){
	server := []Server{
		newSimpleServer("http://www.google.com"),
		newSimpleServer("http://www.yahoo.com"),
		newSimpleServer("http://www.facebook.com"),
	}

	lb := NewLoadBalancer("8000", server)
	handleRedirect := func(rw http.ResponseWriter, req *http.Request){
		lb.serveProxy(rw, req)
	}
	http.HandleFunc("/", handleRedirect)

	fmt.Printf("Server listening on port %s\n", lb.port)
	http.ListenAndServe(":"+lb.port, nil)
}