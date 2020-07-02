package services

import (
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

func GetAuthQuery(query bson.M, c *gin.Context) bson.M {

		// 只获取自己的数据
		query["user_id"] = 1
		return query
}

