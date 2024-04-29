package types // Http
type Http struct {
	Port     int    `mapstructure:"port"`
	Host     string `mapstructure:"host"`
	BasePath string `mapstructure:"base_path"`
}

// Log
type Log struct {
	Level string `mapstructure:"level"`
}

// EnvironmentVariable
type EnvironmentVariable struct {
	Http Http `mapstructure:"http"`
	Log  Log  `mapstructure:"log"`
}

