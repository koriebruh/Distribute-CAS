package p2p

import (
	"fmt"
	"net"
	"sync"
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

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	decoder  Decoder

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
	}
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
	defer conn.Close()
	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {
		fmt.Printf("TCP handshake error: %s\n", err)
		return
	}

	fmt.Fprintf(conn, "Welcome to the server!\n\n")

	msg := &Message{}
	for {
		if err := t.Decoder.Decode(conn, msg); err != nil {
			fmt.Fprintf(conn, "Error reading message: %s\n", err)
			break
		}

		msg.From = conn.RemoteAddr().String()
		fmt.Printf("message received: %+v\n", msg)
	}

}
