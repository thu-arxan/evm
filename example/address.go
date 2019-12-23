package example

// Address is the address
type Address [20]byte

// Bytes is the implementation of interface
func (a *Address) Bytes() []byte {
	return a[:]
}
