package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version     string
	configFile  string
	versionFlag bool
)

func main() {
	versionString := "Janus v" + version
	cobra.OnInitialize(func() {
		if versionFlag {
			fmt.Println(versionString)
			os.Exit(0)
		}
	})

	var RootCmd = &cobra.Command{
		Use:   "janus",
		Short: "Janus is an API Gateway",
		Long: versionString + `
This is a lightweight API Gateway and Management Platform that enables you
to control who accesses your API, when they access it and how they access it.
API Gateway will also record detailed analytics on how your users are interacting
with your API and when things go wrong.
Complete documentation is available at https://hellofresh.gitbooks.io/janus`,
		Run: RunServer,
	}
	RootCmd.Flags().StringVarP(&configFile, "config", "c", "", "Source of a configuration file")
	RootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Print application version")

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
