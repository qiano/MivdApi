package account

import (
	"fmt"
	"testing"
)

func Test_Add(t *testing.T) {

	fmt.Println(Add("admin", "123456", "admin"))
	fmt.Println(Add("manager", "123456", "manager"))
	fmt.Println(Add("user1", "123456", "user"))
	fmt.Println(Add("user2", "123456", "user"))
}

func Test_Login(t *testing.T) {
	fmt.Println(Login("admin", "123456"))
	fmt.Println(Login("admin", "1234564"))
	fmt.Println(Login("admin1", "1234564"))
}
