package controller

import (
	"fmt"
	"github.com/PlutoaCharon/goGateway/dao"
	"github.com/PlutoaCharon/goGateway/dto"
	"github.com/PlutoaCharon/goGateway/middleware"
	"github.com/PlutoaCharon/goGateway/public"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"regexp"
	"strings"
	"time"
)

type ServiceController struct {
}

//ServiceControllerRegister admin路由注册
func ServiceRegister(router *gin.RouterGroup) {
	admin := ServiceController{}
	router.GET("/service_list", admin.ServiceList)               // 列出服务列表具体信息
	router.GET("/service_delete", admin.ServiceDelete)           // 删除指定ID号服务信息
	router.GET("/service_detail", admin.ServiceDetail)           // 查询信息服务信息
	router.POST("/service_add_http", admin.ServiceAddHttp)       // 添加Http服务信息
	router.POST("/service_update_http", admin.ServiceUpdateHttp) // 更新Http服务信息
	router.GET("/service_status", admin.ServiceStatistics)       // 查看流量统计信息
	//router.POST("/service_add_tcp", admin.ServiceAddTcp)         // 添加Tcp服务信息
	//router.POST("/service_update_tcp", admin.ServiceUpdateTcp)   // 更新Tcp服务信息
	//router.POST("/service_add_grpc", admin.ServiceAddGrpc)       // 添加Grpc服务信息
	//router.POST("/service_update_grpc", admin.ServiceUpdateGrpc) // 更新Grpc服务信息

}

// ServiceList godoc
// @Summary 服务列表
// @Description 服务列表
// @Tags 服务管理
// @ID /service/service_list
// @Accept  json
// @Produce  json
// @Param info query string false "关键词"
// @Param page_size query int true "每页个数"
// @Param page_no query int true "当前页数"
// @Success 200 {object} middleware.Response{data=dto.ServiceListItemOutput} "success"
// @Router /service/service_list [get]
func (admin *ServiceController) ServiceList(c *gin.Context) {
	var params = &dto.ServiceListInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	// 从数据库中分页读取信息
	var roleInfo = &dao.ServiceInfo{}
	list, total, err := roleInfo.ServiceList(c, lib.GORMDefaultPool, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	// 格式化输出信息
	var outputList []dto.ServiceListItemOutput
	for _, item := range list {
		detail, err := item.ServiceDetail(c, lib.GORMDefaultPool, &dao.ServiceInfo{
			ID: item.ID,
		})
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			return
		}
		// 1. http后缀接入 cluserIp + cluserPort+Path
		// 2. http域名接入 domain
		// 3. tcp, grpc cluserIp + srevicePort

		// 获取总结点数
		totalNode := len(detail.LoadBalance.GetIPListByModel())

		// 获取qps， qpd
		//serviceCounter, _ := public.FlowCounterHandler.GetCounter(public.FlowServicePrefix + item.ServiceName)
		//qps := serviceCounter.GetQPS()
		//qpd, _ := serviceCounter.GetDayCount(time.Now())

		serviceIP := lib.GetStringConf("base.cluster.cluster_ip")
		servicePort := lib.GetStringConf("base.cluster.cluster_port")
		serviceSSLPort := lib.GetStringConf("base.cluster.cluster_ssl_port")
		serviceAddr := "unknow"

		serviceHttpBaseURL := serviceIP + ":" + servicePort
		if detail.HttpRule.NeedHttps == 1 {
			serviceHttpBaseURL = serviceIP + ":" + serviceSSLPort
		}
		if item.LoadType == public.LoadTypeHTTP && detail.HttpRule.RuleType == 0 {
			serviceAddr = fmt.Sprintf("%s%s", serviceHttpBaseURL, detail.HttpRule.Rule)
		}
		if item.LoadType == public.LoadTypeHTTP && detail.HttpRule.RuleType == 1 {
			serviceAddr = detail.HttpRule.Rule
		}
		if item.LoadType == public.LoadTypeTCP {
			serviceAddr = fmt.Sprintf("%s:%d", serviceIP, detail.TcpRule.Port)
		}
		if item.LoadType == public.LoadTypeGRPC {
			serviceAddr = fmt.Sprintf("%s:%d", serviceIP, detail.GrpcRule.Port)
		}
		outputList = append(outputList, dto.ServiceListItemOutput{
			ID:          item.ID,
			LoadType:    item.LoadType,
			ServiceName: item.ServiceName,
			ServiceDesc: item.ServiceDesc,
			UpdatedAt:   item.UpdatedAt,
			CreatedAt:   item.CreatedAt,
			//QPS:         qps,
			//QPD:         qpd,
			QPS:         0,
			QPD:         0,
			TotalNode:   totalNode,
			ServiceAddr: serviceAddr,
		})
	}
	middleware.ResponseSuccess(c, map[string]interface{}{"list": outputList, "total": total})
	return
}

