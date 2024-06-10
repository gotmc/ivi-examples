// Copyright (c) 2017-2024 The ivi-examples developers. All rights reserved.
// Project site: https://github.com/gotmc/ivi-examples
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gotmc/ivi/fgen"
	"github.com/gotmc/ivi/fgen/keysight/key33220"
	"github.com/gotmc/lxi"
)

func main() {
	log.Println("IVI LXI Keysight 33220A Example Application")

	// Get IP address from CLI flag.
	var ip string
	flag.StringVar(
		&ip,
		"ip",
		"192.168.1.100",
		"IP address of Keysight 33220A",
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
	fg, err := key33220.New(dev, true)
	if err != nil {
		log.Fatalf("IVI instrument eror: %s", err)
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
	if err = fg.Channels[0].DisableOutput(); err != nil {
		log.Fatalf("error disabling output on ch0: %s", err)
	}
	if err = fg.Channels[0].SetAmplitude(2.1); err != nil {
		log.Fatalf("error setting the amplitude on ch0: %s", err)
	}

	// Alternatively, the channel can be assigned to a variable.
	ch := fg.Channels[0]
	if err = ch.SetStandardWaveform(fgen.Sine); err != nil {
		log.Fatalf("error setting the standard waveform: %s", err)
	}
	if err = ch.SetDCOffset(0.1); err != nil {
		log.Fatalf("error setting DC offest: %s", err)
	}
	if err = ch.SetFrequency(2100); err != nil {
		log.Fatalf("error setting frequency: %s", err)
	}

	// Instead of configuring attributes of a standard waveform individually, the
	// standard waveform can be configured using a single method. In this case, a
	// Sine wave with 0.5 Vpp amplitude, 0.0 Vdc offset, 100.0 Hz, and 0.0 phase
	// shift is created.
	if err = ch.ConfigureStandardWaveform(fgen.Sine, 0.5, 0.0, 100.0, 0.0); err != nil {
		log.Fatalf("error configuring standard waveform: %s", err)
	}

	// Configure a burst waveform using the above 100 Hz sine wave with 400 ms
	// on-time and 200 ms off-time for a total period of 600 ms.
	if err = ch.SetOperationMode(fgen.BurstMode); err != nil {
		log.Fatalf("error setting burst mode: %s", err)
	}

	if err = ch.SetBurstCount(4); err != nil {
		log.Fatalf("error setting burst count: %s", err)
	}

	if err = ch.SetStartTriggerSource(fgen.TriggerSourceInternal); err != nil {
		log.Fatalf("error setting internal trigger source: %s", err)
	}

	if err = fg.SetInternalTriggerRate(1 / 0.06); err != nil {
		log.Fatalf("error setting internal trigger rate: %s", err)
	}

	if err = ch.EnableOutput(); err != nil {
		log.Fatalf("error enabling output: %s", err)
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
	itr, err := fg.InternalTriggerRate()
	if err != nil {
		log.Printf("error querying internal trigger rate: %s", err)
	}
	log.Printf("Internal trigger rate = %.1f Hz", itr)

	// Query the start trigger source.
	ts, err := ch.StartTriggerSource()
	if err != nil {
		log.Printf("error querying start trigger source: %s", err)
	}
	log.Printf("Start trigger source = %v", ts)

	// Query the operation mode.
	om, err := ch.OperationMode()
	if err != nil {
		log.Printf("error querying operation mode: %s", err)
	}
	log.Printf("Operation mode = %s", om)
}
