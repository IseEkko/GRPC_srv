package initalize

import (
	"fmt"
	"github.com/spf13/viper"
	"mx_shop/user_srv/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}
func InitConfig() {
	//从配置文件中读取对应的配置
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("%s-pro.yaml", configFilePrefix)
	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}
}
