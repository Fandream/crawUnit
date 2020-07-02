package routes

import (
	"crawunit/model"
	"crawunit/services"
	"crawunit/services/context"
	"crawunit/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

type UserListRequestData struct {
	PageNum  int `form:"page_num"`
	PageSize int `form:"page_size"`
}

type UserRequestData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Email    string `json:"email"`
}

func Login(c *gin.Context) {
	// 绑定请求数据
	var reqData UserRequestData
	if err := c.ShouldBindJSON(&reqData); err != nil {
		HandleError(http.StatusUnauthorized, c, errors.New("not authorized"))
		return
	}

	// 获取用户
	user, err := model.GetUserByUsername(strings.ToLower(reqData.Username))
	if err != nil {
		HandleError(http.StatusUnauthorized, c, errors.New("not authorized"))
		return
	}

	// 校验密码
	encPassword := utils.EncryptPassword(reqData.Password)
	if user.Password != encPassword {
		HandleError(http.StatusUnauthorized, c, errors.New("not authorized"))
		return
	}

	// 获取token
	tokenStr, err := services.MakeToken(&user)
	if err != nil {
		HandleError(http.StatusUnauthorized, c, errors.New("not authorized"))
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    tokenStr,
	})
}

func GetMe(c *gin.Context) {
	ctx := context.WithGinContext(c)
	user := ctx.User()
	ctx.Success(struct {
		*model.User
		Password string `json:"password,omitempty"`
	}{
		User: user,
	}, nil)
}


