//  Copyright 2020 The THU-Arxan Authors
//  This file is part of the evm library.
//
//  The evm library is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Lesser General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  The evm library is distributed in the hope that it will be useful,/
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
//  GNU Lesser General Public License for more details.
//
//  You should have received a copy of the GNU Lesser General Public License
//  along with the evm library. If not, see <http://www.gnu.org/licenses/>.
//

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
