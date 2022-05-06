// Copyright (c) 2022 Red Hat, Inc.

package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/component-base/logs"

	"github.com/clyang82/hub-of-hubs-apis/pkg/server"
	"github.com/spf13/pflag"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	logs.InitLogs()
	defer logs.FlushLogs()

	opts := server.NewOptions()
	opts.AddFlags(pflag.CommandLine)

	clusterCfg, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	dynamicClient, err := dynamic.NewForConfig(clusterCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	dynInformerFactory := dynamicinformer.NewDynamicSharedInformerFactory(dynamicClient, 0)

	apiServerConfig, err := opts.APIServerConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	proxyServer, err := server.NewServer(dynamicClient, nil, apiServerConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	stopCh := make(chan struct{})

	dynInformerFactory.Start(stopCh)

	if err := proxyServer.Run(stopCh); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
