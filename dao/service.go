package dao

type ServiceDetail struct {
	Info          *ServiceInfo   `json:"info" validate:"required" toml:"info"`
	HttpRule      *HttpRule      `json:"http_rule" validate:""  toml:"http_rule"`
	TcpRule       *TcpRule       `json:"tcp_rule" validate:""  toml:"tcp_rule"`
	GrpcRule      *GrpcRule      `json:"grpc_rule" validate:""  toml:"grpc_rule"`
	LoadBalance   *LoadBalance   `json:"load_balance" validate:"" toml:"load_balance"`
	AccessControl *AccessControl `json:"access_control" toml:"access_control"`
}
