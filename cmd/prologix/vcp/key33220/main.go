// Copyright (c) 2017-2020 The ivi developers. All rights reserved.
// Project site: https://github.com/gotmc/ivi
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"log"
	"time"

	"github.com/gotmc/ivi/fgen"
	"github.com/gotmc/ivi/fgen/agilent/ag33220"
	"github.com/gotmc/prologix"
	"github.com/tarm/serial"
)

func main() {

	// Open a serial port.
	cfg := serial.Config{
		Name:        "/dev/tty.usbserial-PX8X3YR6",
		Baud:        115200,
		ReadTimeout: time.Millisecond * 500,
	}
	port, err := serial.OpenPort(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new GPIB controller using the aforementioned serial port and
	// communicating with the instrument at GPIB address 4.
	gpib, err := prologix.NewController(port, 4, true)
	if err != nil {
		log.Fatalf("NewController error: %s", err)
	}
	prologixVer, err := gpib.Version()
	if err != nil {
		log.Fatalf("Unable to determine Prologix controller version: %s", err)
	}
	log.Printf("Using %s", prologixVer)

	// Create a new IVI instance of the Agilent 33220 function generator
	fg, err := ag33220.New(gpib, true)
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
	ch.ConfigureStandardWaveform(fgen.RampUp, 0.4, 0.1, 2340, 0)
	ch.EnableOutput()

	// Query the fg
	freq, err := ch.Frequency()
	if err != nil {
		log.Printf("error querying frequency: %s", err)
	}
	log.Printf("Frequency = %.0f Hz", freq)
	amp, err := ch.Amplitude()
	if err != nil {
		log.Printf("error querying amplitude: %s", err)
	}
	log.Printf("Amplitude = %.3f Vpp", amp)
	wave, err := ch.StandardWaveform()
	if err != nil {
		log.Printf("error querying standard waveform: %s", err)
	}
	log.Printf("Standard waveform = %s", wave)
	mfr, err := fg.InstrumentManufacturer()
	if err != nil {
		log.Printf("error querying instrument manufacturer: %s", err)
	}
	log.Printf("Instrument manufacturer = %s", mfr)
	model, err := fg.InstrumentModel()
	if err != nil {
		log.Printf("error querying instrument model: %s", err)
	}
	log.Printf("Instrument model = %s", model)
	sn, err := fg.InstrumentSerialNumber()
	if err != nil {
		log.Printf("error querying instrument sn: %s", err)
	}
	log.Printf("Instrument S/N = %s", sn)
	fw, err := fg.FirmwareRevision()
	if err != nil {
		log.Printf("error querying firmware revision: %s", err)
	}
	log.Printf("Firmware revision = %s", fw)

	// Return local control to the front panel.
	err = gpib.FrontPanel(true)
	if err != nil {
		log.Fatalf("error setting local control for front panel: %s", err)
	}

	// Discard any unread data on the serial port and then close.
	err = port.Flush()
	if err != nil {
		log.Printf("error flushing serial port: %s", err)
	}
	err = port.Close()
	if err != nil {
		log.Printf("error closing serial port: %s", err)
	}
}
