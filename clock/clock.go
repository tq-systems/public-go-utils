/*
 *  Copyright (c) 2024 TQ-Systems GmbH <license@tq-group.com>, D-82229 Seefeld, Germany. All rights reserved.
 *  Author: Thomas Stöter and the Energy Manager development team
 *
 *  This software code contained herein is licensed under the terms and conditions of the TQ-Systems Software License Agreement Version 1.0.3 or any later version.
 *  You may obtain a copy of the TQ-Systems Software License Agreement in the folder TQS (TQ-Systems Software Licenses) at the following website:
 *  https://www.tq-group.com/en/support/downloads/tq-software-license-conditions/
 *  In case of any license issues please contact license@tq-group.com.
 *  The corresponding license text can also be found in the LICENSE file.
 */

package clock

import "time"

//go:generate mockgen -destination=../mocks/clock/mock_clock.go -build_flags "--mod=mod" -package=clock -source=clock.go Clock

type Clock interface {
	Now() time.Time
	Since(time.Time) time.Duration
}

type SystemCLock struct{}

func (SystemCLock) Now() time.Time { return time.Now() }

func (SystemCLock) Since(then time.Time) time.Duration { return time.Since(then) }
