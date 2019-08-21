package main

import (
	"app/simpleframework"
	"log"
	"net/http"
)

func main() {
	router := simpleframework.Default()

	//get方法测试
	router.GET("/get/h/i/hello", func(c *simpleframework.Context) error {
		return c.String(200, "get——hello")
	})
	router.GET("/get/hi", func(c *simpleframework.Context) error {
		return c.String(200, "get——hi")
	})
	//post方法测试
	router.POST("/post", func(c *simpleframework.Context) error {
		return c.String(200, "post——hello")
	})
	//获取静态页面
	//要把工作路径改为“$GOPATH/src/其他分级文件夹/simpleftamwork/demo”,包路径改为“其他分级文件夹/simpleftamwork/demo”
	router.GET("/html/*", func(c *simpleframework.Context) error {
		staticHandle := http.StripPrefix("/html",
			http.FileServer(http.Dir("./static/html")))
		staticHandle.ServeHTTP(c.Resp, c.Req)
		return nil
	})
	//获取图片
	//路径要求同上
	router.GET("/image/*", func(c *simpleframework.Context) error {
		staticHandle := http.StripPrefix("/image",
			http.FileServer(http.Dir("./static/image")))
		staticHandle.ServeHTTP(c.Resp, c.Req)
		return nil
	})
	//传参
	router.GET("/hello", func(c *simpleframework.Context) error {
		name := c.Req.FormValue("name")
		if name==""{
			return c.String(200,"hello，guest")
		}else {
			return c.String(200,"hello，"+name)
		}
	})
	//路径传GBK参(动态路由)
	router.GET("/hello/:name", func(c *simpleframework.Context) error {
		return c.String(200,"hello，"+c.Req.Form["name"][0])
	})
	//获取json方法1
	router.GET("/json", func(c *simpleframework.Context) error {
		return c.JSON(200, 200, "this is a message1", map[string]string{"data1": "demo1", "data2": "demo2"})
	})
	//获取json方法2
	router.GET("/jsonn", func(c *simpleframework.Context) error {
		return c.JSONN(200, map[string]interface{}{"message": "this is a message2", "status": 200})
	})
	//WebSocketDemo仅限低字节文本
	router.WebSocket("/ws", simpleframework.WebSocketConfig{
		OnOpen: func(wsc *simpleframework.WebSocketContext) error{
			log.Println("ws:open!")
			return nil
		},
		OnClose: func(wsc *simpleframework.WebSocketContext) error {
			log.Println("ws:close!")
			return nil
		},
		OnMessage: func(wsc *simpleframework.WebSocketContext, s string) error {
			log.Println(wsc.Conn.RemoteAddr(),"客户端发送了信息:", s)
			wsc.Send(s,1)
			return nil
		},
		OnError: nil,
	})

	router.Run(":80")
}
