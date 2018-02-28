package main

import (
	"k8s.io/kubernetes/pkg/kubectl/scheme"
	// install all APIs
	"github.com/yamamoto-febc/kube-etcd-helper/command"
	"gopkg.in/urfave/cli.v2"
	apiregistration "k8s.io/kube-aggregator/pkg/apis/apiregistration/install"
	_ "k8s.io/kubernetes/pkg/apis/core/install"
	"log"
	"os"
)

var (
	appName      = "kube-etcd-helper"
	appUsage     = "for tracking etcd(v3) events in kubernetes"
	appCopyright = "Copyright (C) 2018 Kazumichi Yamamoto."
)

func init() {
	apiregistration.Install(scheme.GroupFactoryRegistry, scheme.Registry, scheme.Scheme)
}

func main() {
	app := &cli.App{
		Name:                  appName,
		Usage:                 appUsage,
		HelpName:              appName,
		Copyright:             appCopyright,
		EnableShellCompletion: true,
		Version:               command.Version,
		Flags:                 command.CliFlags,
		Commands:              command.Commands,
	}
	cli.InitCompletionFlag.Hidden = true

	err := app.Run(os.Args)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
