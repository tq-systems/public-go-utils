/*
 * Copyright (c) 2021-2026 TQ-Systems GmbH <license@tq-group.com>, D-82229
 * Seefeld, Germany. All rights reserved.
 * Author: Christoph Krutz and the Energy Manager development team
 *
 * This software is licensed under the TQ-Systems Product Software License
 * Agreement Version 1.0.3 or any later version.
 * You can obtain a copy of the License Agreement in the TQS (TQ-Systems
 * Software Licenses) folder on the following website:
 * https://www.tq-group.com/en/support/downloads/tq-software-license-conditions/
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
