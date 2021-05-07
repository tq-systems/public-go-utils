/*
 * configuration package - config.go
 * Copyright (c) 2018, TQ-Systems GmbH
 * All rights reserved. For further information see LICENSE.
 */

package config

import (
	"encoding/json"
	"os"
)

// ReadJSON reads and parses a JSON configuration file.
func ReadJSON(file string, dst interface{}) error {
	// Open file for reading
	configFile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer configFile.Close()

	// create new decoder and decode file
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(dst)
	if err != nil {
		return err
	}

	return nil
}

// WriteJSON writes data into the specified file.
// If the file does not exist, the file will be created.
func WriteJSON(file string, dst interface{}) error {
	configFile, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	jsonParser := json.NewEncoder(configFile)
	err = jsonParser.Encode(dst)
	configFile.Close()
	if err != nil {
		return err
	}

	return nil
}
