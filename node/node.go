package node

import (
	"net"

	"github.com/goshinobi/go-tor/tor"
)

type Node struct {
	Server    *tor.Tor `json:"server_info"`
	Freshness int      `json:"freshness"`
	IP        *net.IP  `json:"ip"`
}

func New(group ...string) *Node {
	server := tor.New(group)
	if server == nil {
		return nil
	}
	if err := server.Start(); err != nil {
		return nil
	}
	return &Node{
		server,
		0,
		nil,
	}
}

func (p *Node) Start() error {
	return p.Server.Start()
}

func (p *Node) Kill() error {
	return p.Server.Kill()
}

func (p *Node) Reload() error {
	return p.Server.Reload()
}
