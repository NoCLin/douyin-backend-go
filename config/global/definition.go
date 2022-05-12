package global

type Configuration struct {
	Server   `yaml:"server"`
	Redis    `yaml:"redis"`
	Database `yaml:"database"`
}

type Server struct {
	Addr      string `yaml:"addr"`
	URLPrefix string `yaml:"prefix"`
	Mode      string `yaml:"mode"`

	//LimitNum  int    `yaml:"limitNum"`
	//UserMongo bool   `yaml:"useMongo"`
	//UserRedis bool   `yaml:"useRedis"`
}

type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Database struct {
	Type     string `yaml:"type"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     int    `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`

	MaxIdleConns int `yaml:"maxIdleConns"`
	MaxOpenConns int `yaml:"maxOpenConns"`
	//Log          bool `yaml:"log"`
	AutoMigrate bool `yaml:"autoMigrate"`
}
