package simpleframework

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type WebSocketConfig struct {
	OnOpen    func(*WebSocketContext) error
	OnClose   func(*WebSocketContext) error
	OnMessage func(*WebSocketContext, string) error
	OnError   func(*WebSocketContext) error
}

type WebSocketContext struct {
	Conn   net.Conn
	buffer *bufio.ReadWriter
	config WebSocketConfig
}

func (eg *Engine) WebSocket(path string, wsconfig WebSocketConfig) {
	eg.Handle("GET", path, Open)
	eg.wsconfig = wsconfig
}

func Open(c *Context) error {
	r := c.Req
	w := c.Resp
	config := c.wsconfig

	key := r.Header.Get("Sec-WebSocket-Key")
	s := sha1.New()
	s.Write([]byte(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	b := s.Sum(nil)
	accept := base64.StdEncoding.EncodeToString(b)

	hijack := w.(http.Hijacker)
	con, buffer, _ := hijack.Hijack()

	wsc := &WebSocketContext{
		Conn:   con,
		buffer: buffer,
		config: config,
	}

	upgrade := "HTTP/1.1 101 Switching Protocols\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Accept: " + accept + "\r\n\r\n"
	_, err := buffer.Write([]byte(upgrade))
	if err != nil {
		log.Println(err)
	}
	_ = buffer.Flush()

	_ = config.OnOpen(wsc)

	for {
		_ = wsc.config.OnMessage(wsc, wsc.Read())
	}
}
func Close(wsc *WebSocketContext) {
	err := wsc.config.OnClose(wsc)
	if err != nil {
		log.Println(err)
	}
	_ = wsc.Conn.Close()
}

func (wsc *WebSocketContext) Read() string {
	data := make([]byte, 2)
	_, err := wsc.buffer.Read(data)
	if err != nil {
		Close(wsc)
	}

	bin1 := DecBin(int(data[0]))
	bin2 := DecBin(int(data[1]))

	// bin1
	// 0    1    2     3   4 5 6 7
	// FIN RSV1 RSV2 RSV3 opcode(4)

	// bin2
	// 8     9 10 11 12 13 14 15
	// MASK      PayloadLen

	// RSV
	if bin1[1] == 1 || bin1[2] == 1 || bin1[3] == 1 {
		Close(wsc)
	}
	if bin2[0] == 0 {
		Close(wsc)
	}

	opcode := BinDec(bin1[4:])
	payloadLen := BinDec(bin2[1:])

	switch opcode {
	case 1:
		maskingKey := make([]byte, 4)
		_, _ = wsc.buffer.Read(maskingKey)

		payload := make([]byte, payloadLen)
		_, _ = wsc.buffer.Read(payload)

		data := make([]byte, payloadLen)
		for i := 0; i < payloadLen; i++ {
			data[i] = payload[i] ^ maskingKey[i%4]
		}
		return string(data)
	default:
		Close(wsc)
		return ""
	}
}

func (wsc *WebSocketContext) Send(str string, opcode int) {
	dataLen := len(str)

	bin1 := 0x80 | byte(opcode)
	_, err := wsc.Conn.Write([]byte{bin1})
	if err != nil {
		log.Println(err.Error())
	}

	var bin2 byte
	switch {
	case dataLen <= 125:
		bin2 = byte(dataLen)
	default:
		Close(wsc)
	}
	_, err = wsc.Conn.Write([]byte{bin2})

	_, err = wsc.Conn.Write([]byte(str))
	log.Println("该数据帧的真实负载数据长度(bytes):", dataLen)
	log.Println("服务端向客户端发送了:", str)
}

func DecBin(i int) string {
	if i < 0 {
		return ""
	}
	if i == 0 {
		return "0"
	}
	s := ""
	for q := i; q > 0; q = q / 2 {
		m := q % 2
		s = fmt.Sprintf("%v%v", m, s)
	}
	return s
}

func BinDec(b string) int {
	s := strings.Split(b, "")
	l := len(s)
	i := 0
	d := float64(0)
	for i = 0; i < l; i++ {
		f, err := strconv.ParseFloat(s[i], 10)
		if err != nil {
			log.Println(err.Error())
			return -1
		}
		d += f * math.Pow(2, float64(l-i-1))
	}
	return int(d)
}
