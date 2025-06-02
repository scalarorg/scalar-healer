package utils

import (
	"encoding/hex"

	"github.com/btcsuite/btcd/txscript"
	"github.com/scalarorg/scalar-healer/constants"
)

func ValidateLockingScript(script string) error {
	lockingScript, err := hex.DecodeString(script)
	if err != nil {
		return err
	}

	class := txscript.GetScriptClass(lockingScript)
	if class == txscript.NonStandardTy {
		return constants.ErrInvalidLockingScript
	}
	return nil
}
