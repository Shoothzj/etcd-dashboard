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
	"encoding/json"
	"etcd-dashboard/util"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type DecodeResult struct {
	Content string `json:"content"`
}

func (h *Handler) keyDecodeHandler(w http.ResponseWriter, r *http.Request) {
	decodeComponent := r.URL.Query().Get("decodeComponent")
	decodeNamespace := r.URL.Query().Get("decodeNamespace")
	vars := mux.Vars(r)
	key := vars["key"]
	decodeKey, err := util.Base64Decode(key)
	if err != nil {
		logrus.Errorf("base64 decode key %s failed, err: %v", key, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logrus.Infof("begin to decode key %s component: %s, namespace: %s", decodeKey, decodeComponent, decodeNamespace)
	content, err := h.GetKeyContent(r.Context(), decodeKey)
	if err != nil {
		logrus.Errorf("get key %s content failed, err: %v", key, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	str, err := decode(decodeComponent, decodeNamespace, content)
	if err != nil {
		logrus.Errorf("decode key %s content failed component: %s, namespace: %s, err: %v", decodeKey, decodeComponent, decodeNamespace, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := DecodeResult{
		Content: str,
	}
	payload, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(payload)
	if err != nil {
		logrus.Errorf("write response fail. %s", err)
	}
}
