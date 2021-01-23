package global
import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go_chatserver/config"
)
var (
	Config *config.Config = &config.Config{}
	Signal chan bool = make(chan bool, 1)
)

func init(){
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	configFileName := "config.yaml"
	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil{
		zap.S().Error("Error when read config yaml", err.Error())
		panic(err)
	}
	if err := v.Unmarshal(Config);err != nil {
		panic(err)
	}
}
