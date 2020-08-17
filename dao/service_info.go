package dao

import (
	"github.com/PlutoaCharon/goGateway/dto"
	"github.com/PlutoaCharon/goGateway/public"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"time"
)

type ServiceInfo struct {
	ID          int64     `json:"id" gorm:"primary_key"`
	LoadType    int       `json:"load_type" gorm:"column:load_type" description:"负载类型 0=http 1=tcp 2=grpc"`
	ServiceName string    `json:"service_name" gorm:"column:service_name" description:"服务名称"`
	ServiceDesc string    `json:"service_desc" gorm:"column:service_desc" description:"服务描述"`
	UpdatedAt   time.Time `json:"create_at" gorm:"column:create_at" description:"添加时间	"`
	CreatedAt   time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	IsDelete    int8      `json:"is_delete" gorm:"column:is_delete" description:"是否已删除；0：否；1：是"`
}

func (t *ServiceInfo) TableName() string {
	return "gateway_service_info"
}

func (t *ServiceInfo) ServiceDetail(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceDetail, error) {
	info := &ServiceInfo{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(info).Error
	if err != nil {
		return nil, err
	}
	httpRule, _ := (&HttpRule{}).Find(c, tx, (&HttpRule{ServiceID: search.ID}))
	tcpRule, _ := (&TcpRule{}).Find(c, tx, (&TcpRule{ServiceID: search.ID}))
	grpcRule, _ := (&GrpcRule{}).Find(c, tx, (&GrpcRule{ServiceID: search.ID}))
	loadbalance, _ := (&LoadBalance{}).Find(c, tx, (&LoadBalance{ServiceID: search.ID}))
	accessControl, _ := (&AccessControl{}).Find(c, tx, (&AccessControl{ServiceID: search.ID}))
	detail := &ServiceDetail{
		Info:          info,
		HttpRule:      httpRule,
		TcpRule:       tcpRule,
		GrpcRule:      grpcRule,
		LoadBalance:   loadbalance,
		AccessControl: accessControl,
	}
	return detail, err
}

func (t *ServiceInfo) ServiceList(c *gin.Context, tx *gorm.DB, params *dto.ServiceListInput) ([]ServiceInfo, int64, error) {
	var list []ServiceInfo
	var count int64
	pageNo := params.PageNo
	pageSize := params.PageSize

	//limit offset,pagesize
	offset := (pageNo - 1) * pageSize
	query := tx.SetCtx(public.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Select("*")
	query = query.Where("is_delete=?", 0)
	if params.Info != "" {
		query = query.Where(" (service_name like ? or service_desc like ?)", "%"+params.Info+"%", "%"+params.Info+"%")
	}
	err := query.Limit(pageSize).Offset(offset).Order("id desc").Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	errCount := query.Count(&count).Error
	if errCount != nil {
		return nil, 0, err
	}
	return list, count, nil
}