// ServiceDelete godoc
// @Summary 服务删除
// @Description 服务删除
// @Tags 服务管理
// @ID /service/service_delete
// @Accept  json
// @Produce  json
// @Param id query string true "服务ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_delete [get]
func (admin *ServiceController) ServiceDelete(c *gin.Context) {
	params := &dto.ServiceDetailInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	// 读取服务基本信息
	search := &dao.ServiceInfo{
		ID: params.ID,
	}

	info, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	info.IsDelete = 1
	if err := info.Save(c, lib.GORMDefaultPool); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")
	return
}

// ServiceStatistics godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags 服务管理
// @ID /service/service_status
// @Accept  json
// @Produce  json
// @Param id query string true "服务ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_status [get]
func (admin *ServiceController) ServiceStatistics(c *gin.Context) {
	params := &dto.ServiceDetailInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	search := &dao.ServiceInfo{
		ID: params.ID,
	}
	detail, err := search.ServiceDetail(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	counter, _ := public.FlowCounterHandler.GetCounter(public.FlowServicePrefix + detail.Info.ServiceName)

	//今日流量全天小时级访问统计
	var todayStat []int64
	for i := 0; i <= time.Now().In(lib.TimeLocation).Hour(); i++ {
		nowTime := time.Now()
		nowTime = time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		hourStat, _ := counter.GetHourCount(nowTime)
		todayStat = append(todayStat, hourStat)
	}

	//昨日流量全天小时级访问统计
	var yesterdayStat []int64
	for i := 0; i <= 23; i++ {
		nowTime := time.Now().AddDate(0, 0, -1)
		nowTime = time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		hourStat, _ := counter.GetHourCount(nowTime)
		yesterdayStat = append(yesterdayStat, hourStat)
	}
	//yesterdayStat = []int64{
	//	120,
	//	50,
	//	10,
	//	57,
	//	59,
	//	48,
	//	76,
	//	69,
	//	200,
	//	400,
	//	580,
	//	1500,
	//	2500,
	//	2300,
	//	1300,
	//	1700,
	//	1900,
	//	1000,
	//	800,
	//	570,
	//	500,
	//	360,
	//	200,
	//	105,
	//}
	//todayStat = []int64{
	//	78,
	//	23,
	//	78,
	//	123,
	//	325,
	//	378,
	//	456,
	//	478,
	//	500,
	//	800,
	//	760,
	//}
	middleware.ResponseSuccess(c, map[string][]int64{
		"today":     todayStat,
		"yesterday": yesterdayStat,
	})
	return
}

// ServiceAddHttp godoc
// @Summary 添加Http服务信息
// @Description 添加Http服务信息
// @Tags 服务管理
// @ID /service/service_add_http
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceAddHttpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_add_http [post]
// 添加Http服务信息
func (admin *ServiceController) ServiceAddHttp(c *gin.Context) {
	params := &dto.ServiceAddHttpInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	//验证service_name是否被占用
	search := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
	}
	if _, err := search.Find(c, lib.GORMDefaultPool, search); err == nil {
		middleware.ResponseError(c, 2001, errors.New("服务名被占用，请重新输入"))
		return
	}

	//验证rule前缀是否被占用
	ruleSearch := &dao.HttpRule{
		Rule: params.Rule,
	}
	if _, err := ruleSearch.Find(c, lib.GORMDefaultPool, ruleSearch); err == nil {
		middleware.ResponseError(c, 2002, errors.New("服务前缀或域名被占用，请重新输入"))
		return
	}

	//验证rule_type=0时以/开头，rule_type=1时不能出现/
	if params.RuleType == 0 {
		matched, _ := regexp.Match(`^/\S+$`, []byte(params.Rule))
		if !matched {
			middleware.ResponseError(c, 2003, errors.New("路径接入时必须以/开头"))
			return
		}
	}
	if params.RuleType == 1 {
		matched, _ := regexp.Match(`^[0-9a-z-_\.]+$`, []byte(params.Rule))
		if !matched {
			middleware.ResponseError(c, 2004, errors.New("域名接入时，只支持数字、小写字母、中划线、下划线"))
			return
		}
	}

	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2005, errors.New("ip列表与权重设置不匹配"))
		return
	}

	tx := lib.GORMDefaultPool.Begin()
	info := &dao.ServiceInfo{
		LoadType:    public.LoadTypeHTTP,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2002, err)
		return
	}

	loadBalance := &dao.LoadBalance{
		ServiceID:              info.ID,
		RoundType:              params.RoundType,
		IpList:                 params.IpList,
		WeightList:             params.WeightList,
		ForbidList:             params.ForbidList,
		UpstreamConnectTimeout: params.UpstreamConnectTimeout,
		UpstreamHeaderTimeout:  params.UpstreamHeaderTimeout,
		UpstreamIdleTimeout:    params.UpstreamIdleTimeout,
		UpstreamMaxIdle:        params.UpstreamMaxIdle,
		//CheckMethod:            public.DefaultCheckMethod,
		//CheckTimeout:           public.DefaultCheckTimeout,
		//CheckInterval:          public.DefaultCheckInterval,
	}
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2003, err)
		return
	}

	httpRule := &dao.HttpRule{
		ServiceID:      info.ID,
		RuleType:       params.RuleType,
		Rule:           params.Rule,
		NeedHttps:      params.RuleType,
		NeedWebsocket:  params.RuleType,
		NeedStripUri:   params.NeedStripUri,
		UrlRewrite:     params.UrlRewrite,
		HeaderTransfor: params.HeaderTransfor,
	}
	if err := httpRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2003, err)
		return
	}

	accessControl := &dao.AccessControl{
		ServiceID:         info.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		WhiteHostName:     params.WhiteHostName,
		ClientIPFlowLimit: params.ClientIPFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2004, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "")
	return
}

