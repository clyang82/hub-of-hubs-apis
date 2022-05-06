// Copyright (c) 2022 Red Hat, Inc.

package server

import (
	"github.com/clyang82/hub-of-hubs-apis/pkg/server/api"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
)

type Server struct {
	*genericapiserver.GenericAPIServer
}

func NewServer(
	client dynamic.Interface,
	informerFactory dynamicinformer.DynamicSharedInformerFactory,
	apiServerConfig *genericapiserver.Config) (*Server, error) {
	apiServer, err := apiServerConfig.Complete(nil).New("server", genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	if err := api.Install(apiServer, client, informerFactory); err != nil {
		return nil, err
	}

	return &Server{apiServer}, nil
}

func (p *Server) Run(stopCh <-chan struct{}) error {
	return p.GenericAPIServer.PrepareRun().Run(stopCh)
}
