// Copyright (c) 2017-2024 The ivi-examples developers. All rights reserved.
// Project site: https://github.com/gotmc/ivi-examples
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/gotmc/asrl"
	"github.com/gotmc/ivi/dcpwr/keysight/e36xx"
)

var (
	debugLevel uint
	serialPort string
)

func init() {
	// Get the debug level from CLI flag.
	const (
		defaultLevel = 1
		debugUsage   = "USB debug level"
	)
	flag.UintVar(&debugLevel, "debug", defaultLevel, debugUsage)
	flag.UintVar(&debugLevel, "d", defaultLevel, debugUsage+" (shorthand)")

	// Get serial port used to talk with Keysight E3631A.
	flag.StringVar(
		&serialPort,
		"port",
		"/dev/tty.usbserial-AH03IINA",
		"Serial port for Keysight E3631A",
	)
}

func main() {
	// Parse the flags
	flag.Parse()

	// Open the serial port.
	address := fmt.Sprintf("ASRL::%s::9600::8N2::INSTR", serialPort)
	log.Printf("VISA Address = %s", address)
	dev, err := asrl.NewDevice(address)
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	// Create a new IVI instance of the HP/Agilent/Keysight E3631A DC power
	// supply.
	ps, err := e36xx.New(dev, true)
	if err != nil {
		log.Fatalf("IVI instrument error: %s", err)
	}
	log.Print("Created new IVI e36xx instrument")

	log.Print("Sending IVI command InstrumentModel")
	model, err := ps.InstrumentModel()
	if err != nil {
		log.Fatalf("could not determine instrument model: %s", err)
	}
	log.Printf("Instrument model = %s", model)

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
	if err != nil && err != io.EOF {
		log.Print(err)
	}

	v, err := ch6v.VoltageLevel()
	if err != nil && err != io.EOF {
		log.Printf("error reading voltage level: %s", err)
	}
	log.Printf("Output Voltage on 6V channel = %.3f Vdc", v)

	err = ch6v.SetCurrentLimit(1.0)
	if err != nil {
		log.Print(err)
	}
	err = ch6v.EnableOutput()
	if err != nil {
		log.Print(err)
	}
}
