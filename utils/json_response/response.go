package json_response

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"log"
	"net/http"
)

type BaseResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

func Response(c *gin.Context, code int32, msg string, data interface{}) {

	var m map[string]interface{}
	if data != nil {
		err := mapstructure.Decode(data, &m)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		m = make(map[string]interface{})
	}

	m["StatusCode"] = code
	m["StatusMsg"] = msg

	c.JSON(http.StatusOK, m)

}

func OK(c *gin.Context, msg string, data interface{}) {
	Response(c, 0, msg, data)
}

func Error(c *gin.Context, code int32, msg string) {
	Response(c, code, msg, nil)
}
