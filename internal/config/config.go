package config

import (
	"os"
	"encoding/json"
)

const configFileName = ".gatorconfig.json"

type Config struct{
	DB_url string `json:"db_url"`
	Current_username string `json:"current_user_name"`
}


func Read() (Config, error) {
	homePath, err := os.UserHomeDir()
	if err != nil{
		return Config{}, err
	}

	configPath := homePath + "/" +  configFileName

	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	var res Config
	if err = json.Unmarshal(data, &res); err != nil {
		return Config{}, err
	}

	return res, nil
	
}

func (c *Config) SetUser(username string) error {
	c.Current_username = username

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile("/home/andres/.gatorconfig.json", data, 0660)
	if err != nil {
		return err
	}

	return nil
}
