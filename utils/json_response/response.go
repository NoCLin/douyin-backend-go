package json_response

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/sjson"
	"net/http"
)

type BaseResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

func Response(c *gin.Context, code int32, msg string, data interface{}) {

	// 接口中，数据与 status_code 同一个层级，且无法修改，因此先 Marshal 再 set code & msg
	// 如果 数据统一在下一级，如 {"status_code":0,"status_message":"ok","data":{}}
	// 那么可以直接嵌套一个 interface

	//var m map[string]interface{}
	//if data != nil {
	//	err := mapstructure.Decode(data, &m)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//} else {
	//	m = make(map[string]interface{})
	//}
	marshal, err := json.Marshal(data)
	if err != nil {
		return
	}

	marshal, _ = sjson.SetBytes(marshal, "status_code", code)
	marshal, _ = sjson.SetBytes(marshal, "status_msg", msg)
	c.Data(http.StatusOK, "application/json", marshal)
	//c.JSON(http.StatusOK, marshal)

}

func OK(c *gin.Context, msg string, data interface{}) {
	Response(c, 0, msg, data)
}

func Error(c *gin.Context, code int32, msg string) {
	Response(c, code, msg, nil)
}
