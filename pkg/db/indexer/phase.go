package indexer

import (
	"context"

	"github.com/scalarorg/data-models/chains"
)

// SELECT * FROM (
// 	SELECT *, ROW_NUMBER() OVER (
// 			PARTITION BY custodian_group_uid
// 			ORDER BY block_number DESC
// 	) AS row_num
// 	FROM switched_phases
// 	WHERE chain = 'evm|11155111' AND custodian_group_uid IN ('c0b6f4f549aa224fc7b6387fa5f14c77595c83308024f90f5dc0b8afae15be34', 'bffb71bf819ae4cb65188905ac54763a09144bc3a0629808d7142dd5dbd98693')
// ) sub
// WHERE row_num <= 2
// ORDER BY custodian_group_uid, block_number DESC

func (r *IndexerRepository) GetNumberOfLatestSwitchedPhaseEvents(ctx context.Context, numberOfEvents int, chain string, grUID string) ([]chains.SwitchedPhase, error) {
	var events []chains.SwitchedPhase
	err := r.DB.WithContext(ctx).
		Where("chain = ?", chain).
		Where("custodian_group_uid = ?", grUID).
		Order("block_number desc").
		Limit(numberOfEvents).
		Find(&events).Error

	if err != nil {
		return nil, err
	}
	return events, nil
}

func (r *IndexerRepository) GetBatchNumberOfLatestSwitchedPhaseEvents(
	ctx context.Context,
	numberOfEvents int,
	chain string,
	grUID []string) (
	map[string][]chains.SwitchedPhase, error) {
	var switchedPhases []chains.SwitchedPhase
	query := `SELECT * FROM (
		SELECT *, ROW_NUMBER() OVER (
				PARTITION BY custodian_group_uid
				ORDER BY block_number DESC
		) AS row_num
		FROM switched_phases
		WHERE chain = ? AND custodian_group_uid IN (?)
	) sub
	WHERE row_num <= ?
	ORDER BY custodian_group_uid, block_number DESC`
	if err := r.DB.Raw(query, chain, grUID, numberOfEvents).Scan(&switchedPhases).Error; err != nil {
		return nil, err
	}

	switchedPhasesMap := make(map[string][]chains.SwitchedPhase)
	for _, switchedPhase := range switchedPhases {
		if _, ok := switchedPhasesMap[switchedPhase.CustodianGroupUid]; !ok {
			switchedPhasesMap[switchedPhase.CustodianGroupUid] = make([]chains.SwitchedPhase, 0, numberOfEvents)
		}
		switchedPhasesMap[switchedPhase.CustodianGroupUid] = append(switchedPhasesMap[switchedPhase.CustodianGroupUid], switchedPhase)
	}

	return switchedPhasesMap, nil
}

func (r *IndexerRepository) GetBatchLastestSwitchedPhaseEvents(
	ctx context.Context,
	chain string,
	grUID []string) (
	map[string]chains.SwitchedPhase, error) {

	var switchedPhases []chains.SwitchedPhase
	query :=
		`SELECT DISTINCT ON (custodian_group_uid) *
	FROM switched_phases
	WHERE chain = ? AND custodian_group_uid IN (?)
	ORDER BY custodian_group_uid, block_number DESC`
	if err := r.DB.Raw(query, chain, grUID).Scan(&switchedPhases).Error; err != nil {
		return nil, err
	}

	switchedPhasesMap := make(map[string]chains.SwitchedPhase)
	for _, switchedPhase := range switchedPhases {
		switchedPhasesMap[switchedPhase.CustodianGroupUid] = switchedPhase
	}
	return switchedPhasesMap, nil
}
