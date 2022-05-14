package main

import (
	"fmt"
	"github.com/NoCLin/douyin-backend-go/config"
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/gin-gonic/gin"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	addr := G.Config.Server.Addr
	gin.SetMode(G.Config.Server.Mode)

	r := gin.Default()
	initRouter(r)

	err = r.Run(addr)
	if err != nil {
		fmt.Println("启动服务器失败", err)
		panic(err)
	}
}
