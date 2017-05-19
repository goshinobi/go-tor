package cluster

import (
	"encoding/json"
	"errors"

	"github.com/goshinobi/go-tor/node"
)

type Cluster struct {
	Name  *string      `json:name`
	Nodes []*node.Node `json:nodes`
}

func New(name ...string) *Cluster {
	var n *string
	if len(name) != 0 {
		n = &name[0]
	} else {
		n = nil
	}
	return &Cluster{
		Name:  n,
		Nodes: make([]*node.Node, 0),
	}
}

func (c *Cluster) Add() error {
	var n *node.Node
	if c.Name == nil {
		n = node.New()
	} else {
		n = node.New(*c.Name)
	}
	if n == nil {
		return errors.New("can not add new node")
	}

	c.Nodes = append(c.Nodes, n)
	return nil
}

func (c *Cluster) KillAll() error {
	var errs []error
	for _, n := range c.Nodes {
		if err := n.Kill(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	bin, _ := json.Marshal(errs)
	return errors.New(string(bin))
}
