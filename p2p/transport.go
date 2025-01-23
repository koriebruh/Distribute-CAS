package p2p

// Peer represented the remote node
type Peer interface {
	Close() error
}

// Transport handle communication between the node in network can be form
// (TCP, UDP, WebSockets ..)
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC // return only
}
