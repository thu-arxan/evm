package evm

// Here defines some kind of gas costs
const (
	GasZero          uint64 = 0
	GasBase          uint64 = 2
	GasVeryLow       uint64 = 3
	GasLow           uint64 = 5
	GasMid           uint64 = 8
	GasHigh          uint64 = 10
	GasExtCode       uint64 = 700
	GasBalance       uint64 = 400
	GasSload         uint64 = 200
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

	GasSha3          uint64 = 1
	GasGetAccount    uint64 = 1
	GasStorageUpdate uint64 = 1
	GasCreateAccount uint64 = 1

	GasBaseOp  uint64 = 0 // TODO: make this 1
	GasStackOp uint64 = 1

	GasEcRecover     uint64 = 1
	GasSha256Word    uint64 = 1
	GasSha256Base    uint64 = 1
	GasRipemd160Word uint64 = 1
	GasRipemd160Base uint64 = 1
	GasExpModWord    uint64 = 1
	GasExpModBase    uint64 = 1
	GasIdentityWord  uint64 = 1
	GasIdentityBase  uint64 = 1
)
