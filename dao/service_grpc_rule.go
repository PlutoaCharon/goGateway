package dao

import (
	"github.com/PlutoaCharon/goGateway/public"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

type GrpcRule struct {
	ID             int64  `json:"id" gorm:"primary_key"`
	ServiceID      int64  `json:"service_id" gorm:"column:service_id" description:"服务id	"`
	Port           int    `json:"port" gorm:"column:port" description:"端口	"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor" description:"header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue"`
}

func (t *GrpcRule) TableName() string {
	return "gateway_service_grpc_rule"
}

func (t *GrpcRule) Find(c *gin.Context, tx *gorm.DB, search *GrpcRule) (*GrpcRule, error) {
	model := &GrpcRule{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

func (t *GrpcRule) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error; err != nil {
		return err
	}
	return nil
}
