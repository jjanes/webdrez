package config

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
)

type Config struct {
	Socials []Social `json:"socials"`
	Theme   string   `json:"theme"`
}
type Social struct {
	Name string `json:"name"`
	User string `json:"user"`
	Icon template.HTML
	Url  string `json:"url"`
}

func Load(file string) (Config, error) {
	config := Config{}

	data, err := os.ReadFile(file)

	if err != nil {
		return config, err
	}

	if err = json.Unmarshal(data, &config); err != nil {
		return config, err
	}

	fmt.Printf(">> %s", data)

	return config, nil
}
