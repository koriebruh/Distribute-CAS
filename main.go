package main

import (
	"github.com/koribruh/Distribute-CAS/p2p"
	"log"
)

func main() {

	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    "127.0.0.1:3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}

	tr := p2p.NewTCPTransport(tcpOpts)

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}

}
