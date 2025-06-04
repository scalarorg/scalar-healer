package utils

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
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

func ScriptPubKeyToAddress(scriptPubKey []byte, chainParams *chaincfg.Params) (btcutil.Address, error) {
	_, addresses, _, err := txscript.ExtractPkScriptAddrs(scriptPubKey, chainParams)
	if err != nil {
		return nil, err
	}

	if len(addresses) == 0 {
		return nil, fmt.Errorf("no addresses found")
	}

	return addresses[0], nil
}
