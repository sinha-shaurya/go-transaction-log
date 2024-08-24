package database

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"go.uber.org/zap"

	"os"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func getConfigFilePath() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current working directory:", err)
	}
	configPath := filepath.Join(cwd, "/config.yaml")
	return configPath
}

func loadConfig(filepath string) (*Config, error) {
	// read yaml from config
	data, readErr := os.ReadFile(filepath)
	if readErr != nil {
		return nil, fmt.Errorf("%d not found on system", readErr)
	}

	var config MySQLConfig
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal config", err)
	}

	return &config.MysqlConfigObject, nil
}

func OpenConnection(config *Config) (*gorm.DB, error) {
	// Construct the MySQL DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&multiStatements=true",
		config.Username,
		config.Password,
		config.Host,
		config.DBName,
	)

	logConfig := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second * 1.0,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	// Open a GORM connection
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{
		Logger: logConfig,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Enable debug mode if set in config
	if config.Debug {
		db = db.Debug()
	}

	return db, nil
}

func GetDbContext() (*gorm.DB, error) {

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	config, loadErr := loadConfig(getConfigFilePath())
	if loadErr != nil {
		return nil, loadErr
	}

	logger.Info("Loaded config successfully")

	dbContext, err := OpenConnection(config)
	if err != nil {
		return nil, err
	}

	logger.Info("DB initialized successfully")

	sqlDB, err := dbContext.DB()
	if err != nil {
		logger.Sugar().Errorf("Error in getting mysql db %v", err)
		return nil, err
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return dbContext, nil
}
