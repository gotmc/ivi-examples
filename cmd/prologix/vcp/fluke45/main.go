// Copyright (c) 2017-2024 The ivi-examples developers. All rights reserved.
// Project site: https://github.com/gotmc/ivi-examples
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"log"

	"github.com/gotmc/ivi/dmm/fluke/fluke45"
	"github.com/gotmc/prologix"
	"github.com/gotmc/prologix/driver/vcp"
)

func main() {
	serialPort := "/dev/tty.usbserial-PX8X3YR6"
	vcp, err := vcp.NewVCP(serialPort)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new GPIB controller using the aforementioned serial port and
	// communicating with the instrument at GPIB address 5.
	gpib, err := prologix.NewController(vcp, 10, true)
	if err != nil {
		log.Fatalf("NewController error: %s", err)
	}
	prologixVer, err := gpib.Version()
	if err != nil {
		log.Fatalf("Unable to determine Prologix controller version: %s", err)
	}
	log.Printf("Using %s", prologixVer)

	// Create a new IVI instance of the HP/Agilent/Keysight E3631A DC power
	// supply.
	dmm, err := fluke45.New(gpib, true)
	if err != nil {
		log.Fatalf("IVI instrument error: %s", err)
	}

	fcn, err := dmm.MeasurementFunction()
	if err != nil {
		log.Fatalf("error getting measurement function: %s", err)
	}
	log.Printf("MeasurementFunction = %s", fcn)

	// Return local control to the front panel.
	err = gpib.FrontPanel(true)
	if err != nil {
		log.Fatalf("error setting local control for front panel: %s", err)
	}

	// Discard any unread data on the serial port and then close.
	err = vcp.Flush()
	if err != nil {
		log.Printf("error flushing serial port: %s", err)
	}
	err = vcp.Close()
	if err != nil {
		log.Printf("error closing serial port: %s", err)
	}
}
