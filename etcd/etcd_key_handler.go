package etcd

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	v3 "go.etcd.io/etcd/client/v3"
	"net/http"
	"time"
)

func (h *Handler) keysHandler(w http.ResponseWriter, r *http.Request) {
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
	Content []byte `json:"content"`
}

func (h *Handler) keyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	logrus.Infof("begin to get key %s", key)
	content, err := h.GetKeyContent(r.Context(), key)
	if err != nil {
		logrus.Errorf("get key %s content failed, err: %v", key, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := GetKeyResp{
		Content: content,
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
