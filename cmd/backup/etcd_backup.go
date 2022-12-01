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
	"etcd-dashboard/etcd"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/url"
	"os"
)

var (
	action         string
	etcdHost       string
	etcdPort       int
	etcdCert       string
	etcdTlsEnabled bool
	etcdKey        string
	etcdCA         string
	backupPath     string
	rootCmd        = &cobra.Command{
		Use:   "etcd-backup",
		Short: "etcd-backup is a tool for backing up etcd data",
		Run: func(cmd *cobra.Command, args []string) {
			switch action {
			case "backup":
				backup()
			case "restore":
				restore()
			default:
				logrus.Errorf("Unknown action: %v", action)
				os.Exit(1)
			}
		},
	}
)

func backup() {
	handler, err := etcd.NewHandler(etcdHost, etcdPort, etcdTlsEnabled, etcdCert, etcdKey, etcdCA)
	if err != nil {
		logrus.Errorf("create client failed: %v", err)
		os.Exit(1)
	}
	allKeyContent, err := handler.GetAllKeyContent(context.TODO())
	if err != nil {
		logrus.Errorf("get all key content failed: %v", err)
		os.Exit(1)
	}

	if err := os.Mkdir(backupPath, os.ModePerm); err != nil {
		logrus.Errorf("mkdir %s failed: %v", backupPath, err)
		os.Exit(1)
	}

	for _, kv := range allKeyContent {
		humanKey := url.QueryEscape(string(kv.Key))
		name := fmt.Sprintf("%s%s%s", backupPath, string(os.PathSeparator), humanKey)
		file, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
		if err != nil {
			logrus.Errorf("create file(%s) failed: %v", name, err)
		} else {
			_, err = file.Write(kv.Value)
			if err != nil {
				logrus.Errorf("writing file(%s) failed: %v", name, err)
			}
		}
		err = file.Close()
		if err != nil {
			logrus.Errorf("close fd(%s) failed: %v", name, err)
		}
	}
}

func restore() {
	handler, err := etcd.NewHandler(etcdHost, etcdPort, etcdTlsEnabled, etcdCert, etcdKey, etcdCA)
	if err != nil {
		logrus.Errorf("create client failed: %v", err)
		os.Exit(1)
	}

	filePath, err := os.Open(backupPath)
	defer func(filePath *os.File) {
		err := filePath.Close()
		if err != nil {
			logrus.Errorf("close file failed: %v", err)
		}
	}(filePath)
	if err != nil {
		logrus.Errorf("not fund path[%s], get error: %v", backupPath, err)
		os.Exit(1)
	}
	files, err := filePath.ReadDir(-1)
	if err != nil {
		logrus.Errorf("get directory(%s) files failed: %v", backupPath, err)
		os.Exit(1)
	}

	for _, file := range files {
		if file.IsDir() {
			logrus.Warnf("get directory: %s, skip", file.Name())
			continue
		}
		fileName := fmt.Sprintf("%s%s%s", backupPath, string(os.PathSeparator), file.Name())
		bytes, err := os.ReadFile(fileName)
		if err != nil {
			logrus.Errorf("reading file(%s) failed: %v", fileName, err)
		} else {
			humanKey := url.QueryEscape(file.Name())
			err := handler.PutKey(context.TODO(), humanKey, string(bytes))
			if err != nil {
				logrus.Errorf("put key: %s, value: %s to etcd failed: %v", humanKey, string(bytes), err)
			}
		}
	}
}

func main() {
	rootCmd.PersistentFlags().StringVar(&action, "action", "", "Action to perform. Possible values: backup, restore")
	rootCmd.PersistentFlags().StringVar(&etcdHost, "host", "localhost", "etcd host")
	rootCmd.PersistentFlags().IntVar(&etcdPort, "port", 2379, "etcd port")
	rootCmd.PersistentFlags().BoolVar(&etcdTlsEnabled, "tls-enabled", false, "Enable TLS")
	rootCmd.PersistentFlags().StringVar(&etcdCert, "cert", "", "etcd cert")
	rootCmd.PersistentFlags().StringVar(&etcdKey, "key", "", "etcd key")
	rootCmd.PersistentFlags().StringVar(&etcdCA, "ca", "", "etcd ca")
	rootCmd.PersistentFlags().StringVar(&backupPath, "backup-path", "", "Path to backup file")
	if err := rootCmd.Execute(); err != nil {
		logrus.Errorf("Error executing etcd-backup: %v", err)
		os.Exit(1)
	}
}
