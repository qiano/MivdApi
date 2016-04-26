package record

import (
	. "github.com/qshuai162/common/config"
	// "fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

type Record struct {
	Id        int    //编号
	Type      string //类型
	Date      int64  //时间
	PicData   string //照片base64
	AreaData  string //反应区图片base64
	Project   string //测试项目
	ResultMsg string //判定
	LotNo     string //批号
	UserName  string //操作员
	Location  string //地点
	Remark    string //备注
}

var mongodbstr = Config["mongodbHost"]
var dbName = Config["mongodbDbName"]

// var mongodbstr = "121.41.46.25:27017"
// var dbName = "mivd_dev"
var collectionName = "record"

func (r *Record) Save() {
	session, err := mgo.Dial(mongodbstr)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	col := session.DB(dbName).C(collectionName)
	var temp Record
	col.Find(nil).Sort("-$natural").One(&temp)
	r.Id = temp.Id + 1
	col.Insert(r)
}

func GetList(pageIdx, pageSize int, user, role,ty string) (records []Record) {
	session, err := mgo.Dial(mongodbstr)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	col := session.DB(dbName).C(collectionName)
	if role != "user" {
		col.Find(bson.M{"type":ty}).Sort("-$natural").Skip((pageIdx - 1) * pageSize).Limit(pageSize).All(&records)
	} else {
		col.Find(bson.M{"username": user,"type":ty}).Sort("-$natural").Skip((pageIdx - 1) * pageSize).Limit(pageSize).All(&records)
	}

	return
}

func FindById(id string) (record Record) {
	session, err := mgo.Dial(mongodbstr)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()
	ID, _ := strconv.Atoi(id)
	session.DB(dbName).C(collectionName).Find(bson.M{"id": ID}).One(&record)
	return
}

func FindByOid(oid string) (record Record) {
	session, err := mgo.Dial(mongodbstr)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()
	session.DB(dbName).C(collectionName).FindId(oid)
	return
}
