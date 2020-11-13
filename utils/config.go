package utils

import (
	"github.com/spf13/viper"
	"github.com/hashicorp/go-hclog"
)

type Configurations struct {
	ServerPort	string
	DBHost 	string
	DBName  string
	DBUser  string
	DBPass  string
	DBPort  string	
	AccessTokenSecrete []byte
	RefreshTokenSecrete []byte
	JwtExpiration	int
	CustomKeySecrete []byte
}

func NewConfigurations(logger hclog.Logger) *Configurations {

	viper.SetEnvPrefix("USER_AUTH")
	viper.AutomaticEnv()

	viper.SetDefault("SERVER_PORT", ":9090")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_NAME", "db_name")
	viper.SetDefault("DB_USER", "db_user")
	viper.SetDefault("DB_PASSWORD", "db_password")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("ACCESS_JWT_SECRETE_KEY", "superSecretKeyForAccessToken")
	viper.SetDefault("REFRESH_JWT_SECRETE_KEY", "superSecretKeyForRefreshToken")
	viper.SetDefault("JWT_EXPIRATION", 30)
	viper.SetDefault("CUSTOM_SECRETE_KEY", "superSecretKeyForCustomKey")

	configs := &Configurations {
		ServerPort : viper.GetString("SERVER_PORT"),
		DBHost	   : viper.GetString("DB_HOST"),
		DBName	   : viper.GetString("DB_NAME"),
		DBUser 	   : viper.GetString("DB_USER"),
		DBPass	   : viper.GetString("DB_PASSWORD"),
		DBPort	   : viper.GetString("DB_PORT"),
		AccessTokenSecrete : []byte(viper.GetString("ACCESS_JWT_SECRETE_KEY")),
		RefreshTokenSecrete : []byte(viper.GetString("REFRESH_JWT_SECRETE_KEY")),
		JwtExpiration	: viper.GetInt("JWT_EXPIRATION"),
		CustomKeySecrete : []byte(viper.GetString("CUSTOM_SECRETE_KEY")),
	}
	logger.Debug("serve port", configs.ServerPort)
	logger.Debug("db host", configs.DBHost)
	logger.Debug("db name", configs.DBName)
	logger.Debug("db port", configs.DBPort)
	logger.Debug("jwt expiration", configs.JwtExpiration)

	return configs
}