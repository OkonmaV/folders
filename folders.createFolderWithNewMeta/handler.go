package main

import (
	"thin-peak/httpservice"
	"thin-peak/logs/logger"
	"time"

	"github.com/big-larry/suckhttp"
	"github.com/big-larry/suckutils"
)

type CreateFolderWithNewMeta struct {
	createFolder      *httpservice.InnerService
	createMetauser    *httpservice.InnerService
	changeFoldersMeta *httpservice.InnerService
}

func NewCreateFolderWithNewMeta(createFolder *httpservice.InnerService, createMetauser *httpservice.InnerService, changeFoldersMeta *httpservice.InnerService) (*CreateFolderWithNewMeta, error) {
	return &CreateFolderWithNewMeta{createFolder: createFolder, createMetauser: createMetauser, changeFoldersMeta: changeFoldersMeta}, nil
}

func (conf *CreateFolderWithNewMeta) Handle(r *suckhttp.Request, l *logger.Logger) (*suckhttp.Response, error) {

	createFolderResp, err := conf.createFolder.Send(r)
	if err != nil {
		return nil, err
	}
	if i, _ := createFolderResp.GetStatus(); i != 200 {
		return nil, nil
	}
	r.
	createMetauserResp, err := conf.createMetauser.Send(r)

	expires := time.Now().Add(20 * time.Hour).String()
	resp := suckhttp.NewResponse(200, "OK")
	resp.SetHeader(suckhttp.Set_Cookie, suckutils.ConcatFour("koki=", string(tokenResp.GetBody()), ";Expires=", expires))

	return resp, nil
}
