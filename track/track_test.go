package track

import (
	"testing"
    "fmt"
)

func Test_Add(t *testing.T) {
    re:=NewTrackRecord("11111111111","test","test","fdfdfddf",123.123,3213.321)
    
    
	// re.DateTime = time.Now().Unix()
	re.Save()

    list:=GetList(1,10,"test","user")
	fmt.Println(list)
    one:=FindByID("2")
    fmt.Println(one)
}
