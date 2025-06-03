package utils

import (
	"encoding/hex"

	"github.com/btcsuite/btcd/txscript"
	"github.com/scalarorg/scalar-healer/constants"
)

func ValidateLockingScript(script string) ([]byte, error) {
	lockingScript, err := hex.DecodeString(script)
	if err != nil {
		return nil, err
	}

	class := txscript.GetScriptClass(lockingScript)
	if class == txscript.NonStandardTy {
		return nil, constants.ErrInvalidLockingScript
	}
	return lockingScript, nil
}
