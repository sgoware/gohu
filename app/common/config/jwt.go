package config

type JWTConfig struct { // jwt鉴权配置
	SecretKey   string // 密钥
	ExpiresTime int64  // 过期时间,单位:秒
	BufferTime  int64  // 缓冲时间
	Issuer      string // 签发者
}
