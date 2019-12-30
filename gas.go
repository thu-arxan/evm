package evm

// Here defines some kind of gas costs
// TODO: We try our best to support EIP-2200
const (
	GasZero          uint64 = 0
	GasBase          uint64 = 2
	GasVeryLow       uint64 = 3
	GasLow           uint64 = 5
	GasMid           uint64 = 8
	GasHigh          uint64 = 10
	GasExtCode       uint64 = 700
	GasBalance       uint64 = 700 //EIP-1884 change it from 400 to 700
	GasSload         uint64 = 800 //EIP-1884 change it from 200 to 800
	GasJumpDest      uint64 = 1
	GasSset          uint64 = 20000
	GasSclear        uint64 = 5000
	GasSreset        uint64 = 5000
	GasSelfDestruct  uint64 = 5000
	GasCreate        uint64 = 32000
	GasCodeDeposit   uint64 = 200
	GasCall          uint64 = 700
	GasCallValue     uint64 = 9000
	GasCallStipend   uint64 = 2300
	GasNewAccount    uint64 = 25000
	GasExp           uint64 = 10
	GasExpByte       uint64 = 50
	GasMemory        uint64 = 3
	GasTxCreate      uint64 = 32000
	GasTxDataZero    uint64 = 4
	GasTxDataNonZero uint64 = 68
	GasTransaction   uint64 = 21000
	GasLog           uint64 = 375
	GasLogData       uint64 = 8
	GasLogTopic      uint64 = 375
	GasSHA3          uint64 = 30
	GasSHA3Word      uint64 = 6
	GasCopy          uint64 = 3
	GasBlockHash     uint64 = 20
	GasQuadDivisor   uint64 = 20
	GasCreateData    uint64 = 200

	// EIP2200 changes many things of Sstore
	GasSstoreSentryEIP2200      uint64 = 2300  // Minimum gas required to be present for an SSTORE call, not consumed
	GasSstoreNoopEIP2200        uint64 = 800   // Once per SSTORE operation if the value doesn't change.
	GasSstoreDirtyEIP2200       uint64 = 800   // Once per SSTORE operation if a dirty value is changed.
	GasSstoreInitEIP2200        uint64 = 20000 // Once per SSTORE operation from clean zero to non-zero
	GasSstoreInitRefundEIP2200  uint64 = 19200 // Once per SSTORE operation for resetting to the original zero value
	GasSstoreCleanEIP2200       uint64 = 5000  // Once per SSTORE operation from clean non-zero to something else
	GasSstoreCleanRefundEIP2200 uint64 = 4200  // Once per SSTORE operation for resetting to the original non-zero value
	GasSstoreClearRefundEIP2200 uint64 = 15000 // Once per SSTORE operation for clearing an originally existing storage slot
)
