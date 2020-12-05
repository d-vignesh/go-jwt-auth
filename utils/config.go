package utils

import (
	"github.com/hashicorp/go-hclog"
	"github.com/lib/pq"
	"github.com/spf13/viper"
)

// Configurations wraps all the config variables required by the auth service
type Configurations struct {
	ServerAddress string
	DBHost        string
	DBName        string
	DBUser        string
	DBPass        string
	DBPort        string
	DBConn        string
	// AccessTokenSecrete  []byte
	// RefreshTokenSecrete []byte
	AccessTokenPrivateKeyPath  string
	AccessTokenPublicKeyPath   string
	RefreshTokenPrivateKeyPath string
	RefreshTokenPublicKeyPath  string
	JwtExpiration              int
	// CustomKeySecrete    []byte
}

// NewConfigurations returns a new Configuration object
func NewConfigurations(logger hclog.Logger) *Configurations {

	viper.AutomaticEnv()

	dbURL := viper.GetString("DATABASE_URL")
	conn, _ := pq.ParseURL(dbURL)
	logger.Debug("found database url in env, connection string is formed by parsing it")
	logger.Debug("db connection string", conn)

	viper.SetDefault("SERVER_ADDRESS", "0.0.0.0:9090")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_NAME", "bookite")
	viper.SetDefault("DB_USER", "vignesh")
	viper.SetDefault("DB_PASSWORD", "Vickee@14")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("ACCESS_TOKEN_PRIVATE_KEY_PATH", "./access-private.pem")
	viper.SetDefault("ACCESS_TOKEN_PUBLIC_KEY_PATH", "./access-public.pem")
	viper.SetDefault("REFRESH_TOKEN_PRIVATE_KEY_PATH", "./refresh-private.pem")
	viper.SetDefault("REFRESH_TOKEN_PUBLIC_KEY_PATH", "./refresh-public.pem")
	// viper.SetDefault("REFRESH_JWT_SECRETE_KEY", "superSecretKeyForRefreshToken")
	viper.SetDefault("JWT_EXPIRATION", 30)
	// viper.SetDefault("CUSTOM_SECRETE_KEY", "superSecretKeyForCustomKey")

	configs := &Configurations{
		ServerAddress: viper.GetString("SERVER_ADDRESS"),
		DBHost:        viper.GetString("DB_HOST"),
		DBName:        viper.GetString("DB_NAME"),
		DBUser:        viper.GetString("DB_USER"),
		DBPass:        viper.GetString("DB_PASSWORD"),
		DBPort:        viper.GetString("DB_PORT"),
		DBConn:        conn,
		// AccessTokenSecrete:  []byte(viper.GetString("ACCESS_JWT_SECRETE_KEY")),
		// RefreshTokenSecrete: []byte(viper.GetString("REFRESH_JWT_SECRETE_KEY")),
		JwtExpiration: viper.GetInt("JWT_EXPIRATION"),
		// CustomKeySecrete:    []byte(viper.GetString("CUSTOM_SECRETE_KEY")),
		AccessTokenPrivateKeyPath:  viper.GetString("ACCESS_TOKEN_PRIVATE_KEY_PATH"),
		AccessTokenPublicKeyPath:   viper.GetString("ACCESS_TOKEN_PUBLIC_KEY_PATH"),
		RefreshTokenPrivateKeyPath: viper.GetString("REFRESH_TOKEN_PRIVATE_KEY_PATH"),
		RefreshTokenPublicKeyPath:  viper.GetString("REFRESH_TOKEN_PUBLIC_KEY_PATH"),
	}

	// reading heroku provided port to handle deployment with heroku
	port := viper.GetString("PORT")
	if port != nil {
		logger.Debug("using the port allocated by heroku", port)
		configs.ServerAddress = "0.0.0.0:" + port
	}

	logger.Debug("serve port", configs.ServerAddress)
	logger.Debug("db host", configs.DBHost)
	logger.Debug("db name", configs.DBName)
	logger.Debug("db port", configs.DBPort)
	logger.Debug("jwt expiration", configs.JwtExpiration)

	return configs
}
