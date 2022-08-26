package config

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetMysqlDB(namespace string) (db *gorm.DB, err error) {
	if agolloClient == nil {
		return nil, errEmptyConfigClient
	}
	dsn, err := agolloClient.getMysqlDsn(namespace)
	if err != nil {
		return nil, err
	}

	db, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetRedisClient(namespace string) (rdb *redis.Client, err error) {
	if agolloClient == nil {
		return nil, errEmptyConfigClient
	}
	redisOptions, err := agolloClient.newRedisOptions(namespace)
	if err != nil {
		return nil, err
	}

	rdb = redis.NewClient(redisOptions)
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

// getMysqlDsn 返回 mysql DSN
func (c *Agollo) getMysqlDsn(namespace string) (dsn string, err error) {
	v, err := c.GetViper(namespace)
	if err != nil {
		return "", err
	}
	// 使用 mysql 几号服务器(考虑以后将数据库分布式化,使用 TiDB )
	serverNum := v.GetInt("Database.Mysql.ServerNum")

	databaseViper, err := c.GetViper(databaseNamespace)
	if err != nil {
		return "", err
	}

	// 拼接dsn字符串
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Asia%%2FShanghai",
		databaseViper.GetString(fmt.Sprintf("Mysql.Server%d.Username", serverNum)), // 数据库用户名
		databaseViper.GetString(fmt.Sprintf("Mysql.Server%d.Password", serverNum)), // 数据库密码
		databaseViper.GetString(fmt.Sprintf("Mysql.Server%d.Address", serverNum)),  // 数据库地址
		databaseViper.GetString(fmt.Sprintf("Mysql.Server%d.Port", serverNum)),     // 数据库端口
		v.GetString("Database.Mysql.DatabaseName"),                                 // mysql 的数据库名字
		v.GetString("Database.Mysql.DatabaseCharset"),                              // mysql 的数据库使用的字符集
	)
	return dsn, nil
}

// newRedisOptions 返回 *redisOptions
func (c *Agollo) newRedisOptions(namespace string) (*redis.Options, error) {
	v, err := c.GetViper(namespace)
	if err != nil {
		return nil, err
	}
	serverNum := v.GetInt("Database.Redis.ServerNum") // 使用 redis 几号服务器

	databaseViper, err := c.GetViper(databaseNamespace)
	if err != nil {
		return nil, err
	}

	return &redis.Options{
		Addr: fmt.Sprintf("%s:%s",
			databaseViper.GetString(fmt.Sprintf("Redis.Server%d.Address", serverNum)),
			databaseViper.GetString(fmt.Sprintf("Redis.Server%d.Port", serverNum)),
		),
		Password: databaseViper.GetString(fmt.Sprintf("Redis.Server%d.Password", serverNum)),
		DB:       v.GetInt("Database.Redis.DatabaseNum"), // 使用 redis 几号数据库
	}, nil
}
