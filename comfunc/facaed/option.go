package facaed

import "github.com/typ7733965/tool/config"

type ConfigOption func(f *Facade)

func ConfigOptionWithApp(app *config.App) ConfigOption {
	return func(f *Facade) {
		if app != nil {
			f.configs.App = app
		}
	}
}
func ConfigOptionWithMysql(mysql *config.MysqlConfig) ConfigOption {
	return func(f *Facade) {
		if mysql != nil {
			f.configs.mysql = mysql
		}
	}
}
func ConfigOptionWithRedis(redis *config.RedisConfig) ConfigOption {
	return func(f *Facade) {
		if redis != nil {
			f.configs.redis = redis
		}
	}
}
