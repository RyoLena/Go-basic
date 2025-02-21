package config

var Config = config{
	DB: DBConfig{
		//本地
		DSN: "",
	},
	Redis: RedisConfig{
		Addr:     "",
		Password: "",
	},
}
