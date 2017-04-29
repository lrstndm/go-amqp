package testconn

import (
	"bytes"
	"errors"
	"io"
	"net"
	"time"
)

func New(data []byte) *Conn {
	c := &Conn{
		data: bytes.Split(data, []byte("SPLIT\n")),
		done: make(chan struct{}),
		err:  make(chan error, 1),
	}
	return c
}

type Conn struct {
	data [][]byte
	// data         []byte
	done         chan struct{}
	write        chan struct{}
	err          chan error
	readDeadline *time.Timer
}

func (c *Conn) Read(b []byte) (int, error) {
	if len(c.data) == 0 {
		select {
		case <-c.done:
			return 0, io.EOF
		case err := <-c.err:
			return 0, err
		}
	}
	time.Sleep(1 * time.Millisecond)
	// // fmt.Printf("DEBUG: %q\n", c.data[0])
	n := copy(b, c.data[0])
	c.data = c.data[1:]
	// n := copy(b, c.data)
	// c.data = c.data[n:]
	return n, nil
}

func (c *Conn) Write(b []byte) (n int, err error) {
	return len(b), nil
}

func (c *Conn) Close() error {
	close(c.done)
	return nil
}

func (c *Conn) LocalAddr() net.Addr {
	return &net.TCPAddr{net.IP{127, 0, 0, 1}, 49706, ""}
}

func (c *Conn) RemoteAddr() net.Addr {
	return &net.TCPAddr{net.IP{127, 0, 0, 1}, 49706, ""}
}

func (c *Conn) SetDeadline(t time.Time) error {
	return nil
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	if c.readDeadline == nil {
		c.readDeadline = time.AfterFunc(t.Sub(time.Now()), func() {
			select {
			case c.err <- errors.New("timeout"):
			default:
			}
		})
		return nil
	}
	c.readDeadline.Reset(t.Sub(time.Now()))
	return nil
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return nil
}