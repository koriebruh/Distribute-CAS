package p2p

// Peer represented the remote node
type Peer interface {
}

// Transport handle communication between the node in network can be form
// (TCP, UDP, WebSockets ..)
type Transport interface {
	ListenAndAccept() error
}
