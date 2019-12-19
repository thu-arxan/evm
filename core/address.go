package core

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

// AddressFromWord256 convert Word256 to Address
func AddressFromWord256(addr Word256) Address {
	return Address(addr.Word160())
}
