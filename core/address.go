package core

import "fmt"

// Address is Word160
type Address Word160

// Here defines some length
const (
	AddressLength    = Word160Length
	AddressHexLength = 2 * AddressLength
)

// ZeroAddress is zero
var ZeroAddress = Address{}

// Word256 return Word256 of address
func (address Address) Word256() Word256 {
	return Word160(address).Word256()
}

// Bytes return bytes of address
func (address Address) Bytes() []byte {
	return address.Word256().Bytes()
}

// AddressFromBytes returns an address consisting of the first 20 bytes of bs, return an error if the bs does not have length exactly 20
// but will still return either: the bytes in bs padded on the right or the first 20 bytes of bs truncated in any case.
func AddressFromBytes(bs []byte) (address Address, err error) {
	if len(bs) != Word160Length {
		err = fmt.Errorf("slice passed as address '%X' has %d bytes but should have %d bytes",
			bs, len(bs), Word160Length)
		// It is caller's responsibility to check for errors. If they ignore the error we'll assume they want the
		// best-effort mapping of the bytes passed to an address so we don't return here
	}
	copy(address[:], bs)
	return
}
