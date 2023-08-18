/*
 * REST utilities - listen_linux.go
 * Copyright (c) 2020 - 2023 TQ-Systems GmbH <license@tq-group.com>, D-82229 Seefeld, Germany. All rights reserved.
 * Author: Matthias Schiffer and the Energy Manager development team
 *
 * This software code contained herein is licensed under the terms and conditions of
 * the TQ-Systems Product Software License Agreement Version 1.0.1 or any later version.
 * You will find the corresponding license text in the LICENSE file.
 * In case of any license issues please contact license@tq-group.com.
 */

package rest

import (
	"net"
	"os/user"
	"strconv"

	"golang.org/x/sys/unix"
)

func chgrp(path string, group string) error {
	grp, err := user.LookupGroup(group)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(grp.Gid)
	if err != nil {
		return err
	}
	err = unix.Chown(path, -1, gid)
	if err != nil {
		return err
	}
	err = unix.Chmod(path, 0664)
	if err != nil {
		return err
	}

	return nil
}

// Listen returns a rest listener
func Listen(proto string, listen string, group string) (net.Listener, error) {
	if proto == "unix" {
		// No error handling needed: this may fail when the socket file does not exist;
		// if something goes wrong (permissions etc.), the net.Listen call will return
		// an error
		_ = unix.Unlink(listen)
	}

	listener, err := net.Listen(proto, listen)
	if err != nil {
		return nil, err
	}

	if proto == "unix" && group != "" {
		err = chgrp(listen, group)
		if err != nil {
			listener.Close()
			return nil, err
		}
	}

	return listener, nil
}
