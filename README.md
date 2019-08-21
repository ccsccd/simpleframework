# simpleframework
### 使用说明
+ 创建路由
```
router := simpleframework.Default()`
```
+ GET方法，用文本输出的工具方法验证：
```
router.GET("/get/h/i/hello", func(c *simpleframework.Context) error {
	return c.String(200, "get——hello")
})
```
+ POST方法：
```
router.POST("/post", func(c *simpleframework.Context) error {
	return c.String(200, "post——hello")
})
```
+ 用simpleframework.Context取其中的Resp和Req，读取静态资源：
```
//要把工作路径改为“$GOPATH/src/其他分级文件夹/simpleftamwork/demo”,包路径改为“其他分级文件夹/simpleftamwork/demo”
router.GET("/html/*", func(c *simpleframework.Context) error {
	staticHandle := http.StripPrefix("/html",
		http.FileServer(http.Dir("./static/html")))
	staticHandle.ServeHTTP(c.Resp, c.Req)
	return nil
})
router.GET("/image/*", func(c *simpleframework.Context) error {
	staticHandle := http.StripPrefix("/image",
		http.FileServer(http.Dir("./static/image")))
	staticHandle.ServeHTTP(c.Resp, c.Req)
	return nil
})
```
+ 两种方式传参：
```
//直接取simpleframework.Context中Req的参数
router.GET("/hello", func(c *simpleframework.Context) error {
	name := c.Req.FormValue("name")
	if name==""{
		return c.String(200,"hello，guest")
	}else {
		return c.String(200,"hello，"+name)
	}
})
//路径传参(动态路由)，同样从simpleframework.Context中Req取参数
router.GET("/hello/:name", func(c *simpleframework.Context) error {
	return c.String(200,"hello，"+c.Req.Form["name"][0])
})
```
+ json输出的工具方法
```
//获取json方法1
router.GET("/json", func(c *simpleframework.Context) error {
	return c.JSON(200, 200, "this is a message1", map[string]string{"data1": "demo1", "data2": "demo2"})
})
//获取json方法2
router.GET("/jsonn", func(c *simpleframework.Context) error {
	return c.JSONN(200, map[string]interface{}{"message": "this is a message2", "status": 200})
})
```
+ 极简的WebSocketDemo
```
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
```
+ 最后监听端口
```
router.Run(":80")
```
### 备注
+ 暂不支持GET/POST以外的方法，路由树里未添加
+ 极简WebSocket只能每帧传小于125字节的文本
+ WebSocket和GET共用路由树；不能出现重复路径，否者后者覆盖前者
