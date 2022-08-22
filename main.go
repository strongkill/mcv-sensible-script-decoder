package main

import (
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/godaddy-x/jorm/util"
	scriptDecoder "github.com/metasv/metacontract-script-decoder"
	"log"
	"net/http"
)

type Message struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type ScriptMode struct {
	Hex string `json:"hex" bson:"hex" binding:"required"`
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization,authKey")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func main() {

	router := gin.Default()
	router.Use(Cors())
	v2 := router.Group("/v1/mvc-browser")
	{
		v2.POST("/script-decoder", Decoder)
	}

	_ = router.Run(util.AddStr("0.0.0.0:", "3000"))

}

func Decoder(context *gin.Context) {
	var scriptMode ScriptMode
	if err := context.ShouldBindJSON(&scriptMode); err != nil {
		context.JSONP(http.StatusInternalServerError, Message{Code: 1, Data: "Request params is empty"})
	} else {
		if scriptMode.Hex == "" {
			context.JSONP(http.StatusInternalServerError, Message{Code: 1, Data: "Script is empty"})
		} else {
			script, err := hex.DecodeString(scriptMode.Hex)
			if err != nil {
				return
			}

			txo := &scriptDecoder.TxoData{}

			scriptDecoder.DecodeMvcTxo(script, txo)

			data, _ := json.Marshal(txo)
			log.Printf("%v", string(data))
			context.JSONP(http.StatusOK, gin.H{"code": 0, "data": txo})
		}
	}
}
