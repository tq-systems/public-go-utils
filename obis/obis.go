/*
 * OBIS code mapping - obis.go
 * Copyright (c) 2018 - 2023 TQ-Systems GmbH <license@tq-group.com>, D-82229 Seefeld, Germany. All rights reserved.
 * Author: Christoph Krutz and the Energy Manager development team
 *
 * This software code contained herein is licensed under the terms and conditions of
 * the TQ-Systems Product Software License Agreement Version 1.0.1 or any later version.
 * You will find the corresponding license text in the LICENSE file.
 * In case of any license issues please contact license@tq-group.com.
 */
package obis

// ObisLabel maps OBIS code to readable unit
var ObisLabel = map[string]string{
	"1-0:1.4.0*255":  "Active power+",
	"1-0:2.4.0*255":  "Active power-",
	"1-0:3.4.0*255":  "Reactive power+",
	"1-0:4.4.0*255":  "Reactive power-",
	"1-0:9.4.0*255":  "Apparent power+",
	"1-0:10.4.0*255": "Apparent power-",
	"1-0:13.4.0*255": "Power factor",
	"1-0:14.4.0*255": "Supply frequency",
	"1-0:21.4.0*255": "Active power+ (L1)",
	"1-0:22.4.0*255": "Active power- (L1)",
	"1-0:23.4.0*255": "Reactive power+ (L1)",
	"1-0:24.4.0*255": "Reactive power- (L1)",
	"1-0:29.4.0*255": "Apparent power+ (L1)",
	"1-0:30.4.0*255": "Apparent power- (L1)",
	"1-0:31.4.0*255": "Current (L1)",
	"1-0:32.4.0*255": "Voltage (L1)",
	"1-0:33.4.0*255": "Power factor (L1)",
	"1-0:41.4.0*255": "Active power+ (L2)",
	"1-0:42.4.0*255": "Active power- (L2)",
	"1-0:43.4.0*255": "Reactive power+ (L2)",
	"1-0:44.4.0*255": "Reactive power- (L2)",
	"1-0:49.4.0*255": "Apparent power+ (L2)",
	"1-0:50.4.0*255": "Apparent power- (L2)",
	"1-0:51.4.0*255": "Current (L2)",
	"1-0:52.4.0*255": "Voltage (L2)",
	"1-0:53.4.0*255": "Power factor (L2)",
	"1-0:61.4.0*255": "Active power+ (L3)",
	"1-0:62.4.0*255": "Active power- (L3)",
	"1-0:63.4.0*255": "Reactive power+ (L3)",
	"1-0:64.4.0*255": "Reactive power- (L3)",
	"1-0:69.4.0*255": "Apparent power+ (L3)",
	"1-0:70.4.0*255": "Apparent power- (L3)",
	"1-0:71.4.0*255": "Current (L3)",
	"1-0:72.4.0*255": "Voltage (L3)",
	"1-0:73.4.0*255": "Power factor (L3)",
	"1-0:1.8.0*255":  "Active energy+",
	"1-0:2.8.0*255":  "Active energy-",
	"1-0:3.8.0*255":  "Reactive energy+",
	"1-0:4.8.0*255":  "Reactive energy-",
	"1-0:9.8.0*255":  "Apparent energy+",
	"1-0:10.8.0*255": "Apparent energy-",
	"1-0:21.8.0*255": "Active energy+ (L1)",
	"1-0:22.8.0*255": "Active energy- (L1)",
	"1-0:23.8.0*255": "Reactive energy+ (L1)",
	"1-0:24.8.0*255": "Reactive energy- (L1)",
	"1-0:29.8.0*255": "Apparent energy+ (L1)",
	"1-0:30.8.0*255": "Apparent energy- (L1)",
	"1-0:41.8.0*255": "Active energy+ (L2)",
	"1-0:42.8.0*255": "Active energy- (L2)",
	"1-0:43.8.0*255": "Reactive energy+ (L2)",
	"1-0:44.8.0*255": "Reactive energy- (L2)",
	"1-0:49.8.0*255": "Apparent energy+ (L2)",
	"1-0:50.8.0*255": "Apparent energy- (L2)",
	"1-0:61.8.0*255": "Active energy+ (L3)",
	"1-0:62.8.0*255": "Active energy- (L3)",
	"1-0:63.8.0*255": "Reactive energy+ (L3)",
	"1-0:64.8.0*255": "Reactive energy- (L3)",
	"1-0:69.8.0*255": "Apparent energy+ (L3)",
	"1-0:70.8.0*255": "Apparent energy- (L3)",
}
