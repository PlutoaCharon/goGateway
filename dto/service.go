package dto

import (
	"github.com/PlutoaCharon/goGateway/public"
	"github.com/gin-gonic/gin"
	"time"
)

// 具体信息展示
type ServiceListItemOutput struct {
	ID          int64     `json:"id" gorm:"primary_key"`
	LoadType    int       `json:"load_type" gorm:"column:load_type" description:"账号创建时间"`
	ServiceAddr string    `json:"service_addr" gorm:"column:service_addr" description:"服务地址"`
	ServiceName string    `json:"service_name" gorm:"column:service_name" description:"服务名称"`
	ServiceDesc string    `json:"service_desc" gorm:"column:service_desc" description:"服务描述"`
	UpdatedAt   time.Time `json:"create_at" gorm:"column:create_at" description:"添加时间	"`
	CreatedAt   time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	QPS         int64     `json:"qps" description:"每秒请求数"`
	QPD         int64     `json:"qpd" description:"每天请求数"`
	TotalNode   int       `json:"total_node" description:"总节点数"`
}

//---------------------------------------------------------------------------------------------------------------

// 搜索框查找信息参数
type ServiceListInput struct {
	Info     string `json:"info" form:"info" comment:"查找信息" validate:""`
	PageSize int    `json:"page_size" form:"page_size" comment:"页数" validate:"required,min=1,max=999"`
	PageNo   int    `json:"page_no" form:"page_no" comment:"页码" validate:"required,min=1,max=999"`
}

func (params *ServiceListInput) GetValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

//---------------------------------------------------------------------------------------------------------------

// 删除信息根据ID
type ServiceDetailInput struct {
	ID int64 `json:"id" form:"id" comment:"服务ID" validate:"required"`
}

func (params *ServiceDetailInput) GetValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}

//---------------------------------------------------------------------------------------------------------------
