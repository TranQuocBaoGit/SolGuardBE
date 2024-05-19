package configs

import (
	"os"
)

type Config struct{
	ALCHEMY_API_KEY		string `mapstructure:"ALCHEMY_API_KEY"`
	ETHER_SCAN_API 		string `mapstructure:"ETHER_SCAN_API"`
	SERVER				string `mapstructure:"SERVER"`
	CONTRACT_PATH 		string `mapstructure:"CONTRACT_PATH"`
	MONGO_URI 			string `mapstructure:"MONGO_URI"`
	DATABASE_NAME		string `mapstructure:"DATABASE_NAME"`
}

func LoadConfig(path string) (config Config, err error){
	// viper.AddConfigPath(path)
	// viper.SetConfigName(".env")
	// viper.SetConfigType("env")

	// viper.AutomaticEnv()

	// err = viper.ReadInConfig()
	// if err != nil{
	// 	return
	// }

	alchemy := os.Getenv("ALCHEMY_API_KEY")
	etherscan := os.Getenv("ETHER_SCAN_API")
	server := os.Getenv("SERVER")
	mongo := os.Getenv("MONGO_URI")
	database := os.Getenv("DATABASE_NAME")

	// err = viper.Unmarshal(&config)
	return Config{
		ALCHEMY_API_KEY: alchemy,
		ETHER_SCAN_API: etherscan,
		SERVER: server,
		CONTRACT_PATH: "",
		MONGO_URI: mongo,
		DATABASE_NAME: database,
	}, nil
}
