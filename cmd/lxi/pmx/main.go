// Copyright (c) 2017-2024 The ivi-examples developers. All rights reserved.
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
	dcp, err := pmx.New(dev, true)
	if err != nil {
		log.Fatalf("IVI instrument error: %s", err)
	}
	dcp.Reset(ctx)

	// Alternatively, the channel can be assigned to a variable.
	ch := dcp.Channels[0]
	ch.DisableOutput(ctx)
	ch.SetVoltageLevel(ctx, 50)
	ch.ConfigureCurrentLimit(ctx, dcpwr.CurrentTrip, 0.25)
	// The above command is the same as the following two:
	// ch.SetCurrentLimitBehavior(ctx, dcpwr.Trip)
	// ch.SetCurrentLimit(ctx, 0.25)
	ch.ConfigureOVP(ctx, true, 60)
	// The aove command is the same as the following two:
	// ch.SetOVPEnabled(ctx, true)
	// ch.SetOVPLimit(ctx, 60)
	ch.EnableOutput(ctx)

	// Let the power supply settle before we query it.
	time.Sleep(500 * time.Millisecond)
	v, err := ch.VoltageLevel(ctx)
	if err != nil {
		log.Printf("error querying voltage level: %s", err)
	}
	log.Printf("Voltage limit = %.0f V", v)
	measured, err := ch.MeasureVoltage(ctx)
	if err != nil {
		log.Printf("error measuriing the voltage: %s", err)
	}
	log.Printf("Measured voltage = %.3f V", measured)

	// Get information about the power supply
	mfr, err := dcp.InstrumentManufacturer(ctx)
	if err != nil {
		log.Printf("error querying instrument manufacturer: %s", err)
	}
	log.Printf("Instrument manufacturer = %s", mfr)
	model, err := dcp.InstrumentModel(ctx)
	if err != nil {
		log.Printf("error querying instrument model: %s", err)
	}
	log.Printf("Instrument model = %s", model)
	sn, err := dcp.InstrumentSerialNumber(ctx)
	if err != nil {
		log.Printf("error querying instrument sn: %s", err)
	}
	log.Printf("Instrument S/N = %s", sn)
	fw, err := dcp.FirmwareRevision(ctx)
	if err != nil {
		log.Printf("error querying firmware revision: %s", err)
	}
	log.Printf("Firmware revision = %s", fw)
}
