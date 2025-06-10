package sqlc_test

import (
	"testing"

	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

func TestConvertNumericToUint64(t *testing.T) {
	n := uint64(100_000_000)
	num := sqlc.ConvertUint64ToNumeric(n)
	converted := sqlc.ConvertNumericToUint64(num)

	if converted != n {
		t.Errorf("expected %d, got %d", n, converted)
	}
}
