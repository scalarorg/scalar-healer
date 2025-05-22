package testutils_test

import (
	"context"
	"testing"

	"github.com/scalarorg/scalar-healer/pkg/db"
	testutils "github.com/scalarorg/scalar-healer/pkg/test_utils"
)

func TestPrepareTestPostgresDB(t *testing.T) {
	testutils.RunWithTestDB(func(_ context.Context, _ db.HealderAdapter) error {
		return nil
	})
}
