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

package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/embed"
	"testing"
	"time"
)

func Test_backup(t *testing.T) {
	cfg := embed.NewConfig()
	cfg.Dir = "default.etcd"
	e, err := embed.StartEtcd(cfg)
	assert.NoErrorf(t, err, "create etcd embed server failed: %v", err)
	defer e.Close()
	select {
	case <-e.Server.ReadyNotify():
		logrus.Info("etcd server is ready!")
	case <-time.After(60 * time.Second):
		e.Server.Stop() // trigger a shutdown
		logrus.Fatal("etcd server took too long to start!")
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	assert.NoErrorf(t, err, "connect etcd server failed: %v", err)
	defer cli.Close()

	var ctx = context.Background()

	// put key
	_, err = cli.Put(ctx, "hello_world", "great!")
	assert.NoErrorf(t, err, "put key failed: %v", err)

	// get key
	resp, err := cli.Get(ctx, "hello_world")
	assert.NoErrorf(t, err, "get key failed: %v", err)
	assert.Equal(t, 1, len(resp.Kvs))

	// check key
	kv := resp.Kvs[0]
	assert.Equal(t, "great!", string(kv.Value))
	t.Logf("key: %s, value: %s", string(kv.Key), string(kv.Value))
}
