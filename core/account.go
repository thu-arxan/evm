package core

// Account defines an account structure
type Account struct {
	address *Address
	code    []byte
	balance uint64
}

// NewAccount is the constructor of Account
func NewAccount(address *Address) *Account {
	return &Account{
		address: address,
	}
}

// SetCode set code
func (a *Account) SetCode(code []byte) {
	a.code = code
}

// GetAddress return the address
func (a *Account) GetAddress() *Address {
	return a.address
}

// GetBalance return the balance
func (a *Account) GetBalance() uint64 {
	return a.balance
}

// GetCode return the code
func (a *Account) GetCode() []byte {
	return a.code
}

// GetCodeHash return the hash of code
// todo: not implementation yet
func (a *Account) GetCodeHash() []byte {
	return nil
}

// AddBalance add balance
// todo: if overflow
func (a *Account) AddBalance(balance uint64) error {
	a.balance += balance
	return nil
}

// SubBalance sub balance
// todo: avoid overflow
func (a *Account) SubBalance(balance uint64) error {
	a.balance -= balance
	return nil
}
