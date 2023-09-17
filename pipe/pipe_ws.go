package pipe

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type PipeWS struct {
	target  net.Conn
	origin  *websocket.Conn
	inbuf   chan []byte
	backbuf chan []byte
}

// 入口流量
func (p *PipeWS) StreamIn() {
	defer func() {
		log.Println("StreamIn closed")
		close(p.inbuf)
	}()

	for {
		_, readbuf, err := p.origin.ReadMessage()
		fmt.Println(readbuf)
		if err != nil {
			fmt.Println(err.Error(), "StreamIn")
			return
		}
		copy_buf := make([]byte, len(readbuf))
		copy(copy_buf, readbuf)
		fmt.Println("reading ", string(copy_buf))
		p.inbuf <- copy_buf
	}
}

// 返回流量
func (p *PipeWS) StreamBack() {

	defer func() {
		log.Println("StreamBack closed")
		close(p.backbuf)
	}()

	buff := make([]byte, bufferSize)

	for {
		sz, err := p.target.Read(buff)
		if err != nil {
			// 问题直接结束
			fmt.Println(err.Error(), "Streamback")
			return
		}
		copy_buf := make([]byte, sz)
		copy(copy_buf, buff)
		p.backbuf <- copy_buf
	}
}

func (p *PipeWS) Working() {
	defer func() {
		log.Println("Stop working")
		p.target.Close()
		p.origin.Close()
	}()

	p.backbuf = make(chan []byte, 10)
	p.inbuf = make(chan []byte, 10)
	for {
		select {
		case buf, ok := <-p.inbuf:
			if ok {
				log.Println("传递", string(buf))
				p.target.Write(buf)
			} else {
				return
			}
		case buf, ok := <-p.backbuf:
			if ok {
				log.Println("返回", string(buf))
				p.origin.SetWriteDeadline(time.Now().Add(time.Second))
				p.origin.WriteMessage(websocket.BinaryMessage, buf)
			} else {
				return
			}
		}
	}

}

func (p *PipeWS) Run() {
	go p.Working()
	go p.StreamIn()
	go p.StreamBack()
}
