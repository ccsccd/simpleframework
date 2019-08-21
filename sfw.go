package simpleframework

import (
	"sync"
)

type HandlerFunc func(*Context) error

type Engine struct {
	//路由树
	Tree *TreeRouter
	//内容池，复用临时对象，减少内存压力
	Pool                    *sync.Pool
	NotFoundHandler         HandlerFunc
	MethodNotAllowedHandler HandlerFunc
	wsconfig WebSocketConfig
}

func Default() *Engine {
	tree := &TreeRouter{
		component: "/",
		methods:   make(map[string]HandlerFunc),
	}
	eg := &Engine{Tree: tree,
		NotFoundHandler:         NotFoundHandler,
		MethodNotAllowedHandler: MethodNotAllowedHandler,
	}
	//New()函数的作用是当我们从Pool中Get()对象时，如果Pool为空，则先通过New创建一个对象，插入Pool中，然后返回对象
	eg.Pool = &sync.Pool{
		New: func() interface{} {
			return NewContext(nil, nil, eg,eg.wsconfig)
		},
	}
	return eg
}
