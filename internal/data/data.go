package data

import (
	"github.com/dtm-labs/rockscache"
	"github.com/fulltimelink/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewDB, NewRedis, NewRockscache)

// Data .
type Data struct {
	Mysql      *gorm.DB
	RedisCli   *redis.Client
	RocksCache *rockscache.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, mysql *gorm.DB, redisCli *redis.Client, rocksCache *rockscache.Client) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		Mysql:      mysql,
		RedisCli:   redisCli,
		RocksCache: rocksCache,
	}, cleanup, nil
}

// NewDB --  @# 初始化数据库链接
func NewDB(c *conf.Data, logger log.Logger) *gorm.DB {
	defer func() {
		if err := recover(); nil != err {
			log.NewHelper(logger).Errorw("kind", "mysql", "Error", err)
		}
	}()
	// --  @# mysql链接
	db, err := gorm.Open(mysql.Open(c.Database.Source))
	if nil != err {
		panic(err)
	}

	sqlDB, err := db.DB()
	if nil != err {
		panic(err)
	}

	// --  @# 链接池参数设置
	sqlDB.SetMaxIdleConns(int(c.Database.MaxIdleConns))
	sqlDB.SetMaxOpenConns(int(c.Database.MaxOpenConns))
	sqlDB.SetConnMaxLifetime(c.Database.ConnMaxLifeTime.AsDuration())
	log.NewHelper(logger).Infow("Kind", "mysql", "status", "enable")
	return db
}

// NewRedis --  @# 初始化Redis
func NewRedis(c *conf.Data, logger log.Logger) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: c.Redis.Addr,
		//Password: c.Redis.Password,
		//DB: int(c.Redis.DB),
		//PoolSize: int(c.Redis.PoolSize),
		//MinIdleConns: int(c.Redis.MinIdleConns),
		//MaxRetries: int(c.Redis.MaxRetries),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
	})
	// --  @# 开启redis tracing
	if err := redisotel.InstrumentTracing(rdb); err != nil {
		panic(err)
	}

	// --  @# 开启redis metrics
	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		panic(err)
	}
	log.NewHelper(logger).Infow("Kind", "redis", "status", "enable")
	return rdb
}

func NewRockscache(c *conf.Data, logger log.Logger, rdb *redis.Client) *rockscache.Client {
	dc := rockscache.NewClient(rdb, rockscache.NewDefaultOptions())
	log.NewHelper(logger).Infow("Kind", "rockscache", "status", "enable")
	return dc
}
