package pipe

import (
	"fmt"
	"log"
	"net"
)

const (
	bufferSize = 10 * 1024
)

type Pipe struct {
	origin  net.Conn
	target  net.Conn
	inbuf   chan []byte
	backbuf chan []byte
}

// 入口流量
func (p *Pipe) StreamIn() {
	defer func() {
		log.Println("StreamIn closed")
		close(p.inbuf)
	}()

	buff := make([]byte, bufferSize)

	for {
		sz, err := p.origin.Read(buff)
		if err != nil {
			// 问题直接结束
			fmt.Println(err.Error(), "StreamIn")
			return
		}
		copy_buf := make([]byte, sz)
		copy(copy_buf, buff)
		fmt.Println("reading ", string(copy_buf))
		p.inbuf <- copy_buf
	}
}

// 返回流量
func (p *Pipe) StreamBack() {

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

func (p *Pipe) Working() {
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
				// log.Println("传递", string(buf))
				p.target.Write(buf)
			} else {
				return
			}
		case buf, ok := <-p.backbuf:
			if ok {
				// log.Println("返回", string(buf))
				p.origin.Write(buf)
			} else {
				return
			}
		}
	}

}

func (p *Pipe) Run() {
	go p.Working()
	go p.StreamIn()
	go p.StreamBack()
}
