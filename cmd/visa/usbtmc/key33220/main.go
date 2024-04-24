// Copyright (c) 2017-2024 The ivi-examples developers. All rights reserved.
// Project site: https://github.com/gotmc/ivi-examples
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"flag"
	"log"

	"github.com/gotmc/ivi/fgen"
	"github.com/gotmc/ivi/fgen/keysight/key33220"
	_ "github.com/gotmc/usbtmc/driver/google"
	"github.com/gotmc/visa"
	_ "github.com/gotmc/visa/driver/usbtmc"
)

var (
	debugLevel uint
	address    string
)

func init() {
	// Get the debug level from CLI flag.
	const (
		defaultLevel = 1
		debugUsage   = "USB debug level"
	)
	flag.UintVar(&debugLevel, "debug", defaultLevel, debugUsage)
	flag.UintVar(&debugLevel, "d", defaultLevel, debugUsage+" (shorthand)")

	// Get VISA address from CLI flag.
	flag.StringVar(
		&address,
		"visa",
		"USB0::2391::1031::MY44035849::INSTR",
		"VISA address of Keysight 33220A",
	)
}

func main() {
	log.Println("IVI USBTMC Keysight 33220A Example Application")

	// Parse the flags
	flag.Parse()

	// Configure a new VISA resource using the USBTMC driver.
	log.Printf("VISA address = %s", address)
	res, err := visa.NewResource(address)
	if err != nil {
		log.Fatalf("VISA resource %s: %s", address, err)
	}

	// Create a new IVI instance of and reset the Agilent 33220 function
	// generator using the USBTMC device.
	inst, err := key33220.New(res, true)
	if err != nil {
		log.Fatalf("IVI instrument error: %s", err)
	}

	// From here forward, we can use the IVI API for the function generator
	// instead of having to send SCPI or other commands that are specific to this
	// model function generator.

	// Query the instrument manufacturer.
	mfr, err := inst.InstrumentManufacturer()
	if err != nil {
		log.Printf("error querying instrument manufacturer: %s", err)
	}
	log.Printf("Instrument manufacturer = %s", mfr)

	// Query the instrument model.
	model, err := inst.InstrumentModel()
	if err != nil {
		log.Printf("error querying instrument model: %s", err)
	}
	log.Printf("Instrument model = %s", model)

	// Query the instrument's serial number.
	sn, err := inst.InstrumentSerialNumber()
	if err != nil {
		log.Printf("error querying instrument sn: %s", err)
	}
	log.Printf("Instrument S/N = %s", sn)

	// Query the firmware revision.
	fw, err := inst.FirmwareRevision()
	if err != nil {
		log.Printf("error querying firmware revision: %s", err)
	}
	log.Printf("Firmware revision = %s", fw)

	// Channel specific methods can be accessed directly from the instrument
	// using 0-based index to select the desirec channel.
	if err = inst.Channels[0].DisableOutput(); err != nil {
		log.Fatalf("error disabling output on ch0: %s", err)
	}
	if err = inst.Channels[0].SetAmplitude(2.1); err != nil {
		log.Fatalf("error setting the amplitude on ch0: %s", err)
	}

	// Alternatively, the channel can be assigned to a variable.
	ch := inst.Channels[0]
	if err = ch.SetStandardWaveform(fgen.Sine); err != nil {
		log.Fatalf("error setting the standard waveform: %s", err)
	}
	if err = ch.SetDCOffset(0.3); err != nil {
		log.Fatalf("error setting DC offest: %s", err)
	}
	if err = ch.SetFrequency(2230); err != nil {
		log.Fatalf("error setting frequency: %s", err)
	}

	// Instead of configuring attributes of a standard waveform individually, the
	// standard waveform can be configured using a single method.
	if err = ch.ConfigureStandardWaveform(fgen.Sine, 0.5, 0.0, 100, 0); err != nil {
		log.Fatalf("error configuring standard waveform: %s", err)
	}

	// Setup a bursted sinusoidal waveform.
	if err = ch.SetBurstCount(10); err != nil {
		log.Fatalf("error setting burst count: %s", err)
	}
	// Set the code period to 60 ms.
	if err = inst.SetInternalTriggerRate(1 / 0.6); err != nil {
		log.Fatalf("error setting the internal trigger rate: %s", err)
	}
	if err = ch.SetStartTriggerSource(fgen.TriggerSourceInternal); err != nil {
		log.Fatalf("error setting the trigger source: %s", err)
	}
	if err = ch.SetOperationMode(fgen.BurstMode); err != nil {
		log.Fatalf("error setting the operation mode to burst: %s", err)
	}

	// Enable the output.
	if err = ch.EnableOutput(); err != nil {
		log.Fatalf("error enabling the output: %s", err)
	}

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

	// Query the DC offset voltage.
	offset, err := ch.DCOffset()
	if err != nil {
		log.Printf("error querying DC offset: %s", err)
	}
	log.Printf("DC Offset = %.1f mV", 1000*offset)

	// Query the standard waveform.
	wave, err := ch.StandardWaveform()
	if err != nil {
		log.Printf("error querying standard waveform: %s", err)
	}
	log.Printf("Standard waveform = %s", wave)

	// Query the burst count.
	bc, err := ch.BurstCount()
	if err != nil {
		log.Printf("error querying burst count: %s", err)
	}
	log.Printf("Burst count = %d", bc)

	// Query the internal trigger rate.
	itr, err := inst.InternalTriggerRate()
	if err != nil {
		log.Printf("error querying internal trigger rate: %s", err)
	}
	log.Printf("Internal trigger rate = %.1f Hz", itr)

	// Query the trigger source.
	ts, err := ch.StartTriggerSource()
	if err != nil {
		log.Printf("error querying start trigger source: %s", err)
	}
	log.Printf("Start trigger source = %s", ts)

	// Query the operation mode.
	om, err := ch.OperationMode()
	if err != nil {
		log.Printf("error querying operation mode: %s", err)
	}
	log.Printf("Operation mode = %s", om)

	// Close the VISA resource.
	err = res.Close()
	if err != nil {
		log.Printf("Error closing VISA resource: %s", err)
	}
}
