package utils

import (
	"github.com/spf13/viper"
	"github.com/hashicorp/go-hclog"
	"github.com/lib/pq"
)

type Configurations struct {
	ServerPort	string
	DBHost 	string
	DBName  string
	DBUser  string
	DBPass  string
	DBPort  string
	DBConn  string	
	AccessTokenSecrete []byte
	RefreshTokenSecrete []byte
	JwtExpiration	int
	CustomKeySecrete []byte
}

func NewConfigurations(logger hclog.Logger) *Configurations {

	viper.SetEnvPrefix("USER_AUTH")
	viper.AutomaticEnv()

	dbUrl := viper.GetString("DATABASE_URL")
	conn, _ := pq.ParseURL(dbUrl)
	logger.Debug("found database url in env, connection string is formed by parsing it")
	logger.Debug("db connection string", conn)

	viper.SetDefault("SERVER_PORT", "0.0.0.0:9090")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_NAME", "bookite")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "password")
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
		DBConn 	   : conn,
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
	logger.Debug("jwt access token secrete", configs.AccessTokenSecrete)

	return configs
}