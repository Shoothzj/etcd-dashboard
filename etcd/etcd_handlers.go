// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package etcd

import (
	"crypto/tls"
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
		var tlsConfig *tls.Config
		tlsConfig, err = util.NewTLSConfig(certFile, keyFile, caFile)
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
	subRouter.HandleFunc("/keys", h.keyPutHandler).Methods("PUT")
	subRouter.HandleFunc("/keys", h.keysListHandler).Methods("GET")
	subRouter.HandleFunc("/keys/{key:.*}", h.keyHandler).Methods("GET")
	subRouter.HandleFunc("/keys-decode/{key:.*}", h.keyDecodeHandler)
	subRouter.HandleFunc("/keys-delete", h.keysDeleteHandler).Methods("POST")
}
