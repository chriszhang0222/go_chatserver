package global
import (
	"go_chatserver/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)
var (
	Config *config.Config = &config.Config{}
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
