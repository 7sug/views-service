package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func Read() Settings {
	var result Settings
	filePath := fmt.Sprintf(".config/config.json")
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err = json.Unmarshal(fileBytes, &result); err != nil {
		log.Fatalf("Unmarshalling error: %s", err)
	}

	return result
}
