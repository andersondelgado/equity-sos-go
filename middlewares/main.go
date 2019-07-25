package middlewares

import (
	json2 "encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"../util"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var datas util.Response
		t := util.DecodeHeaderToken(c)
		if t.Success == true {
			v := reflect.ValueOf(t.Data)
			values := make(map[string]interface{}, v.NumField())
			m := (t.Data)
			fmt.Println("<- ", v.Interface())
			fmt.Println("-> ", m)
			fmt.Println("<-> ", values)

			json, _ := json2.Marshal(&t.Data)
			fmt.Println("<josn> ", string(json))
			//c.Set("user", values)
			//c.Set("user", t.Data)
			c.Set("user", string(json))

			// return next(c)
			c.Next()
		} else {
			datas.Success = t.Success
			datas.Message = t.Message
			datas.Data = ""

			c.AbortWithStatusJSON(http.StatusUnauthorized, datas)
			return
		}

		// return c.JSON(http.StatusUnauthorized, datas)
	}
}
