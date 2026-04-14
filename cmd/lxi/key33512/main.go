// Copyright (c) 2017-2026 The ivi-examples developers. All rights reserved.
// Project site: https://github.com/gotmc/ivi-examples
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/gotmc/ivi"
	"github.com/gotmc/ivi/fgen"
	"github.com/gotmc/ivi/fgen/keysight/key33000"
	"github.com/gotmc/lxi"
)

func main() {
	log.Println("IVI LXI Keysight 33512B Example Application")

	// Get IP address from CLI flag.
	var ip string
	flag.StringVar(
		&ip,
		"ip",
		"192.168.1.100",
		"IP address of Keysight 33512B",
	)
	flag.Parse()

	ctx := context.Background()

	// Create a new LXI device
	address := fmt.Sprintf("TCPIP0::%s::5025::SOCKET", ip)
	log.Printf("VISA address = %s", address)
	dev, err := lxi.NewDevice(ctx, address)
	if err != nil {
		log.Fatalf("NewDevice error: %s", err)
	}

	// Close the LXI device when done.
	defer dev.Close()

	// Create a new IVI instance and reset the Keysight 33512B function generator
	// using the LXI device.
	fg, err := key33000.New(dev, ivi.WithIDQuery(), ivi.WithReset())
	if err != nil {
		log.Fatalf("IVI instrument error: %s", err)
	}

	// Query the instrument identification.
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

	// --- Channel 1: Configure a 1 kHz sine wave ---

	ch1, err := fg.Channel(0)
	if err != nil {
		log.Fatalf("error getting channel 0: %s", err)
	}
	if err = ch1.DisableOutput(); err != nil {
		log.Fatalf("error disabling output on ch1: %s", err)
	}

	if err = ch1.ConfigureStandardWaveform(fgen.Sine, 0.5, 0.0, 100.0, 0.0); err != nil {
		log.Fatalf("error configuring standard waveform on ch1: %s", err)
	}

	// Configure a burst waveform using the above 100 Hz sine wave with 400 ms
	// on-time and 200 ms off-time for a total period of 600 ms.
	if err = ch1.SetOperationMode(fgen.BurstMode); err != nil {
		log.Fatalf("error setting burst mode: %s", err)
	}

	if err = ch1.SetBurstCount(4); err != nil {
		log.Fatalf("error setting burst count: %s", err)
	}

	if err = ch1.SetStartTriggerSource(fgen.TriggerSourceInternal); err != nil {
		log.Fatalf("error setting internal trigger source: %s", err)
	}

	if err = fg.SetInternalTriggerRate(1 / 0.06); err != nil {
		log.Fatalf("error setting internal trigger rate: %s", err)
	}

	if err = ch1.EnableOutput(); err != nil {
		log.Fatalf("error enabling output on ch1: %s", err)
	}

	// --- Channel 2: Configure a 500 Hz square wave ---

	ch2, err := fg.Channel(1)
	if err != nil {
		log.Fatalf("error getting channel 1: %s", err)
	}
	if err = ch2.DisableOutput(); err != nil {
		log.Fatalf("error disabling output on ch2: %s", err)
	}

	if err = ch2.SetStandardWaveform(fgen.Square); err != nil {
		log.Fatalf("error setting waveform on ch2: %s", err)
	}
	if err = ch2.SetFrequency(500); err != nil {
		log.Fatalf("error setting frequency on ch2: %s", err)
	}
	if err = ch2.SetAmplitude(2.0); err != nil {
		log.Fatalf("error setting amplitude on ch2: %s", err)
	}
	if err = ch2.SetDCOffset(0.5); err != nil {
		log.Fatalf("error setting DC offset on ch2: %s", err)
	}

	if err = ch2.EnableOutput(); err != nil {
		log.Fatalf("error enabling output on ch2: %s", err)
	}

	// --- Query both channels ---

	for i := range 2 {
		ch, err := fg.Channel(i)
		if err != nil {
			log.Fatalf("error getting channel %d: %s", i, err)
		}

		wave, err := ch.StandardWaveform()
		if err != nil {
			log.Printf("ch%d: error querying waveform: %s", i+1, err)
		}

		freq, err := ch.Frequency()
		if err != nil {
			log.Printf("ch%d: error querying frequency: %s", i+1, err)
		}

		amp, err := ch.Amplitude()
		if err != nil {
			log.Printf("ch%d: error querying amplitude: %s", i+1, err)
		}

		offset, err := ch.DCOffset()
		if err != nil {
			log.Printf("ch%d: error querying DC offset: %s", i+1, err)
		}

		enabled, err := ch.OutputEnabled()
		if err != nil {
			log.Printf("ch%d: error querying output enabled: %s", i+1, err)
		}

		log.Printf("CH%d: %s, %.0f Hz, %.3f Vpp, %.3f Vdc offset, enabled=%t",
			i+1, wave, freq, amp, offset, enabled)
	}
}
