package simpleframework

import (
	"encoding/json"
	"log"
	"net/http"
)

type Context struct {
	eg       *Engine
	Resp     http.ResponseWriter
	Req      *http.Request
	wsconfig WebSocketConfig
}

func NewContext(w http.ResponseWriter, r *http.Request, eg *Engine,wsconfig WebSocketConfig) *Context {
	c := &Context{}
	c.eg = eg
	c.Init(w, r)
	c.wsconfig=wsconfig
	return c
}

func (c *Context) Init(w http.ResponseWriter, r *http.Request) {
	c.Resp = w
	c.Req = r
}

func (c *Context) String(code int, body string) error {
	c.Resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Resp.WriteHeader(code)
	_, err := c.Resp.Write([]byte(body))
	if err != nil {
		return err
	}
	log.Println("向页面输出文本")
	return nil
}

func (c *Context) JSON(code, status int, message string, data interface{}) error {
	m := map[string]interface{}{
		"message": message,
		"status":  status,
		"data":    data,
	}
	j, err := json.Marshal(m)
	if err != nil {
		return err
	}
	c.Resp.Header().Set("Content-Type", "application/json")
	c.Resp.WriteHeader(code)
	_, err = c.Resp.Write(j)
	if err != nil {
		return err
	}
	log.Println("向页面输出JSON")
	return nil
}

func (c *Context) JSONN(code int, message map[string]interface{}) error {
	j, err := json.Marshal(message)
	if err != nil {
		return err
	}
	c.Resp.Header().Set("Content-Type", "application/json")
	c.Resp.WriteHeader(code)
	_, err = c.Resp.Write(j)
	if err != nil {
		return err
	}
	log.Println("向页面输出JSONN")
	return nil
}
