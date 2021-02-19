package main

import "net"

type Role uint8

const (
	WATCHER Role = iota
	BLACK
	WHITE
)

const (
	WAITTING = "WAITTING"
	PLAYING  = "PLAYING"
)

type Player struct {
	conn *net.TCPConn
	Name string
	Role Role
}

func (p *Player) Init(conn *net.TCPConn) {
	p.conn = conn
}

func (p *Player) Login(username, password string) error {
	return nil
}

func (p *Player) Send(msg []byte) (int, error) {
	return p.conn.Write(msg)
}
