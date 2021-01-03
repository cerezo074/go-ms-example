package config

import "github.com/spf13/viper"

//Credentials stores all configuration of the application.
//The values are read by viper from a config file or environment variables
type Credentials struct {
	DBDriver           string `mapstructure:"DB_DRIVER"`
	DBSource           string `mapstructure:"DB_SOURCE"`
	ServerAddress      string `mapstructure:"SERVER_ADDRESS"`
	AWSAccessKeyID     string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AWSSecretKey       string `mapstructure:"AWS_SECRET_KEY"`
	AWSS3ProfileRegion string `mapstructure:"AWS_S3_PROFILE_REGION"`
	AWSS3ProfileBucket string `mapstructure:"AWS_S3_PROFILE_BUCKET"`
}

func LoadCredentials(path string) (config Credentials, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
