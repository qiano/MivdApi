package account

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"testing"
)

var mongodbstr string = "121.41.46.25:27017"

func Test_Add(t *testing.T) {
	session, err := mgo.Dial(mongodbstr)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	fmt.Println(Add(session, "admin", "123456", "admin"))
	fmt.Println(Add(session, "manager", "123456", "manager"))
	fmt.Println(Add(session, "user1", "123456", "user"))
	fmt.Println(Add(session, "user2", "123456", "user"))
}

func Test_Login(t *testing.T) {
	session, err := mgo.Dial(mongodbstr)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()
	fmt.Println(Login(session, "admin", "123456"))
	fmt.Println(Login(session, "admin", "1234564"))
	fmt.Println(Login(session, "admin1", "1234564"))
}
