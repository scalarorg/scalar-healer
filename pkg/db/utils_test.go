package db_test

import (
	"fmt"
	"testing"

	"github.com/scalarorg/scalar-healer/pkg/db"
)

func TestConvertUint64ToNumeric(t *testing.T) {
	var x = uint64(1234567890)
	var num = db.ConvertUint64ToNumeric(x)
	fmt.Println(num)
}
