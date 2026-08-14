package main

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	api, _ := url.Parse("http://127.0.0.1:8080")
	proxy := httputil.NewSingleHostReverseProxy(api)
	files, _ := url.Parse("http://127.0.0.1:8081")
	filesProxy := httputil.NewSingleHostReverseProxy(files)
	root := os.Getenv("LABSOS_DASHBOARD_ROOT")
	if root == "" {
		root = "/opt/labsos/dashboard"
	}
	lan := os.Getenv("LABSOS_LAN_CIDR")
	if lan == "" {
		lan = "192.168.0.0/24"
	}
	_, network, err := net.ParseCIDR(lan)
	if err != nil {
		log.Fatal(err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !allowed(r, network) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/") {
			proxy.ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/file-manager") {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, "/file-manager")
			if r.URL.Path == "" {
				r.URL.Path = "/"
			}
			filesProxy.ServeHTTP(w, r)
			return
		}
		path := filepath.Join(root, filepath.Clean("/"+r.URL.Path))
		if info, err := os.Stat(path); err != nil || info.IsDir() {
			path = filepath.Join(root, "index.html")
		}
		http.ServeFile(w, r, path)
	})
	log.Printf("LabsOS dashboard http://0.0.0.0:80 root %s lan %s", root, lan)
	log.Fatal(http.ListenAndServe(":80", handler))
}

func allowed(r *http.Request, network *net.IPNet) bool {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false
	}
	ip := net.ParseIP(host)
	return ip.IsLoopback() || network.Contains(ip)
}
