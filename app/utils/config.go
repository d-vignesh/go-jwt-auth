package utils

import (
	"fmt"

	"github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Configurations wraps all the config variables required by the auth service
type Configurations struct {
	ServerAddress              string
	DBHost                     string
	DBName                     string
	DBUser                     string
	DBPass                     string
	DBPort                     string
	DBConn                     string
	AccessTokenPrivateKeyPath  string
	AccessTokenPublicKeyPath   string
	RefreshTokenPrivateKeyPath string
	RefreshTokenPublicKeyPath  string
	JwtExpiration              int // in minutes
	SendGridApiKey             string
	MailVerifCodeExpiration    int // in hours
	PassResetCodeExpiration    int // in minutes
	MailVerifTemplateID        string
	PassResetTemplateID        string
}

// NewConfigurations returns a new Configuration object
func NewConfigurations(logger *zap.Logger) *Configurations {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		logger.Error(fmt.Sprintf("%s %v", "cannot load config:", err))
	}

	dbURL := viper.GetString("DATABASE_URL")
	conn, _ := pq.ParseURL(dbURL)
	logger.Debug("Parsing db URL connection string")
	logger.Debug(fmt.Sprintf("%s %s", "DB connection string: ", conn))

	viper.SetDefault("SERVER_ADDRESS", "0.0.0.0:9000")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_NAME", "test")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("ACCESS_TOKEN_PRIVATE_KEY_PATH", "./access-private.pem")
	viper.SetDefault("ACCESS_TOKEN_PUBLIC_KEY_PATH", "./access-public.pem")
	viper.SetDefault("REFRESH_TOKEN_PRIVATE_KEY_PATH", "./refresh-private.pem")
	viper.SetDefault("REFRESH_TOKEN_PUBLIC_KEY_PATH", "./refresh-public.pem")
	viper.SetDefault("JWT_EXPIRATION", 30)
	viper.SetDefault("MAIL_VERIFICATION_CODE_EXPIRATION", 24)
	viper.SetDefault("PASSWORD_RESET_CODE_EXPIRATION", 15)
	viper.SetDefault("MAIL_VERIFICATION_TEMPLATE_ID", "")
	viper.SetDefault("PASSWORD_RESET_TEMPLATE_ID", "")

	configs := &Configurations{
		ServerAddress:              viper.GetString("SERVER_ADDRESS"),
		DBHost:                     viper.GetString("DB_HOST"),
		DBName:                     viper.GetString("DB_NAME"),
		DBUser:                     viper.GetString("DB_USER"),
		DBPass:                     viper.GetString("DB_PASSWORD"),
		DBPort:                     viper.GetString("DB_PORT"),
		DBConn:                     conn,
		JwtExpiration:              viper.GetInt("JWT_EXPIRATION"),
		AccessTokenPrivateKeyPath:  viper.GetString("ACCESS_TOKEN_PRIVATE_KEY_PATH"),
		AccessTokenPublicKeyPath:   viper.GetString("ACCESS_TOKEN_PUBLIC_KEY_PATH"),
		RefreshTokenPrivateKeyPath: viper.GetString("REFRESH_TOKEN_PRIVATE_KEY_PATH"),
		RefreshTokenPublicKeyPath:  viper.GetString("REFRESH_TOKEN_PUBLIC_KEY_PATH"),
		SendGridApiKey:             viper.GetString("SENDGRID_API_KEY"),
		MailVerifCodeExpiration:    viper.GetInt("MAIL_VERIFICATION_CODE_EXPIRATION"),
		PassResetCodeExpiration:    viper.GetInt("PASSWORD_RESET_CODE_EXPIRATION"),
		MailVerifTemplateID:        viper.GetString("MAIL_VERIFICATION_TEMPLATE_ID"),
		PassResetTemplateID:        viper.GetString("PASSWORD_RESET_TEMPLATE_ID"),
	}

	// Used for cloud deployment
	port := viper.GetString("PORT")
	if port != "" {
		logger.Debug(fmt.Sprintf("%s %s", "using the port allocated by cloud provider", port))
		configs.ServerAddress = "0.0.0.0:" + port
	}

	logger.Debug(fmt.Sprintf("%s %s", "serve port", configs.ServerAddress))
	logger.Debug(fmt.Sprintf("%s %s", "db host", configs.DBHost))
	logger.Debug(fmt.Sprintf("%s %s", "db name", configs.DBName))
	logger.Debug(fmt.Sprintf("%s %s", "db port", configs.DBPort))
	logger.Debug(fmt.Sprintf("%s %d", "jwt expiration", configs.JwtExpiration))

	return configs
}
