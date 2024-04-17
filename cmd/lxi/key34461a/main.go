// Copyright (c) 2017-2024 The ivi-examples developers. All rights reserved.
// Project site: https://github.com/gotmc/ivi-examples
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gotmc/ivi/dmm"
	"github.com/gotmc/ivi/dmm/keysight/key3446x"
	"github.com/gotmc/lxi"
)

func main() {
	log.Println("IVI LXI Keysight 33461A Example Application")

	// Get IP address from CLI flag.
	var ip string
	flag.StringVar(
		&ip,
		"ip",
		"10.12.100.56",
		"IP address of Keysight 3461A",
	)
	flag.Parse()

	// Create a new LXI device
	address := fmt.Sprintf("TCPIP0::%s::5025::SOCKET", ip)
	log.Printf("VISA address = %s", address)
	dev, err := lxi.NewDevice(address)
	if err != nil {
		log.Fatalf("NewDevice error: %s", err)
	}

	// Close the LXI device when done.
	defer dev.Close()

	// Create a new IVI instance of and reset the Agilent 33220 function
	// generator using the LXI device.
	d, err := key3446x.New(dev, true)
	if err != nil {
		log.Fatalf("IVI instrument eror: %s", err)
	}

	// From here forward, we can use the IVI API for the function generator
	// instead of having to send SCPI or other commands that are specific to this
	// model function generator.

	// Query the instrument manufacturer.
	mfr, err := d.InstrumentManufacturer()
	if err != nil {
		log.Printf("error querying instrument manufacturer: %s", err)
	}
	log.Printf("Instrument manufacturer = %s", mfr)

	// Query the instrument model.
	model, err := d.InstrumentModel()
	if err != nil {
		log.Printf("error querying instrument model: %s", err)
	}
	log.Printf("Instrument model = %s", model)

	// Query the instrument's serial number.
	sn, err := d.InstrumentSerialNumber()
	if err != nil {
		log.Printf("error querying instrument sn: %s", err)
	}
	log.Printf("Instrument S/N = %s", sn)

	// Query the firmware revision.
	fw, err := d.FirmwareRevision()
	if err != nil {
		log.Printf("error querying firmware revision: %s", err)
	}
	log.Printf("Firmware revision = %s", fw)

	// Query the measurement function.
	fcn, err := d.MeasurementFunction()
	if err != nil {
		log.Printf("error querying the measurement function: %s", err)
	}
	log.Printf("Measurement function = %s", fcn)

	// Set the measurement function to DC volts and then query.
	newFcn := dmm.DCVolts
	log.Printf("Setting the measurement function to %s", newFcn)
	err = d.SetMeasurementFunction(newFcn)
	if err != nil {
		log.Printf("error setting the measurement function: %s", err)
	}
	fcn, err = d.MeasurementFunction()
	if err != nil {
		log.Printf("error querying the measurement function: %s", err)
	}
	log.Printf("Measurement function = %s", fcn)

	// Query the range.
	autoRange, rng, err := d.Range()
	if err != nil {
		log.Printf("error querying the range: %s", err)
	}
	log.Printf("Range = %.f V / %s", rng, autoRange)

	// Set the manual range.
	err = d.SetRange(dmm.AutoOff, 10.0)
	if err != nil {
		log.Printf("error setting the manual range to 10V: %s", err)
	}

	// Query the range.
	autoRange, rng, err = d.Range()
	if err != nil {
		log.Printf("error querying the range: %s", err)
	}
	log.Printf("Range = %.f V / %s", rng, autoRange)

	// Set the auto range.
	err = d.SetRange(dmm.AutoOn, 0.0)
	if err != nil {
		log.Printf("error setting the auto range: %s", err)
	}

	// Query the range.
	autoRange, rng, err = d.Range()
	if err != nil {
		log.Printf("error querying the range: %s", err)
	}
	log.Printf("Range = %.f V / %s", rng, autoRange)

	// Set the measurement function to resistance and then query.
	newFcn = dmm.TwoWireResistance
	log.Printf("Setting the measurement function to %s", newFcn)
	err = d.SetMeasurementFunction(newFcn)
	if err != nil {
		log.Printf("error setting the measurement function: %s", err)
	}
	fcn, err = d.MeasurementFunction()
	if err != nil {
		log.Printf("error querying the measurement function: %s", err)
	}
	log.Printf("Measurement function = %s", fcn)

	// Set the manual range.
	err = d.SetRange(dmm.AutoOff, 100e6)
	if err != nil {
		log.Printf("error setting the manual range to 100 MΩ: %s", err)
	}

	// Query the range.
	autoRange, rng, err = d.Range()
	if err != nil {
		log.Printf("error querying the range: %s", err)
	}
	log.Printf("Range = %g ohms / %s", rng, autoRange)

	// Set the measurement function to DC volts and then query.
	newFcn = dmm.DCVolts
	log.Printf("Setting the measurement function to %s", newFcn)
	err = d.SetMeasurementFunction(newFcn)
	if err != nil {
		log.Printf("error setting the measurement function: %s", err)
	}
	fcn, err = d.MeasurementFunction()
	if err != nil {
		log.Printf("error querying the measurement function: %s", err)
	}
	log.Printf("Measurement function = %s", fcn)

	// Set the range to auto.
	log.Println("Enabling auto range")
	err = d.SetRange(dmm.AutoOn, 0.0)
	if err != nil {
		log.Printf("error enabling auto range: %s", err)
	}
	autoRange, rng, err = d.Range()
	if err != nil {
		log.Printf("error querying the range: %s", err)
	}
	log.Printf("Range = %g ohms / %s", rng, autoRange)

	// Read the measurement.
	msr, err := d.ReadMeasurement(100 * time.Millisecond)
	if err != nil {
		log.Printf("error reading the measurement: %s", err)
	}
	log.Printf("Measurement reading #1 = %g V", msr)

	// Read the measurement.
	msr, err = d.ReadMeasurement(100 * time.Millisecond)
	if err != nil {
		log.Printf("error reading the measurement: %s", err)
	}
	log.Printf("Measurement reading #2 = %g V", msr)

	// Read the measurement.
	msr, err = d.ReadMeasurement(100 * time.Millisecond)
	if err != nil {
		log.Printf("error reading the measurement: %s", err)
	}
	log.Printf("Measurement reading #3 = %g V", msr)

	// Set the measurement function to frequency and then query.
	newFcn = dmm.Frequency
	log.Printf("Setting the measurement function to %s", newFcn)
	err = d.SetMeasurementFunction(newFcn)
	if err != nil {
		log.Printf("error setting the measurement function: %s", err)
	}
	fcn, err = d.MeasurementFunction()
	if err != nil {
		log.Printf("error querying the measurement function: %s", err)
	}
	log.Printf("Measurement function = %s", fcn)

	// Read the frequency.
	msr, err = d.ReadMeasurement(100 * time.Millisecond)
	if err != nil {
		log.Printf("error reading the measurement: %s", err)
	}
	log.Printf("Frequency = %g Hz", msr)

	// Set the measurement function to period and then query.
	newFcn = dmm.Period
	log.Printf("Setting the measurement function to %s", newFcn)
	err = d.SetMeasurementFunction(newFcn)
	if err != nil {
		log.Printf("error setting the measurement function: %s", err)
	}
	fcn, err = d.MeasurementFunction()
	if err != nil {
		log.Printf("error querying the measurement function: %s", err)
	}
	log.Printf("Measurement function = %s", fcn)

	// Read the period.
	msr, err = d.ReadMeasurement(100 * time.Millisecond)
	if err != nil {
		log.Printf("error reading the measurement: %s", err)
	}
	log.Printf("Period = %g s", msr)

	// Query the terminals selected.
	term, err := d.SelectedTerminals()
	if err != nil {
		log.Printf("error querying the select terminals: %s", err)
	}
	log.Printf("Select terminals = %s", term)

}
