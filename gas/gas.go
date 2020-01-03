package gas

// Here defines some gas usage
const (
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
