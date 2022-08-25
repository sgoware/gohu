package config

import (
	"fmt"
	"github.com/go-redis/redis/v8"
)

// GetMysqlDsn 返回 mysql DSN
func (c *Agollo) GetMysqlDsn(namespace string) string {
	config := c.client.GetConfig(namespace)
	serverNum := config.GetIntValue("database.mysql.serverNum", 0) // 使用 mysql 几号服务器(考虑以后将数据库分布式化,使用 TiDB )

	databaseConfig := c.client.GetConfig(databaseNamespace)

	// 拼接dsn字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Asia%%2FShanghai",
		databaseConfig.GetValue(fmt.Sprintf("mysql.server%d.username", serverNum)), // 数据库用户名
		databaseConfig.GetValue(fmt.Sprintf("mysql.server%d.password", serverNum)), // 数据库密码
		databaseConfig.GetValue(fmt.Sprintf("mysql.server%d.address", serverNum)),  // 数据库地址
		databaseConfig.GetValue(fmt.Sprintf("mysql.server%d.port", serverNum)),     // 数据库端口
		config.GetValue("database.mysql.databaseName"),                             // mysql 的数据库名字
		config.GetValue("database.mysql.databaseCharset"),                          // mysql 的数据库使用的字符集
	)
	return dsn
}

// NewRedisOptions 返回 *redisOptions
func (c *Agollo) NewRedisOptions(namespace string) *redis.Options {
	config := c.client.GetConfig(namespace)
	serverNum := config.GetIntValue("database.redis.serverNum", 0) // 使用 redis 几号服务器

	databaseConfig := c.client.GetConfig(databaseNamespace)

	return &redis.Options{
		Addr: fmt.Sprintf("%s:%s",
			databaseConfig.GetValue(fmt.Sprintf("redis.server%d.address", serverNum)),
			databaseConfig.GetValue(fmt.Sprintf("redis.server%d.port", serverNum)),
		),
		Password: databaseConfig.GetValue(fmt.Sprintf("redis.server%d.password", serverNum)),
		DB:       config.GetIntValue("database.redis.databaseNum", 0), // 使用 redis 几号数据库
	}
}
