package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

type BackendServer struct {
	ID   string
	PORT string
	IP   string
	Load int
}

type LoadBalancer struct {
	Servers         *[]BackendServer
	PORT            string
	IP              string
	Algorithm       string
	HealthCheckPath string
	NextIndex       int
	mu              sync.Mutex
}

func NewLoadBalancer(servers *[]BackendServer, port string, ip string, algorithm string, healthCheckPath string) *LoadBalancer {
	return &LoadBalancer{
		Servers:         servers,
		PORT:            port,
		IP:              ip,
		Algorithm:       algorithm,
		HealthCheckPath: healthCheckPath,
		NextIndex:       0,
		mu:              sync.Mutex{},
	}
}

func backendServerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Response from Backend Server %s\n", r.URL.Path)
}

func (lb *LoadBalancer) handleHealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	for _, server := range *lb.Servers {
		fmt.Fprintf(w, "{\"ID\":\"%s\",\"IP\":\"%s\",\"PORT\":\"%s\",\"Load\":%d}\n", server.ID, server.IP, server.PORT, server.Load)
	}
}

func (lb *LoadBalancer) getNextServer() *BackendServer {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	var server *BackendServer

	switch lb.Algorithm {
	case "round-robin":
		server = &(*lb.Servers)[lb.NextIndex]
		lb.NextIndex = (lb.NextIndex + 1) % len(*lb.Servers)
	default:
		server = &(*lb.Servers)[0]
	}

	server.Load += 1
	return server
}

func (lb *LoadBalancer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	// if health check request
	print(r.URL.Path, "\n", lb.HealthCheckPath, "\n")
	if r.URL.Path == lb.HealthCheckPath {
		lb.handleHealthCheck(w, r)
		return
	}

	server := lb.getNextServer()
	proxyURL := fmt.Sprintf("http://%s:%s/bs%s", server.IP, server.PORT, server.ID)

	resp, err := http.Get(proxyURL)

	if err != nil {
		http.Error(w, "Backend server error", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func main() {
	backendServers := []BackendServer{
		{ID: "1", PORT: "8081", IP: "127.0.0.1", Load: 0},
		{ID: "2", PORT: "8082", IP: "127.0.0.1", Load: 0},
		{ID: "3", PORT: "8083", IP: "127.0.0.1", Load: 0},
	}

	// run all backend server parallelly
	for _, server := range backendServers {
		go func(s BackendServer) {
			http.HandleFunc("/bs"+s.ID, backendServerHandler)
			http.ListenAndServe(s.IP+":"+s.PORT, nil)
		}(server)
	}

	loadBalancer := NewLoadBalancer(&backendServers, "8080", "127.0.0.1", "round-robin", "/health")

	http.HandleFunc("/", loadBalancer.HandleRequest)

	http.ListenAndServe(loadBalancer.IP+":"+loadBalancer.PORT, nil)

}
