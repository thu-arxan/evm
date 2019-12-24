package core

// Here defines some consts
const (
	Word160Length       = 20
	Word256Word160Delta = 12
)

// Zero160 is the zero of Word160
var Zero160 = Word160{}

// Word160 is bytes which length is Word160Length
type Word160 [Word160Length]byte

// Word256 convert Word160 to Word256
// The function will add zeros before Word160 until its length == 32
func (w Word160) Word256() (word256 Word256) {
	copy(word256[Word256Word160Delta:], w[:])
	return
}

// Bytes return bytes of Word160
func (w Word160) Bytes() []byte {
	return w[:]
}
