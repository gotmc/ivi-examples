// Copyright (c) 2017-2024 The ivi-examples developers. All rights reserved.
// Project site: https://github.com/gotmc/ivi-examples
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"log"

	"github.com/gotmc/ivi/fgen"
	"github.com/gotmc/ivi/fgen/keysight/key33220"
	"github.com/gotmc/usbtmc"
	_ "github.com/gotmc/usbtmc/driver/google"
)

func main() {

	// Create a USBTMC context and set the debug level
	ctx, err := usbtmc.NewContext()
	if err != nil {
		log.Fatalf("Error creating new USB context: %s", err)
	}
	defer ctx.Close()
	ctx.SetDebugLevel(1)

	// Create a new USBTMC device
	dev, err := ctx.NewDevice("USB0::2391::1031::MY44035849::INSTR")
	if err != nil {
		log.Fatalf("NewDevice error: %s", err)
	}
	defer dev.Close()

	// Create a new IVI instance of the Agilent 33220 function generator
	fg, err := key33220.New(dev, true)
	if err != nil {
		log.Fatalf("IVI instrument error: %s", err)
	}

	// Channel specific methods can be accessed directly from the instrument
	// using 0-based index to select the desirec channel.
	fg.Channels[0].DisableOutput()
	fg.Channels[0].SetAmplitude(0.4)

	// Alternatively, the channel can be assigned to a variable.
	ch := fg.Channels[0]
	ch.SetStandardWaveform(fgen.Sine)
	ch.SetDCOffset(0.1)
	ch.SetFrequency(2340)

	// Instead of configuring attributes of a standard waveform individually, the
	// standard waveform can be configured using a single method.
	ch.ConfigureStandardWaveform(fgen.RampUp, 0.4, 0.1, 2360, 0)
	ch.EnableOutput()

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

	// Query the internal trigger period.
	itp, err := ch.InternalTriggerPeriod()
	if err != nil {
		log.Printf("error querying internal trigger period: %s", err)
	}
	log.Printf("Internal trigger period = %.1f ms", 1000*itp)

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

	// Query the serial number.
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
}
