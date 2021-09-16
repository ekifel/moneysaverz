package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	defaultHTTPPort               = "8000"
	defaultHTTPRWTimeout          = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1
	defaultAccessTokenTTL         = 15 * time.Minute
	defaultRefreshTokenTTL        = 24 * time.Hour * 30

	EnvLocal = "local"
	Prod     = "prod"
)

type (
	Config struct {
		Environment string
		Mongo       MongoConfig
		HTTP        HTTPConfig
		Auth        AuthConfig
	}

	MongoConfig struct {
		URI      string
		User     string
		Password string
		Name     string `mapstructure:"databaseName"`
	}

	AuthConfig struct {
		JWT                    JWTConfig
		PasswordSalt           string
		VerificationCodeLength int `mapstructure:"verificationCodeLength"`
	}

	JWTConfig struct {
		AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
		SigningKey      string
	}

	HTTPConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}
)

func Init(configsDir string) (*Config, error) {
	populateDefaults()

	if err := parseConfigFile(configsDir, os.Getenv("APP_ENV")); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("mongo", &cfg.Mongo); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("auth", &cfg.Auth.JWT); err != nil {
		return err
	}

	return nil
}

func setFromEnv(cfg *Config) {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	cfg.Mongo.URI = os.Getenv("MONGO_URI")
	cfg.Mongo.User = os.Getenv("MONGO_USER")
	cfg.Mongo.Password = os.Getenv("MONGO_PASS")

	cfg.Auth.PasswordSalt = os.Getenv("PASSWORD_SALT")
	cfg.Auth.JWT.SigningKey = os.Getenv("JWT_SIGNING_KEY")

	cfg.HTTP.Host = os.Getenv("HTTP_HOST")

	cfg.Environment = os.Getenv("APP_ENV")
}

func parseConfigFile(folder, env string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if env == EnvLocal {
		return nil
	}

	viper.SetConfigName(env)

	return viper.MergeInConfig()
}

func populateDefaults() {
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("http.max_header_megabytes", defaultHTTPMaxHeaderMegabytes)
	viper.SetDefault("http.timeouts.read", defaultHTTPRWTimeout)
	viper.SetDefault("http.timeouts.write", defaultHTTPRWTimeout)
	viper.SetDefault("auth.accessTokenTTL", defaultAccessTokenTTL)
	viper.SetDefault("auth.refreshTokenTTL", defaultRefreshTokenTTL)
}
