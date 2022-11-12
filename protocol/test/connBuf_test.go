package protocol_test

import (
	"sync"
	"testing"

	"github.com/gaukas/passthru/protocol"
)

var (
	cBuf *protocol.ConnBuf
)

func TestConnBuf(t *testing.T) {
	testNewConnBuf(t)
	testWrite(t)
	testRead(t)
	testWriteConcurrent(t)
	testReadConcurrent(t)
	testWrite(t)
	testPeek(t)
	cBuf.Close()
	testPeekAfterClose(t)
	testReadAfterClose(t)
	testReadEmptyBufAfterClose(t)

	// TODO: test downstream writing
}

func testNewConnBuf(t *testing.T) {
	cBuf = protocol.NewConnBuf()
	if cBuf == nil {
		t.Errorf("Error creating ConnBuf")
	}
}

func testWrite(t *testing.T) {
	cBuf.Write([]byte("test"))
}

func testRead(t *testing.T) {
	p := make([]byte, 1024)
	n, err := cBuf.Read(p)
	if err != nil {
		t.Errorf("Error reading from ConnBuf: %s", err)
	}
	if n != 4 {
		t.Errorf("Wrong number of bytes read: %d", n)
	}
	if string(p[:n]) != "test" {
		t.Errorf("Error reading from ConnBuf")
	}
}

func testWriteConcurrent(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg2 := &sync.WaitGroup{}
	wg.Add(5)
	wg2.Add(5)
	for i := 0; i < 5; i++ {
		go func(wg *sync.WaitGroup, wg2 *sync.WaitGroup, t *testing.T) {
			wg.Done()
			wg.Wait()
			n, err := cBuf.Write([]byte("test"))
			if err != nil {
				t.Errorf("Error writing to ConnBuf: %s", err)
			}
			if n != 4 {
				t.Errorf("Wrong number of bytes written: %d", n)
			}
			wg2.Done()
		}(wg, wg2, t)
	}
	wg2.Wait()
}

func testReadConcurrent(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg2 := &sync.WaitGroup{}
	wg.Add(5)
	wg2.Add(5)
	for i := 0; i < 5; i++ {
		go func(wg *sync.WaitGroup, wg2 *sync.WaitGroup, t *testing.T) {
			wg.Done()
			wg.Wait()
			p := make([]byte, 4)
			n, err := cBuf.Read(p)
			if err != nil {
				t.Errorf("Error reading from ConnBuf: %s", err)
			}
			if n != 4 {
				t.Errorf("Wrong number of bytes read: %d", n)
			}
			if string(p[:n]) != "test" {
				t.Errorf("Error reading from ConnBuf")
			}
			wg2.Done()
		}(wg, wg2, t)
	}

	wg2.Wait()
}

func testPeek(t *testing.T) {
	p := make([]byte, 5)
	err := cBuf.Peek(p, 4)
	if err != nil {
		t.Errorf("Error peeking from ConnBuf: %s", err)
	}
	if string(p[:4]) != "test" {
		t.Errorf("Error peeking from ConnBuf")
	}
}

func testPeekAfterClose(t *testing.T) {
	p := make([]byte, 5)
	err := cBuf.Peek(p, 4)
	if err == nil {
		t.Errorf("Peek should return error after close")
	}
}

func testReadAfterClose(t *testing.T) {
	p := make([]byte, 1024)
	n, err := cBuf.Read(p)
	if err != nil {
		t.Errorf("Error reading from ConnBuf: %s", err)
	}
	if n != 4 {
		t.Errorf("Wrong number of bytes read: %d", n)
	}
}

func testReadEmptyBufAfterClose(t *testing.T) {
	p := make([]byte, 1024)
	n, err := cBuf.Read(p)
	if err == nil {
		t.Errorf("Should error when reading from empty buffer after close")
	}
	if n != 0 {
		t.Errorf("Wrong number of bytes read: %d", n)
	}
}
