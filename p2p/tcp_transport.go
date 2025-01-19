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

type TCPTransport struct {
	listenAddress string
	listener      net.Listener
	handshakeFunc HandshakeFunc
	//decoder       de

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		handshakeFunc: NOPHandshakeFunc,
		listenAddress: listenAddr,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.listenAddress)
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

		go t.handleConnection(conn)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	//fmt.Printf("new connection establiesh %+v\n", peer)
	fmt.Printf("New connection established from %s: %+v\n", conn.RemoteAddr().String(), peer)
}
