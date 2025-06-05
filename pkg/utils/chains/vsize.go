package chains

const (
	P2TR_INPUT_SIZE                           uint64 = 58 // 57.5
	P2TR_OUTPUT_SIZE                          uint64 = 43
	P2TR_BUFFER_SIZE                          uint64 = 11 // 10.5
	ESTIMATE_SIGNATURE_COST                   uint64 = 16
	ESTIMATE_ADDITIONAL_P2TR_SCRIPT_PATH_COST uint64 = 60
)

func CalculateVsize(inputs int, outputs int, quorum uint64) uint64 {
	witnessSize := ESTIMATE_SIGNATURE_COST*quorum + ESTIMATE_ADDITIONAL_P2TR_SCRIPT_PATH_COST
	inputsSize := (P2TR_INPUT_SIZE + witnessSize) * uint64(inputs)
	outputsSize := P2TR_OUTPUT_SIZE * uint64(outputs)
	return P2TR_BUFFER_SIZE + inputsSize + outputsSize
}
