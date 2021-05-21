package main

import (
	"net/url"
	"strings"
	"thin-peak/logs/logger"

	"github.com/big-larry/mgo"
	"github.com/big-larry/mgo/bson"
	"github.com/big-larry/suckhttp"
)

type MetaChange struct {
	mgoSession *mgo.Session
	mgoColl    *mgo.Collection
}
type folder struct {
	Id    string   `bson:"_id"`
	Roots []string `bson:"users"`
	Name  string   `bson:"name"`
	Metas []meta   `bson:"type"`
}

type meta struct {
	Type int    `bson:"metatype"`
	Id   string `bson:"metaid"`
}

func NewMetaChange(mgoAddr string, mgoColl string) (*MetaChange, error) {

	mgoSession, err := mgo.Dial(mgoAddr)
	if err != nil {
		logger.Error("Mongo conn", err)
		return nil, err
	}

	mgoCollection := mgoSession.DB("main").C(mgoColl)

	return &MetaChange{mgoSession: mgoSession, mgoColl: mgoCollection}, nil

}

func (conf *MetaChange) Close() error {
	conf.mgoSession.Close()
	return nil
}

func (conf *MetaChange) Handle(r *suckhttp.Request, l *logger.Logger) (*suckhttp.Response, error) {

	// TODO: AUTH

	if !strings.Contains(r.GetHeader(suckhttp.Content_Type), "application/x-www-form-urlencoded") {
		return suckhttp.NewResponse(400, "Bad request"), nil
	}

	formValues, err := url.ParseQuery(string(r.Body))
	if err != nil {
		return suckhttp.NewResponse(400, "Bad Request"), err
	}

	fid := formValues.Get("fid")
	fnewmeta := formValues.Get("fnewmeta")
	if fid == "" || fnewmeta == "" {
		return suckhttp.NewResponse(400, "Bad request"), nil
	}

	query := &bson.M{"_id": fid, "deleted": bson.M{"$exists": false}}

	change := mgo.Change{
		Update:    bson.M{"$addToSet": bson.M{"metas": &meta{Type: 1, Id: fid}}},
		Upsert:    false,
		ReturnNew: true,
		Remove:    false,
	}
	var foo interface{}

	_, err = conf.mgoColl.Find(query).Apply(change, &foo)
	if err != nil {
		if err == mgo.ErrNotFound {
			return suckhttp.NewResponse(400, "Bad request"), nil
		}
		return nil, err
	}

	return suckhttp.NewResponse(200, "OK"), nil
}
