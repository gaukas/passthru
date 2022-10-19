package protocol

import (
	"errors"
	"io"
	"sync"
)

var (
	ErrNotEnoughData = errors.New("not enough data in buffer")
)

// Thread-safe buffer in order to allow inspection of the connection
type ConnBuf struct {
	mutex  sync.RWMutex
	buf    []byte
	closed bool
}

func NewConnBuf() *ConnBuf {
	return &ConnBuf{
		buf: make([]byte, 0),
	}
}

func (cb *ConnBuf) Write(p []byte) (n int, err error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	if cb.closed {
		return 0, io.ErrClosedPipe
	}

	cb.buf = append(cb.buf, p...)
	return len(p), nil
}

func (cb *ConnBuf) Read(p []byte) (n int, err error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	n = copy(p, cb.buf)
	if n == 0 && cb.closed {
		return 0, io.EOF
	}
	cb.buf = cb.buf[n:]
	return
}

func (cb *ConnBuf) Close() error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	cb.closed = true
	return nil
}

// Peek copies at least n bytes from the buffer, or return error
func (cb *ConnBuf) Peek(p []byte, n int) error {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	if cb.closed { // if closed, no need to peek as it can't be used anymore
		return io.EOF
	}

	nRead := copy(p, cb.buf)
	if nRead < n {
		return ErrNotEnoughData
	}
	return nil
}
