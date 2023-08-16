/*
 * config package - config_test.go
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
	"fmt"
)

type file struct {
	Config string `json:"config"`
}

func ExampleReadJSON() {
	var readData file
	err := ReadJSON("./config.json", &readData)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(readData.Config)
	// Output: Hello World!
}

func ExampleWriteJSON() {
	writeData := file{Config: "Hello World 2.0!"}

	err := WriteJSON("/tmp/config.json", &writeData)
	if err != nil {
		fmt.Println(err)
		return
	}

	var readData file
	err = ReadJSON("/tmp/config.json", &readData)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(readData.Config)
	// Output: Hello World 2.0!
}
