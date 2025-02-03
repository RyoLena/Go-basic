package config

var Config = config{
	DB: DBConfig{
		//本地
		DSN: "root:123456@tcp(127.0.0.1:13306)/webook",
	},
	Redis: RedisConfig{
		Addr:     "localhost:6379",
		Password: "Ryo19120705",
	},
}
