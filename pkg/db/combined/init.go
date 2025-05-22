package combined

import (
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/db/healer"
	"github.com/scalarorg/scalar-healer/pkg/db/indexer"
)

type CombinedManager struct {
	*healer.HealerRepository
	*indexer.IndexerRepository
}

var _ db.CombinedAdapter = (*CombinedManager)(nil)

func NewCombinedManager(healerDB *healer.HealerRepository, indexerDB *indexer.IndexerRepository) *CombinedManager {
	return &CombinedManager{
		healerDB,
		indexerDB,
	}
}

func (m *CombinedManager) Close() {
	m.HealerRepository.Close()
	m.IndexerRepository.Close()
}
