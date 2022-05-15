package controller

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/NoCLin/douyin-backend-go/config"
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/gin-gonic/gin"
	"testing"
)

var mock sqlmock.Sqlmock

func TestMain(m *testing.M) {
	fmt.Println("begin")

	config.InitTestConfig()
	gin.SetMode(gin.TestMode)

	mock = G.DBMock
	m.Run()

	fmt.Println("end")
}
