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
