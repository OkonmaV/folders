package main

import (
	"net/url"
	"strings"
	"thin-peak/logs/logger"

	"github.com/big-larry/mgo"
	"github.com/big-larry/suckhttp"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
)

type CreateFolder struct {
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

func NewCreateFolder(mgoAddr string, mgoColl string) (*CreateFolder, error) {

	mgoSession, err := mgo.Dial(mgoAddr)
	if err != nil {
		logger.Error("Mongo conn", err)
		return nil, err
	}

	mgoCollection := mgoSession.DB("main").C(mgoColl)

	return &CreateFolder{mgoSession: mgoSession, mgoColl: mgoCollection}, nil

}

func (conf *CreateFolder) Close() error {
	conf.mgoSession.Close()
	return nil
}

func getRandId() string {
	return xid.New().String()
}

func (conf *CreateFolder) Handle(r *suckhttp.Request, l *logger.Logger) (*suckhttp.Response, error) {

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

	froot := formValues.Get("froot")
	fname := formValues.Get("fname")
	if froot == "" || fname == "" {
		return suckhttp.NewResponse(400, "Bad request"), nil
	}
	// TODO: get metauser
	metaid := "randmetaid"
	//

	// check root meta ?????
	selector := &bson.M{"_id": froot, "deleted": bson.M{"$exists": false}, "$or": []bson.M{{"metas": &meta{Type: 0, Id: metaid}}, {"metas": &meta{Type: 1, Id: metaid}}}}
	var foo interface{}

	err = conf.mgoColl.Find(selector).One(&foo)
	if err != nil {
		if err == mgo.ErrNotFound {
			return suckhttp.NewResponse(403, "Forbidden"), nil
		}
		return nil, err
	}
	//

	finsert := &folder{Id: getRandId(), Roots: []string{froot}, Name: fname, Metas: []meta{{Type: 0, Id: metaid}}}

	err = conf.mgoColl.Insert(finsert)
	if err != nil {
		return nil, err
	}

	return suckhttp.NewResponse(200, "OK"), nil
}
