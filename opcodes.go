package evm

// OpCode is the type of operation code
type OpCode byte

// 0s: Stop and Arithmetic Operations
const (
	STOP OpCode = iota
	ADD
	MUL
	SUB
	DIV
	SDIV
	MOD
	SMOD
	ADDMOD
	MULMOD
	EXP
	SIGNEXTEND
)

// 10s: Comparison & Bitwise Logic Opreations
const (
	LT OpCode = iota + 0x10
	GT
	SLT
	SGT
	EQ
	ISZERO
	AND
	OR
	XOR
	NOT
	BYTE
	SHL
	SHR
	SAR
)

// 20s: SHA3
const (
	SHA3 = 0x20
)

// 30s: Environmental Information
const (
	ADDRESS OpCode = 0x30 + iota
	BALANCE
	ORIGIN
	CALLER
	CALLVALUE
	CALLDATALOAD
	CALLDATASIZE
	CALLDATACOPY
	CODESIZE
	CODECOPY
	GASPRICE
	EXTCODESIZE
	EXTCODECOPY
	RETURNDATASIZE
	RETURNDATACOPY
	// Note: This is new
	EXTCODEHASH
)

// 40s: Block Information
const (
	BLOCKHASH OpCode = 0x40 + iota
	COINBASE
	TIMESTAMP
	NUMBER
	DIFFICULTY
	GASLIMIT
)

// 50s: Stack, Memory, Storage and Flow Operations
const (
	POP OpCode = 0x50 + iota
	MLOAD
	MSTORE
	MSTORE8
	SLOAD
	SSTORE
	JUMP
	JUMPI
	PC
	MSIZE
	GAS
	JUMPDEST
)

// 60s & 70s: Push Operations
const (
	PUSH1 OpCode = 0x60 + iota
	PUSH2
	PUSH3
	PUSH4
	PUSH5
	PUSH6
	PUSH7
	PUSH8
	PUSH9
	PUSH10
	PUSH11
	PUSH12
	PUSH13
	PUSH14
	PUSH15
	PUSH16
	PUSH17
	PUSH18
	PUSH19
	PUSH20
	PUSH21
	PUSH22
	PUSH23
	PUSH24
	PUSH25
	PUSH26
	PUSH27
	PUSH28
	PUSH29
	PUSH30
	PUSH31
	PUSH32
)

// 80s: Duplication Operations
const (
	DUP1 = 0x80 + iota
	DUP2
	DUP3
	DUP4
	DUP5
	DUP6
	DUP7
	DUP8
	DUP9
	DUP10
	DUP11
	DUP12
	DUP13
	DUP14
	DUP15
	DUP16
)

// 90s: Exchange Operations
const (
	SWAP1 = 0x90 + iota
	SWAP2
	SWAP3
	SWAP4
	SWAP5
	SWAP6
	SWAP7
	SWAP8
	SWAP9
	SWAP10
	SWAP11
	SWAP12
	SWAP13
	SWAP14
	SWAP15
	SWAP16
)

// a0s: Logging Operations
const (
	LOG0 OpCode = 0xa0 + iota
	LOG1
	LOG2
	LOG3
	LOG4
)

// f0s: System Operations
const (
	CREATE OpCode = 0xf0 + iota
	CALL
	CALLCODE
	RETURN
	DELEGATECALL
	CREATE2
	STATICCALL = 0xfa

	REVERT       = 0xfd
	INVALID      = 0xfe
	SELFDESTRUCT = 0xff
)

