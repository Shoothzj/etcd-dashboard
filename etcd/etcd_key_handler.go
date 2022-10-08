package etcd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	v3 "go.etcd.io/etcd/client/v3"
	"net/http"
	"time"
)

type PutKeyReq struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (h *Handler) keyPutHandler(w http.ResponseWriter, r *http.Request) {
	var req PutKeyReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logrus.Infof("begin to put key %s", req.Key)
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	_, err = h.client.Put(ctx, req.Key, req.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) keysListHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("begin to list keys")
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	getResponse, err := h.client.Get(ctx, "", v3.WithPrefix(), v3.WithKeysOnly(), v3.WithSerializable())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	keys := make([]string, getResponse.Count)
	for i, kv := range getResponse.Kvs {
		keys[i] = string(kv.Key)
	}
	payload, err := json.Marshal(keys)
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

type GetKeyResp struct {
	Content string `json:"content"`
}

func (h *Handler) keyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	decodeKeyByte, err := base64.StdEncoding.DecodeString(key)
	decodeKey := string(decodeKeyByte)
	logrus.Infof("begin to get key %s", decodeKey)
	content, err := h.GetKeyContent(r.Context(), decodeKey)
	if err != nil {
		logrus.Errorf("get key %s content failed, err: %v", decodeKey, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := GetKeyResp{
		Content: string(content),
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
