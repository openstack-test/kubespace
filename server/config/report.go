package config

// 报表配置
type Report struct {
	PutFileUrl []string `mapstructure:"putFileUrl" json:"putFileUrl" yaml:"putFileUrl"`
	CallChat   []string `mapstructure:"callChat" json:"callChat" yaml:"callChat"`
	GateWay    string   `mapstructure:"gateWay" json:"gateWay" yaml:"gateWay"`
	Spec       string   `mapstructure:"spec" json:"spec" yaml:"spec"`
}
