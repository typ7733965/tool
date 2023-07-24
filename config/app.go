package config

type App struct {
	Name  string `yaml:"name"`  // 应用名称
	Dev   bool   `yaml:"dev"`   // 是否为开发环境
	Stage string `yaml:"stage"` // 环境local,develop,release,production
}
