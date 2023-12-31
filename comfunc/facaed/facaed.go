package facaed

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/typ7733965/tool/comfunc/env"
	"github.com/typ7733965/tool/comfunc/tool/cache"
	"github.com/typ7733965/tool/comfunc/tool/db"
	"github.com/typ7733965/tool/comfunc/tool/http/fiber"
	"github.com/typ7733965/tool/config"
	"gorm.io/gorm"
)

type ShutdownFunc func(context.Context)

type Config struct {
	App   *config.App         `yaml:"app" json:"app"`
	redis *config.RedisConfig `yaml:"redis" json:"redis"`
	mysql *config.MysqlConfig `yaml:"mysql" json:"mysql"`
	fiber *config.FiberConfig `yaml:"http" json:"http"`
}
type Facade struct {
	configs Config
	env     *env.AppInfo

	redisClient redis.UniversalClient
	fiber       *fiber.App
	mysqlClient *gorm.DB
	shutdown    []ShutdownFunc
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
	if f.configs.redis != nil {
		f.redisClient = cache.InitRedis(f.configs.redis)
	}
	if f.configs.fiber != nil {
		f.configs.fiber.Name = f.env.Name
		if f.fiber, err = fiber.NewApp(f.configs.fiber); err != nil {
			return f, err
		}
	}
	if f.configs.mysql != nil {
		db, err := db.NewMysql(f.configs.mysql)
		if err != nil {
			return f, err
		}
		f.mysqlClient = db
	}
	return
}

func (a *Facade) GetEnv() *env.AppInfo {
	return a.env
}
func (a *Facade) GetRedisClient() redis.UniversalClient {
	return a.redisClient
}
func (a *Facade) GetFiber() *fiber.App {
	return a.fiber
}
func (a *Facade) GetMysqlClient() *gorm.DB {
	return a.mysqlClient
}

func (a *Facade) StartFiber() error {
	if a.fiber != nil {
		a.shutdown = append(a.shutdown, func(ctx context.Context) {
			a.fiber.Fiber().Shutdown()
		})
		return a.fiber.Run()
	}
	return fmt.Errorf("fiber is nil, please init fiber first")
}

// 关闭所有启动的服务
func (a *Facade) Shutdown(ctx context.Context) {
	for _, shutdown := range a.shutdown {
		shutdown(ctx)
	}
}