var opCodeNames = map[OpCode]string{
	// 0s: Stop and Arithmetic Operations
	STOP:       "STOP",
	ADD:        "ADD",
	MUL:        "MUL",
	SUB:        "SUB",
	DIV:        "DIV",
	SDIV:       "SDIV",
	MOD:        "MOD",
	SMOD:       "SMOD",
	ADDMOD:     "ADDMOD",
	MULMOD:     "MULMOD",
	EXP:        "EXP",
	SIGNEXTEND: "SIGNEXTEND",

	// 10s: Comparison & Bitwise Logic Opreations
	LT:     "LT",
	GT:     "GT",
	SLT:    "SLT",
	SGT:    "SGT",
	EQ:     "EQ",
	ISZERO: "ISZERO",
	AND:    "AND",
	OR:     "OR",
	XOR:    "XOR",
	NOT:    "NOT",
	BYTE:   "BYTE",
	SHL:    "SHL",
	SHR:    "SHR",
	SAR:    "SAR",

	// 20s: SHA3
	SHA3: "SHA3",

	// 30s: Environmental Information
	ADDRESS:        "ADDRESS",
	BALANCE:        "BALANCE",
	ORIGIN:         "ORIGIN",
	CALLER:         "CALLER",
	CALLVALUE:      "CALLVALUE",
	CALLDATALOAD:   "CALLDATALOAD",
	CALLDATASIZE:   "CALLDATASIZE",
	CALLDATACOPY:   "CALLDATACOPY",
	CODESIZE:       "CODESIZE",
	CODECOPY:       "CODECOPY",
	GASPRICE:       "GASPRICE",
	EXTCODESIZE:    "EXTCODESIZE",
	EXTCODECOPY:    "EXTCODECOPY",
	RETURNDATASIZE: "RETURNDATASIZE",
	RETURNDATACOPY: "RETURNDATACOPY",
	// Note: This is new
	EXTCODEHASH: "EXTCODEHASH",

	// 40s: Block Information
	BLOCKHASH:  "BLOCKHASH",
	COINBASE:   "COINBASE",
	TIMESTAMP:  "TIMESTAMP",
	NUMBER:     "NUMBER",
	DIFFICULTY: "DIFFICULTY",
	GASLIMIT:   "GASLIMIT",

	// 50s: Stack, Memory, Storage and Flow Operations
	POP:      "POP",
	MLOAD:    "MLOAD",
	MSTORE:   "MSTORE",
	MSTORE8:  "MSTORE8",
	SLOAD:    "SLOAD",
	SSTORE:   "SSTORE",
	JUMP:     "JUMP",
	JUMPI:    "JUMPI",
	PC:       "PC",
	MSIZE:    "MSIZE",
	GAS:      "GAS",
	JUMPDEST: "JUMPDEST",

	// 60s & 70s: Push Operations
	PUSH1:  "PUSH1",
	PUSH2:  "PUSH2",
	PUSH3:  "PUSH3",
	PUSH4:  "PUSH4",
	PUSH5:  "PUSH5",
	PUSH6:  "PUSH6",
	PUSH7:  "PUSH7",
	PUSH8:  "PUSH8",
	PUSH9:  "PUSH9",
	PUSH10: "PUSH10",
	PUSH11: "PUSH11",
	PUSH12: "PUSH12",
	PUSH13: "PUSH13",
	PUSH14: "PUSH14",
	PUSH15: "PUSH15",
	PUSH16: "PUSH16",
	PUSH17: "PUSH17",
	PUSH18: "PUSH18",
	PUSH19: "PUSH19",
	PUSH20: "PUSH20",
	PUSH21: "PUSH21",
	PUSH22: "PUSH22",
	PUSH23: "PUSH23",
	PUSH24: "PUSH24",
	PUSH25: "PUSH25",
	PUSH26: "PUSH26",
	PUSH27: "PUSH27",
	PUSH28: "PUSH28",
	PUSH29: "PUSH29",
	PUSH30: "PUSH30",
	PUSH31: "PUSH31",
	PUSH32: "PUSH32",

	// 80s: Duplication Operations
	DUP1:  "DUP1",
	DUP2:  "DUP2",
	DUP3:  "DUP3",
	DUP4:  "DUP4",
	DUP5:  "DUP5",
	DUP6:  "DUP6",
	DUP7:  "DUP7",
	DUP8:  "DUP8",
	DUP9:  "DUP9",
	DUP10: "DUP10",
	DUP11: "DUP11",
	DUP12: "DUP12",
	DUP13: "DUP13",
	DUP14: "DUP14",
	DUP15: "DUP15",
	DUP16: "DUP16",

	// 90s: Exchange Operations
	SWAP1:  "SWAP1",
	SWAP2:  "SWAP2",
	SWAP3:  "SWAP3",
	SWAP4:  "SWAP4",
	SWAP5:  "SWAP5",
	SWAP6:  "SWAP6",
	SWAP7:  "SWAP7",
	SWAP8:  "SWAP8",
	SWAP9:  "SWAP9",
	SWAP10: "SWAP10",
	SWAP11: "SWAP11",
	SWAP12: "SWAP12",
	SWAP13: "SWAP13",
	SWAP14: "SWAP14",
	SWAP15: "SWAP15",
	SWAP16: "SWAP16",

	// a0s: Logging Operations
	LOG0: "LOG0",
	LOG1: "LOG1",
	LOG2: "LOG2",
	LOG3: "LOG3",
	LOG4: "LOG4",

	// f0s: System Operations
	CREATE:       "CREATE",
	CALL:         "CALL",
	CALLCODE:     "CALLCODE",
	RETURN:       "RETURN",
	DELEGATECALL: "DELEGATECALL",
	CREATE2:      "CREATE2",
	STATICCALL:   "STATICCALL",
	REVERT:       "REVERT",
	INVALID:      "INVALID",
	SELFDESTRUCT: "SELFDESTRUCT",
}

func (o OpCode) String() string {
	return opCodeNames[o]
}
