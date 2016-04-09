package imganalyze

import (
	"fmt"
	"testing"
)

func TestPHcolor(t *testing.T) {
	ph := TestPH(DecodeImg("1.jpg"))
	fmt.Println(ph)
}
