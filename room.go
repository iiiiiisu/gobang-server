package main

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"
)

type Room struct {
	Id      int       `json:"id"`
	Name    string    `json:"name"`
	Host    *Player   `json:"host"`
	Players []*Player `json:"players"`
	White   *Player   `json:"white"`
	Black   *Player   `json:"black"`
	State   string    `json:"state"`
}

func (r *Room) Init(p *Player) {
	r.Host = p
	r.Join(p)
	r.State = "准备中"
	// init id
	s := []byte(strconv.FormatInt(time.Now().UnixNano(), 10))
	r.Id, _ = strconv.Atoi(fmt.Sprintf("%x", md5.Sum(s)))
}

func (r *Room) destroy() {

}

func (r *Room) Join(p *Player) {
	if r.Black == nil {
		r.Black = p
	} else if r.White == nil {
		r.White = p
	}
	r.Players = append(r.Players, p)
	r.broadCast([]byte("{}"))
}

func (r *Room) Leave(p *Player) {
	r.removePlayer(p)
	if p == r.Host {
		if len(r.Players) <= 1 {
			r.destroy()
		} else {
			r.Host = nil
			r.selectHost()
		}
	}
	r.broadCast([]byte("{}"))
}

func (r *Room) removePlayer(p *Player) {
	var i int
	for i, _ = range r.Players {
		if r.Players[i] != p {
			continue
		}
		break
	}
	r.Players = append(r.Players[:i], r.Players[i+1:]...)
	if r.Black == p {
		r.Black = nil
	} else if r.White == p {
		r.White = nil
	}
}

func (r *Room) selectHost() {
	if r.Black != nil {
		r.Host = r.Black
	} else if r.White != nil {
		r.Host = r.White
	} else {
		r.Host = r.Players[0]
	}
}

func (r *Room) broadCast(msg []byte) {
	for i, _ := range r.Players {
		r.Players[i].Send(msg)
	}
}
