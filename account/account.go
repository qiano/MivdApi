package account

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	UserName string
	Password string
	Role     string
}

var dbName = "mivd"
var collectionName = "account"

func Add(session *mgo.Session, username, password, role string) (code int, msg string) {
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

func Login(session *mgo.Session, username, password string) (ac Account) {
	session.DB(dbName).C(collectionName).Find(bson.M{"username": username, "password": password}).One(&ac)
	return
}
