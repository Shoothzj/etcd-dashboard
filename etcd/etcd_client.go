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
	"context"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
	v3 "go.etcd.io/etcd/client/v3"
	"time"
)

func (h *Handler) PutKey(ctx context.Context, key string, val string) error {
	_, err := h.client.Put(ctx, key, val)
	if err != nil {
		logrus.Errorf("Error put key. %s", err)
		return err
	}
	return nil
}

func (h *Handler) DeleteKey(ctx context.Context, key string) error {
	_, err := h.client.Delete(ctx, key)
	if err != nil {
		logrus.Errorf("Error delete key. %s", err)
		return err
	}
	return nil
}

func (h *Handler) KeysList(ctx context.Context) []string {
	logrus.Info("begin to list keys")
	clientCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	getResponse, err := h.client.Get(clientCtx, "", v3.WithPrefix(), v3.WithKeysOnly(), v3.WithSerializable())
	if err != nil {
		logrus.Errorf("Error list keys. %s", err)
		return nil
	}
	keys := make([]string, getResponse.Count)
	for i, kv := range getResponse.Kvs {
		keys[i] = string(kv.Key)
	}
	return keys
}

func (h *Handler) GetKeyContent(ctx context.Context, key string) ([]byte, error) {
	clientCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	resp, err := h.client.Get(clientCtx, key)
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, nil
	}
	return resp.Kvs[0].Value, nil
}

func (h *Handler) GetAllKeyContent(ctx context.Context) ([]*mvccpb.KeyValue, error) {
	clientCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	resp, err := h.client.Get(clientCtx, "", v3.WithPrefix(), v3.WithSerializable())
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, nil
	}
	return resp.Kvs, nil
}
