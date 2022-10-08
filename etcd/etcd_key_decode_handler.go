package etcd

import (
	"encoding/json"
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
	logrus.Infof("begin to get key %s component: %s, namespace: %s", key, decodeComponent, decodeNamespace)
	content, err := h.GetKeyContent(r.Context(), key)
	if err != nil {
		logrus.Errorf("get key %s content failed, err: %v", key, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	str, err := decode(decodeComponent, decodeNamespace, content)
	if err != nil {
		logrus.Errorf("decode key %s content failed component: %s, namespace: %s, err: %v", key, decodeComponent, decodeNamespace, err)
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
