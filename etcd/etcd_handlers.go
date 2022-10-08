package etcd

import (
	"etcd-dashboard/util"
	"fmt"
	"github.com/gorilla/mux"
	v3 "go.etcd.io/etcd/client/v3"
	"time"
)

type Handler struct {
	client *v3.Client
}

func NewHandler(host string, port int, tlsEnabled bool, certFile string, keyFile string, caFile string) (*Handler, error) {
	var client *v3.Client
	var err error
	if tlsEnabled {
		tlsConfig, err := util.NewTLSConfig(certFile, keyFile, caFile)
		if err != nil {
			return nil, err
		}
		client, err = v3.New(v3.Config{
			Endpoints:   []string{fmt.Sprintf("%s:%d", host, port)},
			DialTimeout: 5 * time.Second,
			TLS:         tlsConfig,
		})
	} else {
		client, err = v3.New(v3.Config{
			Endpoints:   []string{fmt.Sprintf("%s:%d", host, port)},
			DialTimeout: 5 * time.Second,
		})
	}
	if err != nil {
		return nil, err
	}
	return &Handler{
		client: client,
	}, nil
}

func (h *Handler) Handle(subRouter *mux.Router) {
	subRouter.HandleFunc("/keys", h.keysHandler)
	subRouter.HandleFunc("/keys/{key}", h.keyHandler)
	subRouter.HandleFunc("/keys/{key}/decode", h.keyHandler)
}
