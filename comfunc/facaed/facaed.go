package facaed

import (
	"github.com/go-redis/redis"
	"github.com/typ7733965/tool/comfunc/env"
	"github.com/typ7733965/tool/comfunc/tool/cache"
	"github.com/typ7733965/tool/comfunc/tool/http/fiber"
	"github.com/typ7733965/tool/config"
)

type Config struct {
	App   *config.App         `yaml:"app" json:"app"`
	Redis *config.RedisConfig `yaml:"redis" json:"redis"`
	Mysql *config.MysqlConfig `yaml:"mysql" json:"mysql"`
	fiber *config.FiberConfig `yaml:"http" json:"http"`
}
type Facade struct {
	configs Config
	env     *env.AppInfo

	redisClient redis.UniversalClient
	fiber       *fiber.App
}

func InitApp(opts ...ConfigOption) (f *Facade, err error) {
	if len(opts) == 0 {
		return
	}
	f = &Facade{}
	for _, opt := range opts {
		opt(f)
	}

	if f.configs.App != nil {
		f.env = env.InitEnv(f.configs.App)
	}
	if f.configs.Redis != nil {
		f.redisClient = cache.InitRedis(f.configs.Redis)
	}
	if f.configs.fiber != nil {
		f.configs.fiber.Name = f.env.Name
		if f.fiber, err = fiber.NewApp(f.configs.fiber); err != nil {
			return f, err
		}
	}
	if f.configs.Mysql != nil {
		// todo
	}
	return
}
func (a *Facade) GetEnv() *env.AppInfo {
	return a.env
}
func (a *Facade) GetRedisClient() redis.UniversalClient {
	return a.redisClient
}
