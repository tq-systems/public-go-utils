/*
 * config package - config.go
 * Copyright (c) 2018 - 2023 TQ-Systems GmbH <license@tq-group.com>, D-82229 Seefeld, Germany. All rights reserved.
 * Author: Christoph Krutz and the Energy Manager development team
 *
 * This software code contained herein is licensed under the terms and conditions of
 * the TQ-Systems Product Software License Agreement Version 1.0.1 or any later version.
 * You will find the corresponding license text in the LICENSE file.
 * In case of any license issues please contact license@tq-group.com.
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
