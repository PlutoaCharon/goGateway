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
)

//AdminRegister admin路由注册
type AdminController struct {
}

func AdminRegister(router *gin.RouterGroup) {
	admin := AdminController{}
	router.GET("/admin_info", admin.AdminInfo)
	router.POST("/change_pwd", admin.ChangePwd)
}

// ListPage godoc
// @Summary 登陆信息
// @Description 登陆信息
// @Tags 管理员接口
// @ID /admin/admin_info
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.AdminInfoOutput} "success"
// @Router /admin/admin_info [get]
// 获取登录管理员信息
func (admin *AdminController) AdminInfo(c *gin.Context) {
	// 1. 读取sessionKey对应json， 转换为结构体
	// 2. 取出数据然后封装 输出
	session := sessions.Default(c)
	adminInfoStr := session.Get(public.AdminInfoSessionKey)
	sessionInfo := &dto.AdminSession{}
	if err := json.Unmarshal([]byte(adminInfoStr.(string)), sessionInfo); err != nil {
		middleware.ResponseError(c, 200, err)
		return
	}
	output := &dto.AdminInfoOutput{
		ID:           sessionInfo.ID,
		Name:         sessionInfo.UserName,
		Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Introduction: "I am a super administrator",
		Roles:        []string{"admin"},
		LoginTime:    sessionInfo.LoginTime,
	}
	middleware.ResponseSuccess(c, output)
	return
}

// ListPage godoc
// @Summary 修改密码
// @Description 修改密码
// @Tags 管理员接口
// @ID /admin/change_pwd
// @Accept  json
// @Produce  json
// @Param body body dto.AdminChangePwdInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin/change_pwd [post]
// 修改密码
func (admin *AdminController) ChangePwd(c *gin.Context) {
	params := &dto.AdminChangePwdInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	// 1. session读取用户信息到结构体 sessInfo
	// 2. sessInfo.ID 读取数据库信息 adminInfo
	// 3. params.password + adminInfo.salt sha256 saltPassword
	// 4. saltPassword ==> adminInfo.password 执行数据保存

	// session读取用户信息到结构体 sessInfo
	session := sessions.Default(c)
	adminInfoStr := session.Get(public.AdminInfoSessionKey)
	sessionInfo := &dto.AdminSession{}
	if err := json.Unmarshal([]byte(adminInfoStr.(string)), sessionInfo); err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	// 读取数据库信息 adminInfo
	search := &dao.GatewayAdmin{
		ID: sessionInfo.ID,
	}
	adminInfo, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	// 生成新密码savePassword 并保存
	savePassword := public.GenSaltPassword(params.Password, adminInfo.Salt)
	adminInfo.Password = savePassword
	if err := adminInfo.Save(c, lib.GORMDefaultPool); err != nil {
		middleware.ResponseError(c, 2004, err)
		return
	}
	middleware.ResponseSuccess(c, "")
	return
}
