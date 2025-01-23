package p2p

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// TCPPeer represent node lain yang terhubung melalui koneksi TCP.
type TCPPeer struct {
	// Connection TCP ke node lain.
	conn net.Conn
	// True if we're starting conn (outbound).
	// False if we receive conn (inbound).
	outbound bool
	wg       *sync.WaitGroup
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn: conn, outbound: outbound,
	}
}

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(peer Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC),
	}
}

// Consume return read only channel for reading incoming message receive from another peer
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("DYUM TCP accept error : %s\n", err)
		}

		fmt.Printf("new connection established from %s: %+v\n", conn.RemoteAddr().String(), conn)

		go t.handleConnection(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConnection(conn net.Conn) {
	var err error

	defer func() {
		fmt.Printf("dropping the connection: %s\n", err)
		defer conn.Close()
	}()

	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			fmt.Fprintln(conn, "error dropping connection!") // Kirim error ke client
			time.Sleep(100 * time.Millisecond)
			return
		}
	}

	fmt.Fprintln(conn, "Welcome to the server!")
	msg := RPC{}
	for {
		if err = t.Decoder.Decode(conn, &msg); err != nil {
			fmt.Fprintf(conn, "Error reading message: %s\n", err)
			return
		}

		msg.From = conn.RemoteAddr().String()
		t.rpcch <- msg
	}

}
