package protocol

import (
	"errors"
	"io"
	"sync"
        "github.com/gaukas/passthru/internal/logger"
)

var (
	ErrNotEnoughData = errors.New("not enough data in buffer")
)

// Thread-safe buffer in order to allow inspection of the connection
type ConnBuf struct {
	mutex  sync.RWMutex
	buf    []byte
	closed bool

	// Performance improvement
	bufCleared bool      // safe guard to prevent anything write to downstream while buffer remains.
	downstream io.Writer // if set, will write to this writer instead of the buffer
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

	if cb.downstream != nil {
		if !cb.bufCleared {
			if len(cb.buf) > 0 {
				_, err = cb.downstream.Write(cb.buf)
				if err != nil {
					logger.Errorf("Caught error %v", err)
					return 0, err
				}

				// clear the buffer
				cb.buf = cb.buf[:0]
			}
			// set the flag to prevent anything write to the buffer
			cb.bufCleared = true
		}
		return cb.downstream.Write(p)
	}
	return cb.writeBufferLocked(p)
}

// must be called when caller holds the lock
func (cb *ConnBuf) writeBufferLocked(p []byte) (n int, err error) {
	cb.buf = append(cb.buf, p...)
	return len(p), nil
}

func (cb *ConnBuf) Read(p []byte) (n int, err error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	// doesn't allow read if there is a downstream
	if cb.downstream != nil {
		err = io.ErrClosedPipe
		return
	}

	n = copy(p, cb.buf)
	// keep reading until the buffer is either non-empty or closed
	for n == 0 && !cb.closed {
		n = copy(p, cb.buf)
	}

	if n == 0 { // means cb.closed is true
		err = io.EOF
		return
	}

	// means n > 0
	cb.buf = cb.buf[n:] // advance the buffer

	return
}

func (cb *ConnBuf) Close() error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	cb.closed = true

	if cb.downstream != nil { // close the downstream if it is a closer
		if closer, ok := cb.downstream.(io.Closer); ok {
			return closer.Close()
		}
	}

	return nil
}

// Peek copies at least n bytes from the buffer, or return error
func (cb *ConnBuf) Peek(p []byte, n int) error {
	cb.mutex.RLock() // ReadLock since it is static
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

func (cb *ConnBuf) SetDownstream(w io.Writer) error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	cb.downstream = w

	// if anything left in the buffer, write it to the downstream
	if len(cb.buf) > 0 {
		_, err := cb.downstream.Write(cb.buf)
		// clear the buffer
		cb.buf = cb.buf[:0]
		return err
	}
	return nil
}
