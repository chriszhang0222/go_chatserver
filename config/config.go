package config

type RedisConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int `mapstructure:"port" json:"port"`
	Db int `mapstructure:db json:"db"`
}

type Config struct{
	Redis RedisConfig `mapstructure:"redis" json:"redis"`
	Port int `mapstructure:"port" json:"port"`
}
