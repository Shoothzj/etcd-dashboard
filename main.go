package main

import (
	"etcd-dashboard/etcd"
	"flag"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

var etcdHost = flag.String("etcd-host", "localhost", "etcd host")
var etcdPort = flag.Int("etcd-port", 2379, "etcd port")
var etcdTlsEnabled = flag.Bool("etcd-tls-enabled", false, "etcd tls enabled")
var etcdCertFile = flag.String("etcd-cert-file", "", "etcd cert file")
var etcdKeyFile = flag.String("etcd-key-file", "", "etcd key file")
var etcdCaFile = flag.String("etcd-ca-file", "", "etcd ca file")

func main() {
	flag.Parse()
	r := mux.NewRouter()
	r.Use(accessControlMiddleware)
	handler, err := etcd.NewHandler(*etcdHost, *etcdPort, *etcdTlsEnabled, *etcdCertFile, *etcdKeyFile, *etcdCaFile)
	if err != nil {
		panic(err)
	}
	handler.Handle(r.PathPrefix("/api/etcd").Subrouter())
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:10001",
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}
	err = srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func accessControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}
