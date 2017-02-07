package ph

import (
	"fmt"
	"github.com/qshuai162/MivdApi_Trail/imganalyze"
	"testing"	

)

func TestPHcolor(t *testing.T) {
	ph := TestPH(imganalyze.DecodeImg("1.jpg"))
	fmt.Println(ph)
}
