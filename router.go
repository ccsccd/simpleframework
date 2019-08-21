package simpleframework

import (
	"log"
	"net/http"
	"strings"
)

func (eg *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//从Pool中Get()内容对象,类型断言为*Context
	c := eg.Pool.Get().(*Context)
	log.Printf("c:%p\n",c)
	c.Init(w, r)
	c.Req.ParseForm()
	log.Println("进入路由:"+r.URL.Path)
	node := eg.Tree.Search(strings.Split(r.URL.Path, "/")[1:], c.Req.Form)
	if node != nil && node.methods[r.Method] != nil {
		if err := node.methods[r.Method](c); err != nil {
			log.Println(err.Error())
		}
	} else if node != nil && node.methods[r.Method] == nil {
		if err := eg.MethodNotAllowedHandler(c); err != nil {
			log.Println(err)
		}
	} else {
		if err := eg.NotFoundHandler(c); err != nil {
			log.Println(err)
		}
	}
	//使用完后，我们把item放回池中，让对象可以重用
	eg.Pool.Put(c)
}

func (eg *Engine) Run(addr string) {
	log.Println("准备监听端口")
	if err := http.ListenAndServe(addr, eg); err != nil {
		log.Fatal(err.Error())
	}
}

func NotFoundHandler(c *Context) error {
	http.NotFound(c.Resp, c.Req)
	log.Println("向页面输出404")
	return nil
}

func MethodNotAllowedHandler(ctx *Context) error {
	http.Error(ctx.Resp, "405 method not allowed", 405)
	log.Println("向页面输出405")
	return nil
}

func (eg *Engine) Handle(method string, path string, handler HandlerFunc) {
	eg.Tree.Insert(method, path, handler)
}

func (eg *Engine) GET(path string, handler HandlerFunc) {
	eg.Handle("GET", path, handler)
}

func (eg *Engine) POST(path string, handler HandlerFunc) {
	eg.Handle("POST", path, handler)
}
