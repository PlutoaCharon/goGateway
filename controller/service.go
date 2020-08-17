package controller

import (
	"fmt"
	"github.com/PlutoaCharon/goGateway/dao"
	"github.com/PlutoaCharon/goGateway/dto"
	"github.com/PlutoaCharon/goGateway/middleware"
	"github.com/PlutoaCharon/goGateway/public"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
)

type ServiceController struct {
}

//ServiceControllerRegister admin路由注册
func ServiceRegister(router *gin.RouterGroup) {
	admin := ServiceController{}
	router.GET("/service_list", admin.ServiceList)     // 列出服务列表具体信息
	router.GET("/service_delete", admin.ServiceDelete) // 删除指定ID号服务信息
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
// @Summary 服务列表
// @Description 服务列表
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
