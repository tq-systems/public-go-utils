/*
 * cmd/omitemptyremover - main.go
 * Copyright (c) 2018 - 2023 TQ-Systems GmbH <license@tq-group.com>, D-82229 Seefeld, Germany. All rights reserved.
 * Author: Marcel Matzat and the Energy Manager development team
 *
 * This software code contained herein is licensed under the terms and conditions of
 * the TQ-Systems Product Software License Agreement Version 1.0.1 or any later version.
 * You will find the corresponding license text in the LICENSE file.
 * In case of any license issues please contact license@tq-group.com.
 */
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/tq-systems/public-go-utils/log"
)

var dir = flag.String("d", "./", "directory")

// run ensures that deferred calls are also executed in case of an error
func run() error {
	flag.Parse()

	files, err := os.ReadDir(*dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	for _, f := range files {
		info, err := f.Info()
		if err != nil {
			return fmt.Errorf("failed to get file info: %v", err)
		}
		if err = removeOmitEmpty(info); err != nil {
			return fmt.Errorf("failed removeOmitEmpty: %v", err)
		}
	}
	return nil

}

func removeOmitEmpty(f os.FileInfo) error {
	if !f.IsDir() && strings.HasSuffix(f.Name(), ".pb.go") {
		fileName := *dir + string(os.PathSeparator) + f.Name()
		input, err := os.ReadFile(fileName)
		if err != nil {
			return fmt.Errorf("failed reading file: %v", err)
		}

		allLines := strings.Split(string(input), "\n")
		log.Info("Removing omitempty in: ", fileName)
		replacer := strings.NewReplacer(",omitempty", "")
		for i, line := range allLines {
			allLines[i] = replacer.Replace(line)
		}
		output := strings.Join(allLines, "\n")
		//READ, WRITE
		err = os.WriteFile(fileName, []byte(output), 0666)
		if err != nil {
			return fmt.Errorf("failed writing file: %v", err)
		}
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Errorf("Failed to init app: %v", err)
		os.Exit(1)
	}

}
