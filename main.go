package main

import (
	"bytes"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/godaddy-x/jorm/util"
	scriptDecoder "github.com/metasv/metacontract-script-decoder"
	"net/http"
)

type Message struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type ScriptMode struct {
	Type string `json:"type" bson:"type" binding:"required"`
	Hex  string `json:"hex" bson:"hex" binding:"required"`
}
type MetaId struct {
	Protocol string `json:"protocol"`
	Data     string `json:"data"`
}

func getMetaIdProtocol(pkScript []byte) []byte {
	script := pkScript[:len(pkScript)-5]
	flagTagStart := bytes.IndexByte(script, scriptDecoder.OP_DATA_26)
	flagTagend := bytes.IndexByte(script, scriptDecoder.OP_DATA_2)
	return script[flagTagStart+1 : flagTagend]
}
func getMetaIdFlag(pkScript []byte) []byte {
	script := pkScript[:len(pkScript)-5]
	flagTagStart := bytes.IndexByte(script, scriptDecoder.OP_DATA_6)
	flagTagend := bytes.IndexByte(script, scriptDecoder.OP_DATA_26)
	return script[flagTagStart+1 : flagTagend]
}
func hasMetaIdFlag(pkScript []byte) bool {
	return bytes.Equal(getMetaIdFlag(pkScript), []byte("metaid"))
	//return bytes.Contains(script,[]byte("metaid"))
}

func DecodeMetaId(pkScript []byte, metaId *MetaId) bool {

	ret := false
	if hasMetaIdFlag(pkScript) {
		metaId.Protocol = string(getMetaIdProtocol(pkScript))
		metaId.Data = string(pkScript)
		ret = true
	}

	return ret
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
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(Cors())
	v2 := router.Group("/v1/mvc-browser")
	{
		v2.POST("/script-decoder", Decoder)
	}

	_ = router.Run(util.AddStr("0.0.0.0:", "9030"))

}

func Decoder(context *gin.Context) {
	origin := context.Request.Header.Get("Origin")
	if origin != "https://api-mvc.metasv.com" {
		context.JSONP(http.StatusInternalServerError, Message{Code: 1, Data: "ERROR"})
		return
	}

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
			switch scriptMode.Type {
			case "metaid":
				metaId := &MetaId{}
				DecodeMetaId(script, metaId)
				//log.Printf("%v",metaId)
				context.JSONP(http.StatusOK, gin.H{"code": 0, "data": metaId})
			default:
				txo := &scriptDecoder.TxoData{}

				scriptDecoder.DecodeMvcTxo(script, txo)

				//data, _ := json.Marshal(txo)
				//log.Printf("%v", string(data))
				context.JSONP(http.StatusOK, gin.H{"code": 0, "data": txo})
			}

		}
	}
}
