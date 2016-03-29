package record

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Record struct {
	Id        string //编号
	Time      int64  //时间
	PicData   string //照片base64
	AreaData  string //反应区图片base64
	Test      string //测试项目
	Judgement string //判定
	LotNo     string //批号
	Operator  string //操作员
	Location  string //地点
	Remark    string //备注
}

var dbName = "mivd"
var collectionName = "record"

func (r *Record) Save(session *mgo.Session) {
	session.DB(dbName).C(collectionName).Insert(r)
}

func GetList(session *mgo.Session, pageIdx, pageSize int) (records []Record) {
	session.DB(dbName).C(collectionName).Find(nil).Skip((pageIdx - 1) * pageSize).Limit(pageSize).All(records)
	return
}

func FindById(session *mgo.Session, id string) (record Record) {
	session.DB(dbName).C(collectionName).Find(bson.M{"Id": id}).One(record)
	return
}

func FindByOid(session *mgo.Session, oid string) (record Record) {
	session.DB(dbName).C(collectionName).FindId(oid)
	return
}
