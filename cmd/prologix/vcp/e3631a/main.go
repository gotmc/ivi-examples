// Copyright (c) 2017-2023 The ivi-examples developers. All rights reserved.
// Project site: https://github.com/gotmc/ivi-examples
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"io"
	"log"

	"github.com/gotmc/ivi/dcpwr/keysight/e36xx"
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
	gpib, err := prologix.NewController(vcp, 5, true)
	if err != nil {
		log.Fatalf("NewController error: %s", err)
	}

	// Query the GPIB instrument address.
	addr, err := gpib.InstrumentAddress()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("GPIB instrument address = %d", addr)

	// Query the Prologix controller version.
	prologixVer, err := gpib.Version()
	if err != nil {
		log.Fatalf("Unable to determine Prologix controller version: %s", err)
	}
	log.Printf("Using %s", prologixVer)

	// Query the auto mode (i.e., read after write).
	auto, err := gpib.ReadAfterWrite()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Read after write = %t", auto)

	// Query the read timeout
	timeout, err := gpib.ReadTimeout()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Read timeout = %d ms", timeout)

	// Determine if the SRQ is asserted.
	srq, err := gpib.ServiceRequest()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Service request asserted = %t", srq)

	// Send the Selected Device Clear (SDC) message
	err = gpib.ClearDevice()
	if err != nil {
		log.Printf("error clearing device: %s", err)
	}

	// Query the identification of the function generator.
	idn, err := gpib.Query("*idn?")
	if err != nil && err != io.EOF {
		log.Fatalf("error querying serial port: %s", err)
	}
	log.Printf("query idn = %s", idn)

	// Create a new IVI instance of the HP/Agilent/Keysight E3631A DC power
	// supply.
	ps, err := e36xx.New(gpib, true)
	if err != nil {
		log.Fatalf("IVI instrument error: %s", err)
	}
	log.Print("Created new IVI e36xx instrument")

	// log.Print("Sending IVI command InstrumentModel")
	// model, err := ps.InstrumentModel()
	// if err != nil {
	// 	log.Fatalf("could not determine instrument model: %s", err)
	// }
	// log.Printf("Instrument model = %s", model)

	// Channel specific methods can be accessed directly from the instrument
	// using 0-based index to select the desired channel.
	log.Printf("Grab first channel, which is the 6V channel.")
	ch6v := ps.Channels[0]
	err = ch6v.DisableOutput()
	if err != nil {
		log.Print(err)
	}

	desiredVoltage := 5.0
	log.Printf("Set the voltage to %.2f V", desiredVoltage)
	err = ch6v.SetVoltageLevel(desiredVoltage)
	if err != nil {
		log.Print(err)
	}

	v, err := ch6v.VoltageLevel()
	if err != nil {
		log.Printf("error reading voltage level: %s", err)
	}
	log.Printf("Voltage = %f", v)

	err = ch6v.SetCurrentLimit(1.0)
	if err != nil {
		log.Print(err)
	}
	err = ch6v.EnableOutput()
	if err != nil {
		log.Print(err)
	}

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
