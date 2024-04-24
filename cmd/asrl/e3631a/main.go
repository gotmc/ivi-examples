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
	"time"

	"github.com/gotmc/asrl"
	"github.com/gotmc/ivi/dcpwr/keysight/e36xx"
)

var (
	serialPort string
	baudRate   int
)

func init() {
	// Get serial port used to talk with Keysight E3631A.
	flag.StringVar(
		&serialPort,
		"port",
		"/dev/tty.usbserial-AH03IINA",
		"Serial port for Keysight E3631A",
	)
	flag.IntVar(
		&baudRate,
		"baud",
		9600,
		"Serial port baud rate for Keysight E3631A",
	)
}

func main() {
	// Parse the flags
	flag.Parse()

	// Open the serial port.
	address := fmt.Sprintf("ASRL::%s::%d::8N2::INSTR", serialPort, baudRate)
	log.Printf("VISA Address = %s", address)
	dev, err := asrl.NewDevice(address)
	if err != nil {
		log.Fatal(err)
	}

	dev.HWHandshaking = true

	// Create a new IVI instance of the HP/Agilent/Keysight E3631A DC power
	// supply. Reset the E3631A in order to clear any previous errors.
	ps, err := e36xx.New(dev, true)
	if err != nil {
		log.Fatalf("IVI instrument error: %s", err)
	}
	log.Print("Created new IVI e36xx instrument")

	// Clear and reset the device.
	if err = ps.Clear(); err != nil {
		log.Fatalf("error clearing device: %v", err)
	}
	if err = ps.Reset(); err != nil {
		log.Fatalf("error resetting device: %v", err)
	}

	time.Sleep(500 * time.Millisecond)
	if err = dev.Command("syst:rem"); err != nil {
		log.Fatalf("error setting to remote: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

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

	// Set the output voltage
	desiredVoltage := 5.0
	log.Printf("Set the voltage to %.2f Vdc", desiredVoltage)
	err = ch6v.SetVoltageLevel(desiredVoltage)
	if err != nil && err != io.EOF {
		log.Print(err)
	}

	// Set the current limit
	desiredCurrent := 1.0
	log.Printf("Set the current limit to %.2f Adc", desiredCurrent)
	err = ch6v.SetCurrentLimit(desiredCurrent)
	if err != nil {
		log.Print(err)
	}

	// Enable the 6V output
	log.Printf("Enable 6V output")
	err = ch6v.EnableOutput()
	if err != nil {
		log.Print(err)
	}

	// Query the output voltage setting.
	log.Printf("Query the set output voltage level")
	v, err := ch6v.VoltageLevel()
	if err != nil && err != io.EOF {
		log.Printf("error reading voltage level: %s", err)
	}
	log.Printf("Output voltage on 6V channel = %.3f Vdc", v)

	// Query the current limit.
	log.Printf("Query current limit")
	curr, err := ch6v.CurrentLimit()
	if err != nil {
		log.Print(err)
	}
	log.Printf("Current limit on 6V channel = %.3f Adc", curr)

	// Measure the output voltage.
	log.Println("Measure the output voltage")
	vMsr, err := ch6v.MeasureVoltage()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Measured voltage = %.3f Vdc", vMsr)

	// Measure the output current.
	log.Println("Measure the output current")
	cMsr, err := ch6v.MeasureCurrent()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Measured current = %.3f Adc", cMsr)

	if _, err = dev.Write([]byte("system:local\n")); err != nil {
		log.Fatalf("error setting to local: %v", err)
	}
	dev.Close()

}
