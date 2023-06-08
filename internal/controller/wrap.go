package controller

import (
	"ai-smart/internal/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"reflect"
	"time"
)

func GinHandelWrap2(obj interface{}) gin.HandlerFunc {
	f := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)

	// check func
	// 检查类型是否是方法
	if f.Kind() != reflect.Func {
		panic("not a func.")
	}
	// 检查参数的数量为1个或2个
	if typ.NumIn() > 2 || typ.NumIn() < 1 {
		panic("func in num not equal 1 or 2.")
	}
	// 检查第一个参数是否为 *gin.Context
	if typ.In(0) != reflect.TypeOf(&gin.Context{}) {
		panic("func first param not *gin.Context.")
	}
	// 如果有第二个参数，判断第二个参数类型是否为指针
	if typ.NumIn() == 2 && typ.In(1).Kind() != reflect.Ptr {
		panic("func second param is not a ptr.")
	}
	// 判断返回参数个数
	if typ.NumOut() != 1 {
		panic("func out num not equal 1.")
	}

	// 获取 *model.BaseResponseInterface 的类型
	tp1 := reflect.TypeOf((*model.BaseResponseInterface)(nil)).Elem()
	// 检查返回的参数是否实现了，model.BaseResponseInterface接口
	if !typ.Out(0).Implements(tp1) {
		panic("func out param not base response.")
	}

	return func(c *gin.Context) {
		in := []reflect.Value{reflect.ValueOf(c)}
		var req interface{}

		// 解析请求
		if typ.NumIn() == 2 {
			// 获得第二个参数的类型
			secondType := typ.In(1)
			// 值类型为结构体
			tmp := reflect.New(secondType.Elem()).Interface()
			if err := getRequest(c, tmp); err != nil {
				log.Fatalf("wrapper get request err:%v", err)
			}
			// 鉴定参数
			if err := validator.New().Struct(tmp); err != nil {
				c.AbortWithStatusJSON(http.StatusOK, model.ParamErrRsp)
				return
			}
			log.Printf("get req:%+v,type:%T", tmp, tmp)
			req = tmp
			in = append(in, reflect.ValueOf(tmp))
		}

		begin := time.Now()
		ans := f.Call(in)[0].Interface()
		c.JSON(http.StatusOK, ans)
		cost := time.Since(begin)

		log.Printf("uri:%v request:%v response:%v cost:%v", c.Request.RequestURI, req, ans, cost)
		if v, ok := ans.(model.BaseResponseInterface); ok {
			if v.GetErr() != nil {
				log.Fatalf("when deal uri:%v,req:%v,appear err:%+v", c.Request.RequestURI, req, v.GetErr())
			}
		}
	}
}

func getRequest(c *gin.Context, req interface{}) error {
	if c.Request.Method == http.MethodPost {
		body, err := c.GetRawData()
		if err != nil {
			return err
		}
		return json.Unmarshal(body, req)
	} else if c.Request.Method == http.MethodGet {
		return c.ShouldBindQuery(req)
	}
	return nil
}
