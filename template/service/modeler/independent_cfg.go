package modeler

// IndependentCfg 个性配置
type IndependentCfg struct {
	Name string `mapstructure:"name"`
}

// Init 用于初始化实例
func (i *IndependentCfg) Init() error {
	// 业务代码

	return nil
}
