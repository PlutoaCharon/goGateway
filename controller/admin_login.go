package controller

import (
	"encoding/json"
	"github.com/PlutoaCharon/goGateway/dao"
	"github.com/PlutoaCharon/goGateway/dto"
	"github.com/PlutoaCharon/goGateway/middleware"
	"github.com/PlutoaCharon/goGateway/public"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"time"
)

type AdminLogin struct {
}

func AdminLoginRegister(router *gin.RouterGroup) {
	admin := &AdminLogin{}
	router.POST("/login", admin.Login)
	router.POST("/logout", admin.Logout)
}

// ListPage godoc
// @Summary 登录接口
// @Description 登录接口
// @Tags 管理员接口
// @ID /admin_login/login
// @Accept  json
// @Produce  json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin_login/login [post]
func (AdminLogin *AdminLogin) Login(c *gin.Context) {

	// 1. params.UserName 取得管理员信息 admininfo
	// 2. admininfo.salt + params.Password sha256 => saltPassword
	// 3. saltPassword == admininfo.password
	params := &dto.AdminLoginInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	adminInfo, err := (&dao.GatewayAdmin{}).LoginCheck(c, lib.GORMDefaultPool, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	adminSession := &dto.AdminSession{
		ID:        adminInfo.ID,
		LoginTime: time.Now(),
		UserName:  adminInfo.UserName,
	}
	session := sessions.Default(c)
	adminBts, _ := json.Marshal(adminSession)
	session.Set(public.AdminInfoSessionKey, string(adminBts))
	_ = session.Save()
	output := &dto.AdminLoginOutput{Token: adminInfo.UserName}
	middleware.ResponseSuccess(c, output)
	return
}

// Logout godoc
// @Summary 退出接口
// @Description 退出接口
// @Tags 管理员接口
// @ID /admin_login/logout
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin_login/logout [post]
func (AdminLogin *AdminLogin) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	_ = session.Save()
	middleware.ResponseSuccess(c, "")
	return
}
