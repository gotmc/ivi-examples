// Copyright (c) 2017-2024 The ivi-examples developers. All rights reserved.
// Project site: https://github.com/gotmc/ivi-examples
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gotmc/asrl"
	"github.com/gotmc/ivi"
	"github.com/gotmc/ivi/fgen"
	"github.com/gotmc/ivi/fgen/srs/ds345"
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
	defer dev.Close()

	// Create a new IVI instance of and reset the SRS DS345 function
	// generator using the serial port.
	fg, err := ds345.New(dev, true)
	if err != nil {
		log.Fatalf("IVI instrument error: %s", err)
	}

	// From here forward, we can use the IVI API for the function generator
	// instead of having to send SCPI or other commands that are specific to this
	// model function generator.

	// Query the instrument manufacturer.
	mfr, err := fg.InstrumentManufacturer()
	if err != nil {
		log.Printf("error querying instrument manufacturer: %s", err)
	}
	log.Printf("Instrument manufacturer = %s", mfr)

	// Query the instrument model.
	model, err := fg.InstrumentModel()
	if err != nil {
		log.Printf("error querying instrument model: %s", err)
	}
	log.Printf("Instrument model = %s", model)

	// Query the instrument's serial number.
	sn, err := fg.InstrumentSerialNumber()
	if err != nil {
		log.Printf("error querying instrument sn: %s", err)
	}
	log.Printf("Instrument S/N = %s", sn)

	// Query the firmware revision.
	fw, err := fg.FirmwareRevision()
	if err != nil {
		log.Printf("error querying firmware revision: %s", err)
	}
	log.Printf("Firmware revision = %s", fw)

	// Channel specific methods can be accessed directly from the instrument
	// using 0-based index to select the desired channel.
	if err = fg.Channels[0].DisableOutput(); err != nil && err != ivi.ErrFunctionNotSupported {
		log.Fatalf("error disabling output on ch0: %s", err)
	}
	if err = fg.Channels[0].SetAmplitude(0.5); err != nil {
		log.Fatalf("error setting the amplitude on ch0: %s", err)
	}

	// Alternatively, the channel can be assigned to a variable.
	ch := fg.Channels[0]
	if err = ch.SetStandardWaveform(fgen.Sine); err != nil {
		log.Fatalf("error setting the standard waveform: %s", err)
	}
	if err = ch.SetDCOffset(0.2); err != nil {
		log.Fatalf("error setting DC offest: %s", err)
	}
	if err = ch.SetFrequency(2350); err != nil {
		log.Fatalf("error setting frequency: %s", err)
	}

	// Instead of configuring attributes of a standard waveform individually, the
	// standard waveform can be configured using a single method.
	if err = ch.ConfigureStandardWaveform(fgen.Sine, 0.5, 0.1, 100, 0); err != nil {
		log.Fatalf("error configuring standard waveform: %s", err)
	}
	if err = ch.EnableOutput(); err != nil && err != ivi.ErrFunctionNotSupported {
		log.Fatalf("error enabling output: %s", err)
	}

	if err = ch.SetBurstCount(40); err != nil {
		log.Fatalf("error setting burst count: %s", err)
	}
	if err = ch.SetTriggerSource(fgen.OldTriggerSourceInternal); err != nil {
		log.Fatalf("error setting internal trigger source: %s", err)
	}
	if err = fg.SetInternalTriggerRate(1 / 0.6); err != nil {
		log.Fatalf("error setting internal trigger rate: %s", err)
	}
	if err = ch.SetOperationMode(fgen.BurstMode); err != nil {
		log.Fatalf("error setting operation mode to burst: %s", err)
	}

	// Query the waveform.
	wave, err := ch.StandardWaveform()
	if err != nil {
		log.Printf("error querying standard waveform: %s", err)
	}
	log.Printf("Standard waveform = %s", wave)

	// Query the frequency.
	freq, err := ch.Frequency()
	if err != nil {
		log.Printf("error querying frequency: %s", err)
	}
	log.Printf("Frequency = %.0f Hz", freq)

	// Query the amplitude.
	amp, err := ch.Amplitude()
	if err != nil {
		log.Printf("error querying amplitude: %s", err)
	}
	log.Printf("Amplitude = %.3f Vpp", amp)

	// Query the DC offset.
	offset, err := ch.DCOffset()
	if err != nil {
		log.Printf("error querying DC offset: %s", err)
	}
	log.Printf("DC Offset = %.1f mV", 1000*offset)

	// Query the burst count.
	bc, err := ch.BurstCount()
	if err != nil {
		log.Printf("error querying burst count: %s", err)
	}
	log.Printf("Burst count = %d", bc)

	// Query the internal trigger rate.
	itr, err := fg.InternalTriggerRate()
	if err != nil {
		log.Printf("error querying internal trigger rate: %s", err)
	}
	log.Printf("Internal trigger rate = %.3g Hz", itr)

	// Query the trigger source.
	ts, err := ch.TriggerSource()
	if err != nil {
		log.Printf("error querying trigger source: %s", err)
	}
	log.Printf("Trigger source = %s", ts)

	// Query the operation mode.
	om, err := ch.OperationMode()
	if err != nil {
		log.Printf("error querying operation mode: %s", err)
	}
	log.Printf("Operation mode = %s", om)
}
