package core

// Account defines an account structure
type Account struct {
	address *Address
	code    []byte
	balance uint64
	nonce   uint64
}

// NewAccount is the constructor of Account
func NewAccount(address *Address) *Account {
	return &Account{
		address: address,
	}
}

// GetAddress return the address
func (a *Account) GetAddress() *Address {
	return a.address
}

// GetBalance return the balance
func (a *Account) GetBalance() uint64 {
	return a.balance
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

// GetCode return the code
func (a *Account) GetCode() []byte {
	return a.code
}

// SetCode set code
func (a *Account) SetCode(code []byte) {
	a.code = code
}

// GetCodeHash return the hash of code
// todo: not implementation yet
func (a *Account) GetCodeHash() []byte {
	return nil
}

// GetNonce return the nonce of account
func (a *Account) GetNonce() uint64 {
	return a.nonce
}

// SetNonce set the nonce of account
func (a *Account) SetNonce(nonce uint64) {
	a.nonce = nonce
}
