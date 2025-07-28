package util

import "github.com/spf13/viper"

type Config struct {
	Port      string    `mapstructure:"port"`      //Leader will run on this port
	Log       Log       `mapstructure:"log"`       //Log config
	Heartbeat Heartbeat `mapstructure:"heartbeat"` //Heartbeat config
}

type Log struct {
	Level    int    `mapstructure:"level"`     //Log level
	IsSave   bool   `mapstructure:"is_save"`   //Whether to save logs
	SavePath string `mapstructure:"save_path"` //The path where the logs are saved
}

type Heartbeat struct {
	Overtime int `mapstructure:"overtime"` //Heartbeat overtime in seconds
}

func ReadConfig(confPath, configName string) Config {
	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath(confPath)

	if err := viper.ReadInConfig(); err != nil {
		InitFatalErrorHandel(err)
		return Config{}
	}

	//unmarshal config
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		InitFatalErrorHandel(err)
		return Config{}
	}

	return config
}
