package precompile

// Here defines some gas usage
const (
	GasEcrecover     uint64 = 3000
	GasSha256Base    uint64 = 60 // Base price for a SHA256 operation
	GasSha256PerWord uint64 = 12 // Per-word price for a SHA256 operation
)
