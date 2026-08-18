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
	"time"

	"github.com/gotmc/ivi"
	"github.com/gotmc/ivi/dcpwr"
	"github.com/gotmc/ivi/dcpwr/keysight/e36000"
	"github.com/gotmc/lxi"
)

func main() {
	log.Println("IVI LXI Keysight E36102B Example Application")

	// Get the IP address and I/O timeout from CLI flags.
	var ip string
	var timeout time.Duration
	flag.StringVar(
		&ip,
		"ip",
		"192.168.1.100",
		"IP address of Keysight E36102B",
	)
	flag.DurationVar(
		&timeout,
		"timeout",
		5*time.Second,
		"I/O timeout applied to each instrument operation",
	)
	flag.Parse()

	// Bound the initial TCP dial with the same timeout so an unreachable
	// instrument fails fast instead of hanging on the default OS connect
	// timeout.
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create a new LXI device. The E36100B series listens for raw SCPI socket
	// sessions on TCP port 5025.
	address := fmt.Sprintf("TCPIP0::%s::5025::SOCKET", ip)
	log.Printf("VISA address = %s", address)
	log.Printf("I/O timeout = %s", timeout)
	dev, err := lxi.NewDevice(ctx, address)
	if err != nil {
		log.Fatalf("NewDevice error: %s", err)
	}

	// Close the LXI device when done.
	defer func() {
		if err := dev.Close(); err != nil {
			log.Printf("error closing LXI device: %s", err)
		}
	}()

	// Create a new IVI instance of and reset the Keysight E36102B DC power
	// supply using the LXI device. ivi.WithTimeout applies the same timeout to
	// each subsequent driver method call.
	ps, err := e36000.New(dev, ivi.WithReset(), ivi.WithTimeout(timeout))
	if err != nil {
		log.Fatalf("IVI instrument error: %s", err)
	}
	defer func() {
		if err := ps.Close(); err != nil {
			log.Printf("error closing IVI driver: %s", err)
		}
	}()

	// From here forward, we can use the IVI API for the DC power supply instead
	// of having to send SCPI or other commands that are specific to this model
	// power supply.

	// Query the instrument manufacturer.
	mfr, err := ps.InstrumentManufacturer()
	if err != nil {
		log.Printf("error querying instrument manufacturer: %s", err)
	}
	log.Printf("Instrument manufacturer = %s", mfr)

	// Query the instrument model.
	model, err := ps.InstrumentModel()
	if err != nil {
		log.Printf("error querying instrument model: %s", err)
	}
	log.Printf("Instrument model = %s", model)

	// Query the instrument's serial number.
	sn, err := ps.InstrumentSerialNumber()
	if err != nil {
		log.Printf("error querying instrument sn: %s", err)
	}
	log.Printf("Instrument S/N = %s", sn)

	// Query the firmware revision.
	fw, err := ps.FirmwareRevision()
	if err != nil {
		log.Printf("error querying firmware revision: %s", err)
	}
	log.Printf("Firmware revision = %s", fw)

	// The E36102B is a single output, 6 V / 5 A supply, so the driver reports
	// one output channel named "Output".
	log.Printf("Output channel count = %d", ps.OutputChannelCount())

	// Channel specific methods are accessed using the Channel method with a
	// 0-based index to select the desired output channel.
	ch, err := ps.Channel(0)
	if err != nil {
		log.Fatalf("error getting channel 0: %s", err)
	}
	log.Printf("Configuring channel %q", ch.Name())

	// Turn the output off while configuring it.
	if err = ch.DisableOutput(); err != nil {
		log.Fatalf("error disabling output: %s", err)
	}

	// Set the output to 3.3 V and regulate at a 0.5 A current limit rather
	// than tripping the output off when the load reaches it.
	const (
		desiredVoltage = 3.3
		currentLimit   = 0.5
		ovpLimit       = 4.0
	)
	if err = ch.SetVoltageLevel(desiredVoltage); err != nil {
		log.Fatalf("error setting voltage level: %s", err)
	}
	if err = ch.ConfigureCurrentLimit(dcpwr.CurrentRegulate, currentLimit); err != nil {
		log.Fatalf("error configuring current limit: %s", err)
	}
	// The above call is the same as the following two:
	// ch.SetCurrentLimitBehavior(dcpwr.CurrentRegulate)
	// ch.SetCurrentLimit(currentLimit)

	// Arm over-voltage protection above the programmed level so the supply
	// shuts the output down if it ever runs away.
	if err = ch.ConfigureOVP(true, ovpLimit); err != nil {
		log.Fatalf("error configuring OVP: %s", err)
	}

	if err = ch.EnableOutput(); err != nil {
		log.Fatalf("error enabling output: %s", err)
	}

	// Let the power supply settle before measuring the output.
	time.Sleep(500 * time.Millisecond)

	// Query the programmed settings.
	v, err := ch.VoltageLevel()
	if err != nil {
		log.Printf("error querying voltage level: %s", err)
	}
	log.Printf("Voltage level = %.3f V", v)

	curr, err := ch.CurrentLimit()
	if err != nil {
		log.Printf("error querying current limit: %s", err)
	}
	log.Printf("Current limit = %.3f A", curr)

	behavior, err := ch.CurrentLimitBehavior()
	if err != nil {
		log.Printf("error querying current limit behavior: %s", err)
	}
	log.Printf("Current limit behavior = %s", behavior)

	ovpEnabled, err := ch.OVPEnabled()
	if err != nil {
		log.Printf("error querying OVP enabled: %s", err)
	}
	ovp, err := ch.OVPLimit()
	if err != nil {
		log.Printf("error querying OVP limit: %s", err)
	}
	log.Printf("OVP = %.3f V, enabled = %t", ovp, ovpEnabled)

	enabled, err := ch.OutputEnabled()
	if err != nil {
		log.Printf("error querying output enabled: %s", err)
	}
	log.Printf("Output enabled = %t", enabled)

	// Measure the actual output.
	vMsr, err := ch.MeasureVoltage()
	if err != nil {
		log.Printf("error measuring the voltage: %s", err)
	}
	log.Printf("Measured voltage = %.3f V", vMsr)

	cMsr, err := ch.MeasureCurrent()
	if err != nil {
		log.Printf("error measuring the current: %s", err)
	}
	log.Printf("Measured current = %.3f A", cMsr)

	// Report which regulation mode the output settled into, read from the
	// instrument's operation status register.
	for _, state := range []dcpwr.OutputState{dcpwr.ConstantVoltage, dcpwr.ConstantCurrent} {
		active, err := ch.QueryOutputState(state)
		if err != nil {
			log.Printf("error querying the %s state: %s", state, err)

			continue
		}
		log.Printf("Output state %s = %t", state, active)
	}

	// Leave the supply in a safe state with the output off.
	if err = ch.DisableOutput(); err != nil {
		log.Printf("error disabling output: %s", err)
	}
}
