package database

type MySQLConfig struct {
	MysqlConfigObject Config `yaml:"mysql"`
}

type Config struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Debug    bool   `yaml:"debug"`
	DBName   string `yaml:"dbname"`
}
