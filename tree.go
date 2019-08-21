package simpleframework

import (
	"log"
	"net/url"
	"strings"
)

type TreeRouter struct {
	children  []*TreeRouter
	param     byte
	component string
	methods   map[string]HandlerFunc
}

func (t *TreeRouter) Insert(method string, path string, handler HandlerFunc) {
	//分割路径段
	components := strings.Split(path, "/")[1:]
Next:
	//遍历路径段
	for _, component := range components {
		//遍历子节点
		for _, child := range t.children {
			//检查当前子节点是否已存在该中途路径,若存在,指针移动到该子节点上,直接遍历下一个路径段
			if child.component == component {
				t = child
				continue Next
			}
		}
		//创建一个新节点,节点里有当前处理的路径段
		newNode := &TreeRouter{component: component,
			methods: make(map[string]HandlerFunc)}
		//检查当前处理的路径段是否有特殊参数
		if component[0] == ':' || component[0] == '*' {
			newNode.param = component[0]
		}
		//把新节点挂在父节点上
		t.children = append(t.children, newNode)
		//指针移动到新节点
		t = newNode
	}
	t.methods[method] = handler
	log.Println("路径插完了")
}

func (t *TreeRouter) Search(components []string, params url.Values) *TreeRouter {
Next:
	//遍历传入的路径段
	for _, component := range components {
		//遍历已存的子节点
		for _, child := range t.children {
			//只有当传入路径段和子节点存的匹配,或者子节点含特殊参数时,才进行处理和深入,否则一直遍历当前层次的子节点
			if child.component == component || child.param == ':' || child.param == '*' {
				//如果子节点参数是*,直接返回子节点,把其中的方法直接交给net/http包的函数处理
				if child.param == '*' {
					//防止把一些异常参数交给net/http包的函数
					if component == "*" || component == ":" {
						return nil
					} else {
						log.Println("脱离路由树")
						return child
					}
				}
				//如果子节点参数是:,把子节点的路由段加进参数集合,继续下面的代码
				if child.param == ':' {
					params.Add(child.component[1:], component)
				}
				//指针移动到当前处理的子节点
				t = child
				//换下一个路径段
				continue Next
			}
		}
		//当前层次的子节点遍历完了没有找到,404(必须路径段数与节点层数相同,才有可能叶子404,否者提前出第一层for循环)
		log.Println("没找到该叶子路径")
		return nil
	}
	//路径段数小于节点层数会导致提前出第一层for循环,若当前处理的节点没有方法,算做中途404
	if t.methods["GET"] == nil && t.methods["POST"] == nil {
		log.Println("没找到该中途路径")
		return nil
	} else {
		//完全匹配后200
		log.Println("找到路径了")
		return t
	}
}