// 修改Http服务信息
func (admin *ServiceController) ServiceUpdateHttp(c *gin.Context) {
	params := &dto.ServiceUpdateHttpInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	//验证rule_type=0时以/开头，rule_type=1时不能出现/
	if params.RuleType == 0 {
		matched, _ := regexp.Match(`^/\S+$`, []byte(params.Rule))
		if !matched {
			middleware.ResponseError(c, 2002, errors.New("路径接入时必须以/开头"))
			return
		}
	}
	if params.RuleType == 1 {
		matched, _ := regexp.Match(`^[0-9a-z-_\.]+$`, []byte(params.Rule))
		if !matched {
			middleware.ResponseError(c, 2002, errors.New("域名接入时，只支持数字、小写字母、中划线、下划线"))
			return
		}
	}

	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2002, errors.New("ip列表与权重设置不匹配"))
		return
	}

	tx := lib.GORMDefaultPool.Begin()
	service := &dao.ServiceInfo{
		ID: params.ID,
	}
	detail, err := service.ServiceDetail(c, lib.GORMDefaultPool, service)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	info := detail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2002, err)
		return
	}

	loadBalance := &dao.LoadBalance{}
	if detail.LoadBalance != nil {
		loadBalance = detail.LoadBalance
	}
	fmt.Println("params.UpstreamConnectTimeout", params.UpstreamConnectTimeout)
	loadBalance.ServiceID = info.ID
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.ForbidList = params.ForbidList
	loadBalance.UpstreamConnectTimeout = params.UpstreamConnectTimeout
	loadBalance.UpstreamHeaderTimeout = params.UpstreamHeaderTimeout
	loadBalance.UpstreamIdleTimeout = params.UpstreamIdleTimeout
	loadBalance.UpstreamMaxIdle = params.UpstreamMaxIdle
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2003, err)
		return
	}

	httpRule := &dao.HttpRule{}
	if detail.HttpRule != nil {
		httpRule = detail.HttpRule
	}
	httpRule.ServiceID = info.ID
	httpRule.RuleType = params.RuleType
	httpRule.Rule = params.Rule
	httpRule.NeedHttps = params.NeedHttps
	httpRule.NeedWebsocket = params.NeedWebsocket
	httpRule.NeedStripUri = params.NeedStripUri
	httpRule.UrlRewrite = params.UrlRewrite
	httpRule.HeaderTransfor = params.HeaderTransfor
	httpRule.HeaderTransfor = params.HeaderTransfor
	httpRule.HeaderTransfor = params.HeaderTransfor
	httpRule.HeaderTransfor = params.HeaderTransfor
	httpRule.HeaderTransfor = params.HeaderTransfor
	httpRule.HeaderTransfor = params.HeaderTransfor
	if err := httpRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2003, err)
		return
	}

	accessControl := &dao.AccessControl{}
	if detail.AccessControl != nil {
		accessControl = detail.AccessControl
	}
	accessControl.ServiceID = info.ID
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.WhiteHostName = params.WhiteHostName
	accessControl.ClientIPFlowLimit = params.ClientIPFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2004, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "")
	return
}

// ServiceDetail godoc
// @Summary 服务详细
// @Description 服务详细
// @Tags 服务管理
// @ID /service/service_detail
// @Accept  json
// @Produce  json
// @Param id query dto.ServiceDetailInput true "服务ID"
// @Success 200 {object} middleware.Response{data=dao.ServiceDetail} "success"
// @Router /service/service_detail [get]
// 查询服务详细
func (admin *ServiceController) ServiceDetail(c *gin.Context) {
	params := &dto.ServiceDetailInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	search := &dao.ServiceInfo{
		ID: params.ID,
	}
	detail, err := search.ServiceDetail(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	middleware.ResponseSuccess(c, detail)
	return
}
