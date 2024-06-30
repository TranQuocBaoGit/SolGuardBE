package configs

import "github.com/joho/godotenv"

type Config struct{
	ALCHEMY_API_KEY		string `mapstructure:"ALCHEMY_API_KEY"`
	ETHER_SCAN_API 		string `mapstructure:"ETHER_SCAN_API"`
	SERVER				string `mapstructure:"SERVER"`
	CONTRACT_PATH 		string `mapstructure:"CONTRACT_PATH"`
	MONGO_URI 			string `mapstructure:"MONGO_URI"`
	DATABASE_NAME		string `mapstructure:"DATABASE_NAME"`
	OPENAI_API_KEY 		string `mapstructure:"OPENAI_API_KEY"`
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
	// err = viper.Unmarshal(&config)

	// alchemy := os.Getenv("ALCHEMY_API_KEY")
	// etherscan := os.Getenv("ETHER_SCAN_API")
	// server := os.Getenv("SERVER")
	// mongo := os.Getenv("MONGO_URI")
	// database := os.Getenv("DATABASE_NAME")

	envFile, _ := godotenv.Read(".env")
	alchemy := envFile["ALCHEMY_API_KEY"]
	etherscan := envFile["ETHER_SCAN_API"]
	server := envFile["SERVER"]
	mongo := envFile["MONGO_URI"]
	database := envFile["DATABASE_NAME"]
	openaiKey := envFile["OPENAI_API_KEY"]


	return Config{
		ALCHEMY_API_KEY: alchemy,
		ETHER_SCAN_API: etherscan,
		SERVER: server,
		CONTRACT_PATH: "",
		MONGO_URI: mongo,
		DATABASE_NAME: database,
		OPENAI_API_KEY: openaiKey,
	}, nil
}
