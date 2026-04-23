/*
 * Copyright (c) 2023-2026 TQ-Systems GmbH <license@tq-group.com>, D-82229
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

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/tq-systems/public-go-utils/v3/log"
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
