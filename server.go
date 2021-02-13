package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

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

type RequestMsg struct {
	Cmd  string `json:"cmd"`
	Data string `json:"data"`
}

type ResponseMsg struct {
	Msg    string `json:"msg"`
	Result string `json:"result"`
}

type Player struct {
	conn *net.TCPConn
	Name string
	Role Role
}

func (p *Player) Init(conn *net.TCPConn, name string) {
	p.conn = conn
	p.Name = name
}

func (p *Player) Send(msg []byte) (int, error) {
	return p.conn.Write(msg)
}

type Room struct {
	Id      int       `json:"id"`
	Name    string    `json:"name"`
	Players []*Player `json:"players"`
	White   *Player   `json:"white"`
	Black   *Player   `json:"black"`
	State   string    `json:"state"`
}

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      uint
	Rooms     map[int]*Room
}

func (s *Server) Init(name string, port uint) {
	s.Name = name
	s.IPVersion = `tcp4`
	s.IP = `0.0.0.0`
	if port != 0 {
		s.Port = port
	} else {
		s.Port = 8901
	}
	// test rooms
	s.Rooms = make(map[int]*Room)
	var r *Room
	r = new(Room)
	r.Id = 1
	r.Name = "新手房"
	r.State = "1/2"
	s.addRoom(r)
	r = new(Room)
	r.Id = 2
	r.Name = "高手房"
	r.State = "2/2"
	s.addRoom(r)
	r = new(Room)
	r.Id = 3
	r.Name = "私人房"
	r.State = "游戏中"
	s.addRoom(r)
}

func (s *Server) Run() {
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		log.Println("Resolve TCP Address Error:  ", err)
		return
	}

	listenner, err := net.ListenTCP(s.IPVersion, addr)
	defer listenner.Close()
	if err != nil {
		log.Println("listen ", s.IPVersion, "err: ", err)
		return
	}
	log.Println("Server was running at ", addr)
	for {
		conn, err := listenner.AcceptTCP()
		if err != nil {
			log.Println("Accept Tcp Error: ", err)
			continue
		}
		log.Println("Accept: ", conn.RemoteAddr())
		go s.handler(conn)
	}
}

func (s *Server) handler(conn *net.TCPConn) {
	chess := new(Chess)
	player := new(Player)
	player.Init(conn, "")
	for {
		buf := make([]byte, 2048)
		cnt, err := conn.Read(buf)
		if err == io.EOF {
			log.Println(conn.RemoteAddr(), "disconnected unexpected")
			break
		}
		if err != nil {
			log.Println("err :", err)
			time.Sleep(time.Second * 3)
			continue
		}
		msg := new(RequestMsg)
		if err := json.Unmarshal(buf[:cnt], msg); err != nil {
			log.Println(err)
		}
		log.Println(msg.Cmd)
		switch strings.ToUpper(msg.Cmd) {
		case "PING":
			r := &ResponseMsg{
				Msg: "PONG",
			}
			msg, _ := json.Marshal(r)
			player.Send(msg)
		case "EXIT":
			log.Println(conn.RemoteAddr(), " exit")
			r := &ResponseMsg{
				Msg: "EXIT",
			}
			msg, _ := json.Marshal(r)
			player.Send(msg)
			return
		case "START":
			chess.Init()
		case "RESET":
			chess.Reset()
		case "GET ROOMS":
			result, _ := json.Marshal(s.Rooms)
			r := &ResponseMsg{
				Msg:    "ROOMS",
				Result: string(result),
			}
			msg, _ := json.Marshal(r)
			player.Send(msg)
		case "JOIN ROOM":
			var data map[string]int
			json.Unmarshal([]byte(msg.Data), &data)
			id := data["id"]
			log.Println(id)
			if r, ok := s.Rooms[id]; ok {
				if r.Black == nil {
					r.Black = player
				} else if r.White == nil {
					r.White = player
				}
				resp := &ResponseMsg{
					Msg: "JOIN ROOM",
				}
				info := fmt.Sprintf("%d, %s, %s\n", r.Id, r.Name, r.State)
				if r.Black != nil {
					info += fmt.Sprintf("b %s\n", r.Black.Name)
				}
				if r.White != nil {
					info += fmt.Sprintf("w %s\n", r.White.Name)
				}
				resp.Result = info
				msg, _ := json.Marshal(resp)
				log.Println(msg)
				player.Send(msg)
			}
		case "PLACE PIECE":
			type TurnData struct {
				Turn []byte `json:"turn"`
				PosX uint8  `json:"x"`
				PoxY uint8  `json:"y"`
			}
			data := &TurnData{}
			json.Unmarshal([]byte(msg.Data), data)
			if data.Turn[0] == chess.Turn {
				chess.PlacePiece(int(data.PosX), int(data.PoxY))
			}
		default:
			log.Println(msg)
		}
	}
}

func (s *Server) addRoom(r *Room) {
	if _, ok := s.Rooms[r.Id]; ok {
		log.Println("Already Exists")
	} else {
		s.Rooms[r.Id] = r
	}
}
