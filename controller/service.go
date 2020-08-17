package controller

import (
	"github.com/PlutoaCharon/goGateway/dao"
	"github.com/PlutoaCharon/goGateway/dto"
	"github.com/PlutoaCharon/goGateway/middleware"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"time"
)

type ServiceController struct {
}

//ServiceControllerRegister admin路由注册
func ServiceRegister(router *gin.RouterGroup) {
	admin := ServiceController{}
	router.GET("/service_list", admin.ServiceList)
}

func (admin *ServiceController) ServiceList(c *gin.Context) {
	params := &dto.ServiceListInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	roleInfo := &dao.ServiceInfo{}
	list, total, err := roleInfo.ServiceList(c, lib.GORMDefaultPool, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	outputList := []dto.ServiceListItemOutput{}
	for _, item := range list {
		detail, err := item.ServiceDetail(c, lib.GORMDefaultPool, &dao.ServiceInfo{
			ID: item.ID,
		})
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			return
		}

		totalNode := len(detail.LoadBalance.GetIPListByModel())
		serviceCounter, _ := public.FlowCounterHandler.GetCounter(public.FlowServicePrefix + item.ServiceName)
		qps := serviceCounter.GetQPS()
		qpd, _ := serviceCounter.GetDayCount(time.Now())

		serviceIP := lib.GetStringConf("base.cluster.cluster_ip")
		servicePort := lib.GetStringConf("base.cluster.cluster_port")
		serviceSSLPort := lib.GetStringConf("base.cluster.cluster_ssl_port")
		serviceHttpBaseURL := serviceIP + ":" + servicePort
		if detail.HttpRule.NeedHttps == 1 {
			serviceHttpBaseURL = serviceIP + ":" + serviceSSLPort
		}
		serviceAddr := "unknow"
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
			QPS:         qps,
			QPD:         qpd,
			TotalNode:   totalNode,
			ServiceAddr: serviceAddr,
		})
	}
	middleware.ResponseSuccess(c, map[string]interface{}{"list": outputList, "total": total})
	return
}
