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
