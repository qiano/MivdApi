package account

import (
	. "github.com/qshuai162/MivdApi/common/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	UserName string
	Password string
	Role     string
}

var mongodbstr = Config["mongodbHost"]
var dbName = Config["mongodbDbName"]

// var mongodbstr = "121.41.46.25:27017"
// var dbName = "mivd_dev"
var collectionName = "account"

func Add(username, password, role string) (code int, msg string) {
	session, err := mgo.Dial(mongodbstr)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	col := session.DB(dbName).C(collectionName)
	count, err := col.Find(bson.M{"username": username}).Count()
	if err != nil {
		return 2, err.Error()
	}
	if count > 0 {
		return 1, "username is exist:" + username
	} else {
		col.Insert(&Account{username, password, role})
		return 0, "add user success"
	}
}

func Login(username, password string) (ac Account) {
	session, err := mgo.Dial(mongodbstr)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()
	session.DB(dbName).C(collectionName).Find(bson.M{"username": username, "password": password}).One(&ac)
	return
}
