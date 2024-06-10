// Copyright (c) 2017-2024 The ivi-examples developers. All rights reserved.
// Project site: https://github.com/gotmc/ivi-examples
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gotmc/ivi/scope"
	"github.com/gotmc/ivi/scope/keysight/infiniivision"
	"github.com/gotmc/lxi"
)

func main() {
	log.Println("IVI LXI Keysight InfiniiVision MSO-X 3024A Example Application")

	// Get IP address from CLI flag.
	var ip string
	flag.StringVar(
		&ip,
		"ip",
		"192.168.1.100",
		"IP address of Keysight InfiniiVision MSO-X 3024A",
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

	// Create a new IVI instance of and reset the Keysight InfiniiVision oscilloscope
	// using the LXI device.
	scope1, err := infiniivision.New(dev, true)
	if err != nil {
		log.Fatalf("IVI instrument eror: %s", err)
	}

	// From here forward, we can use the IVI API for the oscilloscope instead of
	// having to send SCPI or other commands that are specific to this model
	// oscilloscope.

	// Query the instrument manufacturer.
	mfr, err := scope1.InstrumentManufacturer()
	if err != nil {
		log.Printf("error querying instrument manufacturer: %s", err)
	}
	log.Printf("Instrument manufacturer = %s", mfr)

	// Query the instrument model.
	model, err := scope1.InstrumentModel()
	if err != nil {
		log.Printf("error querying instrument model: %s", err)
	}
	log.Printf("Instrument model = %s", model)

	// Query the instrument's serial number.
	sn, err := scope1.InstrumentSerialNumber()
	if err != nil {
		log.Printf("error querying instrument sn: %s", err)
	}
	log.Printf("Instrument S/N = %s", sn)

	// Query the firmware revision.
	fw, err := scope1.FirmwareRevision()
	if err != nil {
		log.Printf("error querying firmware revision: %s", err)
	}
	log.Printf("Firmware revision = %s", fw)

	i := scope1.ChannelCount()
	log.Printf("Channel count = %d", i)

	ch1 := scope1.Channels[0]
	if err = ch1.SetInputImpedance(50.0); err != nil {
		panic(err)
	}

	// Set the total vertical range to 800 mV for CH1.
	if err = ch1.SetVerticalRange(0.8); err != nil {
		panic(err)
	}

	// Set the vertical offset to 0 V for CH1.
	if err = ch1.SetVerticalOffset(0.0); err != nil {
		panic(err)
	}

	// Set the vertical coupling to DC for CH1.
	if err = ch1.SetVerticalCoupling(scope.DCVerticalCoupling); err != nil {
		panic(err)
	}

	// Set the probe attenuation to 1:1 for CH1.
	if err = ch1.SetProbeAttenuation(1.0); err != nil {
		panic(err)
	}

	// Disable CH1.
	if err = ch1.SetChannelEnabled(false); err != nil {
		panic(err)
	}

	// Set the vertical range, vertical offset, vertical coupling, disable auto
	// probe attentuation, set the probe attenuation, and enable the channel in
	// one command.
	if err = ch1.Configure(0.8, 0.0, scope.DCVerticalCoupling, false, 1.0, true); err != nil {
		panic(err)
	}

	// Query the acquisition type.
	acqType, err := scope1.AcquisitionType()
	if err != nil {
		panic(err)
	}
	log.Printf("Acquisition type = %v", acqType)

	// Query the acquisition record length.
	recordLength, err := scope1.AcquisitionRecordLength()
	if err != nil {
		panic(err)
	}
	log.Printf("Acquisition record length = %d", recordLength)

	// Query the acquisition sample rate.
	sampleRate, err := scope1.AcquisitionSampleRate()
	if err != nil {
		panic(err)
	}
	log.Printf("Acquisition sample rate = %g samples/sec", sampleRate)

	// Query the acquistion time per record.
	timePerRecord, err := scope1.AcquisitionTimePerRecord()
	if err != nil {
		panic(err)
	}
	log.Printf("Acquisition time per record = %s", timePerRecord)

	// Set the acquisition time per record.
	timePerRecord = 120 * time.Millisecond
	log.Printf("Setting the acquisition time per record to %s", timePerRecord)
	if err = scope1.SetAcquisitionTimePerRecord(timePerRecord); err != nil {
		panic(err)
	}

	// Query the acquistion time per record.
	timePerRecord, err = scope1.AcquisitionTimePerRecord()
	if err != nil {
		panic(err)
	}
	log.Printf("Acquisition time per record = %s", timePerRecord)

	// Set the trigger type to edge.
	log.Printf("Setting the trigger type to edge")
	if err = scope1.SetTriggerType(scope.EdgeTrigger); err != nil {
		panic(err)
	}

	// Query the trigger type.
	triggerType, err := scope1.TriggerType()
	if err != nil {
		panic(err)
	}
	log.Printf("Trigger type = %s", triggerType)

	// Set the trigger level to 100 mV.
	log.Printf("Setting the trigger level to 15 mV")
	if err = scope1.SetTriggerLevel(0.015); err != nil {
		panic(err)
	}

	// Query the trigger level.
	triggerLevel, err := scope1.TriggerLevel()
	if err != nil {
		panic(err)
	}
	log.Printf("Trigger level = %g V", triggerLevel)

	// Set the trigger holdoff to 40 ms since we have four bursts with 10 ms
	// period (1 / 100 Hz) from the function generator.
	if err = scope1.SetTriggerHoldoff(40 * time.Millisecond); err != nil {
		panic(err)
	}

	// // Set the trigger delay to 50 ms.
	// if err = scope1.SetAcquisitionStartTime(50 * time.Millisecond); err != nil {
	// 	panic(err)
	// }

	// Query the trigger delay.
	delay, err := scope1.AcquisitionStartTime()
	if err != nil {
		panic(err)
	}
	log.Printf("Acquisition start time delay = %s", delay)

	// Measure the peak-to-peak voltage
	vpp, err := ch1.FetchWaveformMeasurement(scope.VoltagePeakToPeak)
	if err != nil {
		panic(err)
	}
	log.Printf("Voltage peak-to-peak = %.3f Vpp", vpp)

	// Measure the frequency
	freq, err := ch1.FetchWaveformMeasurement(scope.Frequency)
	if err != nil {
		panic(err)
	}
	log.Printf("Frequency = %.3f Hz", freq)
	// ConfigureTrigger(triggerType TriggerType, holdoff time.Duration) error
}
