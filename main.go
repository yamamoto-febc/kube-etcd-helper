package main

import (
	"k8s.io/kubectl/pkg/scheme"
	// install all APIs
	"log"
	"os"

	"kube-etcd-helper/command"

	cli "github.com/urfave/cli/v2"
	apiregistration "k8s.io/kube-aggregator/pkg/apis/apiregistration/install"
)

var (
	appName      = "kube-etcd-helper"
	appUsage     = "for tracking etcd(v3) events in kubernetes"
	appCopyright = "Copyright (C) 2018 Kazumichi Yamamoto."
)

func init() {
	apiregistration.Install(scheme.Scheme)
}

func main() {
	app := &cli.App{
		Name:      appName,
		Usage:     appUsage,
		HelpName:  appName,
		Copyright: appCopyright,
		Version:   command.Version,
		Flags:     command.CliFlags,
		Commands:  command.Commands,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
