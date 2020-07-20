package util

type RedisConfig struct {
	Host string
	Port string
	Password string
}

var Redis = new(RedisConfig)

func init() {
	//Redis.Host = "127.0.0.1"
	Redis.Port = "6379"
	Redis.Password = ""

	//Redis.Host = os.Getenv("Redis_Host")
	//Redis.Port = os.Getenv("Redis_Port")

}
