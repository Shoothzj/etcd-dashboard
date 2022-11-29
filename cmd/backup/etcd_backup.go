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
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
}

func restore() {
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
