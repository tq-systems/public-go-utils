/*
 * REST authentication utilities - auth.go
 * Copyright (c) 2020 - 2023 TQ-Systems GmbH <license@tq-group.com>, D-82229 Seefeld, Germany. All rights reserved.
 * Author: Matthias Schiffer and the Energy Manager development team
 *
 * This software code contained herein is licensed under the terms and conditions of
 * the TQ-Systems Product Software License Agreement Version 1.0.1 or any later version.
 * You will find the corresponding license text in the LICENSE file.
 * In case of any license issues please contact license@tq-group.com.
 */
package auth

import "github.com/godbus/dbus/v5"

// User type
type User struct {
	Name  string
	Roles []string
}

// HasRole checks if a user has a role
func (u User) HasRole(role interface{}) bool {
	switch role := role.(type) {
	case []string:
		for _, userRole := range u.Roles {
			for _, routeRole := range role {
				if userRole == routeRole {
					return true
				}
			}
		}
		return false
	case string:
		for _, roleInRoles := range u.Roles {
			if roleInRoles == role {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// ValidateAuthToken validates a token
func ValidateAuthToken(token string) (User, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return User{}, err
	}

	validator := conn.Object("com.tq_group.tq_em.web_login1", "/com/tq_group/tq_em/web_login1")

	call := validator.Call("com.tq_group.tq_em.web_login1.ValidateAuthToken", 0, token)
	if call.Err != nil {
		return User{}, call.Err
	}

	var user User
	err = call.Store(&user.Name, &user.Roles)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
