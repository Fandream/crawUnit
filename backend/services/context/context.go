package context

import (
	"crawlab/constants"
	"crawlab/errors"
	"crawlab/model"
	"fmt"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	errors2 "github.com/pkg/errors"
	"net/http"
	"runtime/debug"
)

type Context struct {
	*gin.Context
}

func (c *Context) User() *model.User {
	userIfe, exists := c.Get(constants.ContextUser)
	if !exists {
		return nil
	}
	user, ok := userIfe.(*model.User)
	if !ok {
		return nil
	}
	return user
}
func (c *Context) Success(data interface{}, metas ...interface{}) {
	var meta interface{}
	if len(metas) == 0 {
		meta = gin.H{}
	} else {
		meta = metas[0]
	}
	if data == nil {
		data = gin.H{}
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "success",
		"data":    data,
		"meta":    meta,
		"error":   "",
	})
}
func (c *Context) Failed(err error, variables ...interface{}) {
	c.failed(err, http.StatusOK, variables...)
}
