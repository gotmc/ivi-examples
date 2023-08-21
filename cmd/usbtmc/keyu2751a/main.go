// Copyright (c) 2017-2020 The ivi developers. All rights reserved.
// Project site: https://github.com/gotmc/ivi
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"fmt"
	"log"

	"github.com/gotmc/ivi/swtch/keysight/u2751a"
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
	dev, err := ctx.NewDeviceByVIDPID(0x0957, 0x3D18)
	if err != nil {
		log.Fatalf("NewDevice error: %s", err)
	}
	defer dev.Close()

	// Create a new IVI instance of the Keysight U2751A switch matrix.
	reset := true
	standalone := false
	sw, err := u2751a.New(dev, reset, standalone)
	if err != nil {
		log.Fatalf("IVI instrument error: %s", err)
	}

	numChannels := sw.ChannelCount()
	log.Printf("U2751A has %d channels", numChannels)

	// Determine instrument model.
	model, err := sw.InstrumentModel()
	if err != nil {
		log.Printf("error querying instrument model: %s", err)
	}
	log.Printf("Instrument model = %s", model)

	// Get a channel by ID and determine the wiremode.
	idx := 0
	ch, err := sw.ChannelByID(idx)
	if err != nil {
		log.Fatalf("Could not find channel %d: %s", idx, err)
	}
	wireMode := ch.WireMode()
	log.Printf("Channel %d contains %d wires.", idx, wireMode)

	// Set some virtual names.
	vn := map[string]string{
		"Row1": "dmmblack",
		"Row2": "dmmred",
		"Col1": "pin1",
		"Col2": "pin2",
	}

	err = sw.SetVirtualNames(vn)
	if err != nil {
		log.Fatal(err)
	}

	// Get a row and a column and set the row to a source channel.
	row1, err := sw.Channel("Row1")
	if err != nil {
		log.Fatal(err)
	}
	row1.SetSourceChannel(true)
	col2, err := sw.Channel("Col2")
	if err != nil {
		log.Fatal(err)
	}
	col2.SetSourceChannel(false)
	log.Printf("Row1 is source channel: %t", row1.IsSourceChannel())
	log.Printf("Col2 is source channel: %t", col2.IsSourceChannel())

	// Make a connection
	err = sw.Connect("Row1", "Col2")
	if err != nil {
		log.Fatalf("could not connect Row1 and Col2: %s", err)
	}
	log.Printf("Connected Row1 (source channel) and Col2 (non-source channel)")

	// Try to make an invalid connection.
	err = sw.Connect("Col1", "Col2")
	if err != nil {
		log.Printf("error trying to connect Col1 and Col2: %s", err)
	}

	// Determine the relay cycle counts
	rows := []string{"101:108", "201:208", "301:308", "401:408"}
	for i, row := range rows {
		q := fmt.Sprintf("diag:rel:cycl? (@%s)", row)
		resp, err := dev.Query(q)
		if err != nil {
			log.Printf("error querying relay cycle counts on row %d: %s", i+1, err)
		}
		log.Printf("Row %d cycle counts = %s", i+1, resp)
	}

}
