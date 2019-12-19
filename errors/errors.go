package errors

// Sink are useful for recording errors but continuing with computation. Implementations may choose to just store
// the first error pushed and ignore subsequent ones or they may record an error trace. Pushing a nil error should have
// no effects.
type Sink interface {
	// Push an error to the error. If a nil error is passed then that value should not be pushed. Returns true if error
	// is non nil.
	PushError(error)
	Error() error
}

// Maybe is an implementation of Sink
type Maybe struct {
	err error
}

// NewMaybe is the constructor of Sink
func NewMaybe() Sink {
	return &Maybe{
		err: nil,
	}
}

// PushError is the implementation of interface
func (m *Maybe) PushError(err error) {
	if err != nil && m.err == nil {
		m.err = err
	}
}

// Error is the implementation of interface
func (m *Maybe) Error() error {
	return m.err
}
