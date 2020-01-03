package gas

// Here defines some gas usage
const (
	Zero          uint64 = 0
	Base          uint64 = 2
	VeryLow       uint64 = 3
	Low           uint64 = 5
	Mid           uint64 = 8
	High          uint64 = 10
	ExtCode       uint64 = 700
	ExtcodeHash   uint64 = 700 //EIP-1884 change it from 400 to 700
	Balance       uint64 = 700 //EIP-1884 change it from 400 to 700
	Sload         uint64 = 800 //EIP-1884 change it from 200 to 800
	JumpDest      uint64 = 1
	Sset          uint64 = 20000
	Sclear        uint64 = 5000
	Sreset        uint64 = 5000
	SelfDestruct  uint64 = 5000
	Create        uint64 = 32000
	CodeDeposit   uint64 = 200
	Call          uint64 = 700
	CallValue     uint64 = 9000
	CallStipend   uint64 = 2300
	NewAccount    uint64 = 25000
	Exp           uint64 = 10
	ExpByte       uint64 = 50
	Memory        uint64 = 3
	TxCreate      uint64 = 32000
	TxDataZero    uint64 = 4
	TxDataNonZero uint64 = 68
	Transaction   uint64 = 21000
	Log           uint64 = 375
	LogData       uint64 = 8
	LogTopic      uint64 = 375
	SHA3          uint64 = 30
	SHA3Word      uint64 = 6
	Copy          uint64 = 3
	BlockHash     uint64 = 20
	QuadDivisor   uint64 = 20
	CreateData    uint64 = 200

	CallNewAccount uint64 = 25000 // Paid for CALL when the destination address didn't exist prior.

	// EIP2200 changes many things of Sstore
	SstoreSentryEIP2200      uint64 = 2300  // Minimum gas required to be present for an SSTORE call, not consumed
	SstoreNoopEIP2200        uint64 = 800   // Once per SSTORE operation if the value doesn't change.
	SstoreDirtyEIP2200       uint64 = 800   // Once per SSTORE operation if a dirty value is changed.
	SstoreInitEIP2200        uint64 = 20000 // Once per SSTORE operation from clean zero to non-zero
	SstoreInitRefundEIP2200  uint64 = 19200 // Once per SSTORE operation for resetting to the original zero value
	SstoreCleanEIP2200       uint64 = 5000  // Once per SSTORE operation from clean non-zero to something else
	SstoreCleanRefundEIP2200 uint64 = 4200  // Once per SSTORE operation for resetting to the original non-zero value
	SstoreClearRefundEIP2200 uint64 = 15000 // Once per SSTORE operation for clearing an originally existing storage slot

	Ecrecover          uint64 = 3000
	Sha256Base         uint64 = 60  // Base price for a SHA256 operation
	Sha256PerWord      uint64 = 12  // Per-word price for a SHA256 operation
	Ripemd160Base      uint64 = 600 // Base price for a RIPEMD160 operation
	Ripemd160PerWord   uint64 = 120 // Per-word price for a RIPEMD160 operation
	IdentityBase       uint64 = 15  // Base price for a data copy operation
	IdentityPerWord    uint64 = 3   // Per-work price for a data copy operation
	ModExpQuadCoeffDiv uint64 = 20  // Divisor for the quadratic particle of the big int modular exponentiation

	Bn256AddByzantium             uint64 = 500    // Byzantium gas needed for an elliptic curve addition
	Bn256AddIstanbul              uint64 = 150    // Gas needed for an elliptic curve addition
	Bn256ScalarMulByzantium       uint64 = 40000  // Byzantium gas needed for an elliptic curve scalar multiplication
	Bn256ScalarMulIstanbul        uint64 = 6000   // Gas needed for an elliptic curve scalar multiplication
	Bn256PairingBaseByzantium     uint64 = 100000 // Byzantium base price for an elliptic curve pairing check
	Bn256PairingBaseIstanbul      uint64 = 45000  // Base price for an elliptic curve pairing check
	Bn256PairingPerPointByzantium uint64 = 80000  // Byzantium per-point price for an elliptic curve pairing check
	Bn256PairingPerPointIstanbul  uint64 = 34000  // Per-point price for an elliptic curve pairing check
)
