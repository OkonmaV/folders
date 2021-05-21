package main

import (
	"context"
	"lib"
	"thin-peak/httpservice"
)

type config struct {
	Configurator string
	Listen       string
}

var thisServiceName httpservice.ServiceName = "conf.createfolderwithnewmeta"
var createFolderServiceName httpservice.ServiceName = "conf.createfolder"
var createMetauserServiceName httpservice.ServiceName = "conf.createmetauser"
var changeFoldersMetaServiceName httpservice.ServiceName = "conf.changefoldersmeta"

func (c *config) GetListenAddress() string {
	return c.Listen
}
func (c *config) GetConfiguratorAddress() string {
	return c.Configurator
}
func (c *config) CreateHandler(ctx context.Context, connectors map[httpservice.ServiceName]*httpservice.InnerService) (httpservice.HttpService, error) {

	return NewCreateFolderWithNewMeta(connectors[createFolderServiceName], connectors[createMetauserServiceName], connectors[changeFoldersMetaServiceName])
}

func main() {
	httpservice.InitNewService(thisServiceName, false, 5, &config{}, lib.ServiceNameCookieTokenGen)
}
