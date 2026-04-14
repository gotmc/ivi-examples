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
	"github.com/gotmc/ivi/dcpwr/kikusui/pmx"
	"github.com/gotmc/lxi"
)

func main() {

	// Get IP address from CLI flag.
	var ip string
	flag.StringVar(
		&ip,
		"ip",
		"192.168.1.100",
		"IP address of Kikusui PMX DC power supply",
	)
	flag.Parse()

	ctx := context.Background()

	// Create a new LXI device
	address := fmt.Sprintf("TCPIP0::%s::5025::SOCKET", ip)
	dev, err := lxi.NewDevice(ctx, address)
	if err != nil {
		log.Fatalf("NewDevice error: %s", err)
	}
	defer dev.Close()

	// Create a new IVI instance of the KIKUSUI PMW power supply and reset.
	dcp, err := pmx.New(dev, ivi.WithIDQuery(), ivi.WithReset())
	if err != nil {
		log.Fatalf("IVI instrument error: %s", err)
	}
	if err = dcp.Reset(); err != nil {
		log.Fatalf("error resetting instrument: %s", err)
	}

	// Get the first channel.
	ch, err := dcp.Channel(0)
	if err != nil {
		log.Fatalf("error getting channel 0: %s", err)
	}
	if err = ch.DisableOutput(); err != nil {
		log.Fatalf("error disabling output: %s", err)
	}
	if err = ch.SetVoltageLevel(50); err != nil {
		log.Fatalf("error setting voltage level: %s", err)
	}
	if err = ch.ConfigureCurrentLimit(dcpwr.CurrentTrip, 0.25); err != nil {
		log.Fatalf("error configuring current limit: %s", err)
	}
	// The above command is the same as the following two:
	// ch.SetCurrentLimitBehavior(dcpwr.Trip)
	// ch.SetCurrentLimit(0.25)
	if err = ch.ConfigureOVP(true, 60); err != nil {
		log.Fatalf("error configuring OVP: %s", err)
	}
	// The above command is the same as the following two:
	// ch.SetOVPEnabled(true)
	// ch.SetOVPLimit(60)
	if err = ch.EnableOutput(); err != nil {
		log.Fatalf("error enabling output: %s", err)
	}

	// Let the power supply settle before we query it.
	time.Sleep(500 * time.Millisecond)
	v, err := ch.VoltageLevel()
	if err != nil {
		log.Printf("error querying voltage level: %s", err)
	}
	log.Printf("Voltage limit = %.0f V", v)
	measured, err := ch.MeasureVoltage()
	if err != nil {
		log.Printf("error measuriing the voltage: %s", err)
	}
	log.Printf("Measured voltage = %.3f V", measured)

	// Get information about the power supply
	mfr, err := dcp.InstrumentManufacturer()
	if err != nil {
		log.Printf("error querying instrument manufacturer: %s", err)
	}
	log.Printf("Instrument manufacturer = %s", mfr)
	model, err := dcp.InstrumentModel()
	if err != nil {
		log.Printf("error querying instrument model: %s", err)
	}
	log.Printf("Instrument model = %s", model)
	sn, err := dcp.InstrumentSerialNumber()
	if err != nil {
		log.Printf("error querying instrument sn: %s", err)
	}
	log.Printf("Instrument S/N = %s", sn)
	fw, err := dcp.FirmwareRevision()
	if err != nil {
		log.Printf("error querying firmware revision: %s", err)
	}
	log.Printf("Firmware revision = %s", fw)
}
