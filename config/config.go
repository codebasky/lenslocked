package config

import (
	"os"
	"strconv"

	"github.com/codebasky/lenslocked/model"
)

type ServerConfig struct {
	Address string
	AuthKey string
}

type Config struct {
	DBCfg     model.PostgresConfig
	SMTPCfg   model.SMTPConfig
	ServerCfg ServerConfig
}

func getServerDefault() ServerConfig {
	return ServerConfig{
		Address: ":3000",
		AuthKey: "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX",
	}
}
func getEnvOrDefault(name string, defValue string) string {
	val, status := os.LookupEnv(name)
	if !status {
		val = defValue
	}
	return val
}

func LoadConfig() Config {
	cfg := Config{}

	dbCfg := model.DefaultPostgresConfig()
	cfg.DBCfg.Host = getEnvOrDefault("DB_HOST", dbCfg.Host)
	cfg.DBCfg.Port = getEnvOrDefault("DB_PORT", dbCfg.Port)
	cfg.DBCfg.User = getEnvOrDefault("DB_USER", dbCfg.User)
	cfg.DBCfg.Password = getEnvOrDefault("DB_PASSWORD", dbCfg.Password)
	cfg.DBCfg.Database = getEnvOrDefault("DB_DATABASE", dbCfg.Database)
	cfg.DBCfg.SSLMode = getEnvOrDefault("DB_SSLMODE", dbCfg.SSLMode)

	smtpCfg := model.DefaultEmailConfig()
	cfg.SMTPCfg.Host = getEnvOrDefault("SMTP_HOST", smtpCfg.Host)
	sport := strconv.Itoa(smtpCfg.Port)
	port, _ := strconv.Atoi(getEnvOrDefault("SMTP_HOST", sport))
	cfg.SMTPCfg.Port = port
	cfg.SMTPCfg.User = getEnvOrDefault("SMTP_USER", smtpCfg.User)
	cfg.SMTPCfg.Password = getEnvOrDefault("SMTP_PASSWORD", smtpCfg.Password)

	srvCfg := getServerDefault()
	cfg.ServerCfg.Address = getEnvOrDefault("WEBSERVER_PORT", srvCfg.Address)
	cfg.ServerCfg.AuthKey = getEnvOrDefault("WEBSERVER_AUTH_KEY", srvCfg.AuthKey)

	return cfg
}
