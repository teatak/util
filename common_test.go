package util

import (
	"fmt"
	"testing"

	"github.com/teatak/util/common"
)

func TestRandom(t *testing.T) {
	ok := common.ValidateMobile("+1:3465759126")
	fmt.Println(ok)

	ok = common.ValidateMobile("+86:18611822358")
	fmt.Println(ok)
}
