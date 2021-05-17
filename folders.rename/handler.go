package main

import (
	"net/url"
	"strings"
	"thin-peak/logs/logger"

	"github.com/big-larry/mgo"
	"github.com/big-larry/mgo/bson"
	"github.com/big-larry/suckhttp"
)

type RenameFolder struct {
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

func NewRenameFolder(mgoAddr string, mgoColl string) (*RenameFolder, error) {

	mgoSession, err := mgo.Dial(mgoAddr)
	if err != nil {
		logger.Error("Mongo conn", err)
		return nil, err
	}

	mgoCollection := mgoSession.DB("main").C(mgoColl)

	return &RenameFolder{mgoSession: mgoSession, mgoColl: mgoCollection}, nil

}

func (conf *RenameFolder) Close() error {
	conf.mgoSession.Close()
	return nil
}

func (conf *RenameFolder) Handle(r *suckhttp.Request, l *logger.Logger) (*suckhttp.Response, error) {

	cookie, ok := r.GetCookie("koki")
	if cookie == "" || !ok { // TODO: нужна ли проверка на "" ?
		return suckhttp.NewResponse(401, "Unauthorized"), nil
	}

	// TODO: AUTH

	if !strings.Contains(r.GetHeader(suckhttp.Content_Type), "application/x-www-form-urlencoded") {
		return suckhttp.NewResponse(400, "Bad request"), nil
	}

	formValues, err := url.ParseQuery(string(r.Body))
	if err != nil {
		return suckhttp.NewResponse(400, "Bad Request"), err
	}

	fid := formValues.Get("fid")
	fnewname := formValues.Get("fnewname")
	if fid == "" || fnewname == "" {
		return suckhttp.NewResponse(400, "Bad request"), nil
	}
	// TODO: get metauser
	metaid := "randmetaid"
	//

	query := &bson.M{"_id": fid, "deleted": bson.M{"$exists": false}, "$or": []bson.M{{"metas": &meta{Type: 0, Id: metaid}}, {"metas": &meta{Type: 1, Id: metaid}}}}

	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"name": fnewname}},
		Upsert:    false,
		ReturnNew: true,
		Remove:    false,
	}

	var foo interface{}

	_, err = conf.mgoColl.Find(query).Apply(change, &foo)
	if err != nil {
		if err == mgo.ErrNotFound {
			return suckhttp.NewResponse(403, "Forbidden"), nil
		}
		return nil, err
	}

	// if info.Updated != 1 {
	// 	return suckhttp.NewResponse(403, "Forbidden"), nil
	// }

	return suckhttp.NewResponse(200, "OK"), nil
}
