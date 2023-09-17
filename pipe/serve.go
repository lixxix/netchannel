package pipe

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

func StartTCP(beginp string, endp string) {
	ln, err := net.Listen("tcp", beginp)
	if err != nil {
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			return
		}

		log.Println("进入监听")
		// 接入成功
		// 将接入到目标地址
		taddr, err := net.ResolveTCPAddr("tcp", endp)
		if err != nil {
			conn.Close()
			continue
		}
		tconn, err := net.DialTCP("tcp", nil, taddr)
		if err != nil {
			conn.Close()
			continue
		}
		err = tconn.SetKeepAlive(true)
		if err != nil {

			continue
		}

		err = tconn.SetKeepAlivePeriod(time.Second * 30)
		if err != nil {
			continue
		}

		PILE := &Pipe{
			origin: conn,
			target: tconn,
		}
		log.Println("启动")
		PILE.Run()
	}
}

func ServeWs(w http.ResponseWriter, r *http.Request, endp string) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 接入成功
	// 将接入到目标地址
	taddr, err := net.ResolveTCPAddr("tcp", endp)
	if err != nil {
		conn.Close()
		return
	}
	tconn, err := net.DialTCP("tcp", nil, taddr)
	if err != nil {
		conn.Close()
		return
	}
	err = tconn.SetKeepAlive(true)
	if err != nil {

		return
	}

	err = tconn.SetKeepAlivePeriod(time.Second * 30)
	if err != nil {
		return
	}

	PILE := &PipeWS{
		origin: conn,
		target: tconn,
	}
	log.Println("启动WS")
	PILE.Run()
}

func StartWS(beginp string, endp string) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(w, r, endp)
	})

	err := http.ListenAndServe(beginp, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
