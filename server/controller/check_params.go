package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kubespace/server/global"
)

func CheckParams(c *gin.Context, ptr interface{}) error {
	if ptr == nil {
		return nil
	}
	switch t := ptr.(type) {
	case string:
		if t != "" {
			panic(t)
		}
	case error:
		panic(t.Error())
	}
	if err := c.ShouldBindJSON(&ptr); err != nil {
		global.GVA_LOG.Warn(fmt.Sprintf("解析参数出错：", err))
		return err
	}
	return nil
}
