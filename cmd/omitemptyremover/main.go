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
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var dir = flag.String("d", "./", "directory")

func main() {
	flag.Parse()

	files, err := ioutil.ReadDir(*dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		removeOmitEmpty(f)
	}

}

func removeOmitEmpty(f os.FileInfo) {
	if !f.IsDir() && strings.HasSuffix(f.Name(), ".pb.go") {
		fileName := *dir + string(os.PathSeparator) + f.Name()
		input, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Panic(err)
		}

		allLines := strings.Split(string(input), "\n")
		log.Println("Removing omitempty in: ", fileName)
		replacer := strings.NewReplacer(",omitempty", "")
		for i, line := range allLines {
			allLines[i] = replacer.Replace(line)
		}
		output := strings.Join(allLines, "\n")
		//READ, WRITE
		err = ioutil.WriteFile(fileName, []byte(output), 0666)
		if err != nil {
			log.Panic(err)
		}
	}
}
