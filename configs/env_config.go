package configs

import "github.com/spf13/viper"

type Config struct{
	ALCHEMY_API_KEY		string `mapstructure:"ALCHEMY_API_KEY"`
	ETHER_SCAN_API 		string `mapstructure:"ETHER_SCAN_API"`
	SERVER				string `mapstructure:"SERVER"`
	CONTRACT_PATH 		string `mapstructure:"CONTRACT_PATH"`
	MONGO_URI 			string `mapstructure:"MONGO_URI"`
	DATABASE_NAME		string `mapstructure:"DATABASE_NAME"`
}

func LoadConfig(path string) (config Config, err error){
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil{
		return
	}

	err = viper.Unmarshal(&config)
	return
}
