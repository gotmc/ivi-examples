// Copyright (c) 2017-2020 The ivi developers. All rights reserved.
// Project site: https://github.com/gotmc/ivi
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"log"

	"github.com/gotmc/ivi/fgen"
	"github.com/gotmc/ivi/fgen/keysight/key33220"
	_ "github.com/gotmc/usbtmc/driver/google"
	"github.com/gotmc/visa"
	_ "github.com/gotmc/visa/driver/usbtmc"
)

func main() {

	// Configure a new VISA resource using the USBTMC driver.
	address := "usb0::2391::1031::MY44035849::INSTR"
	res, err := visa.NewResource(address)
	if err != nil {
		log.Fatalf("VISA resource %s: %s", address, err)
	}

	// Setup the IVI driver for the Keysight 33220A function generator.
	fg, err := key33220.New(res, true)
	if err != nil {
		log.Fatalf("IVI instrument error: %s", err)
	}

	// From here forward, we can use the IVI API for the function generator
	// instead of having to send SCPI or other commands that are specific to this
	// model function generator.

	// Grab the output channel and disable the output.
	ch := fg.Channels[0]
	ch.DisableOutput()

	// Shortcut to configure standard waveform in one command.
	ch.ConfigureStandardWaveform(fgen.Sine, 0.25, 0.07, 2340, 0)

	// Setup a bursted sinusoidal waveform.
	ch.SetBurstCount(131)
	ch.SetInternalTriggerPeriod(0.112) // code period = 112 ms
	ch.SetTriggerSource(fgen.InternalTrigger)
	ch.SetOperationMode(fgen.BurstMode)

	// Enable the output.
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

	// Query the DC offset voltage.
	offset, err := ch.DCOffset()
	if err != nil {
		log.Printf("error querying DC offset: %s", err)
	}
	log.Printf("DC Offset = %.3f V", offset)

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
	log.Printf("Internal trigger period = %f", itp)

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

	// Close the VISA resource.
	err = res.Close()
	if err != nil {
		log.Printf("Error closing VISA resource: %s", err)
	}
}
