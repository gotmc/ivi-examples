# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This repository contains example applications demonstrating the [gotmc/ivi](https://github.com/gotmc/ivi) ecosystem for controlling test and measurement instruments from Go. Each example pairs an IVI instrument driver with a specific transport (LXI, USBTMC, VISA, Prologix GPIB, or serial/ASRL) and lives under `cmd/<transport>/<instrument>/main.go`.

## Common Commands

```bash
just check              # go fmt + go vet
just unit               # unit tests with -race -short -cover
just unit -run TestName # run a single test
just lint               # golangci-lint (requires .golangci.yaml)
just cover              # HTML coverage report
just tidy               # go mod tidy
just updateall          # update all dependencies
```

Run examples via Justfile recipes (each requires hardware or a valid network target):

```bash
just k33220lxi 192.168.1.100    # LXI Keysight 33220A function generator
just k34461lxi 10.12.100.56     # LXI Keysight 34461A DMM
just k33220usb                  # USBTMC Keysight 33220A
just k33220gpib /dev/tty.usbserial-PX8X3YR6  # Prologix GPIB
just ds345 /dev/tty.usbserial   # ASRL SRS DS345
```

Build a single example manually:

```bash
cd cmd/lxi/key33220 && go build -o key33220
```

## Architecture

### Example structure

Every example follows the same pattern:

1. Parse CLI flags for the connection target (IP address, serial port, or VISA address).
2. Create a transport-specific device/resource (e.g., `lxi.NewDevice`, `visa.NewResource`, `prologix.NewController`).
3. Wrap it in an IVI driver (e.g., `key33220.New(dev, true)` where `true` resets on init).
4. Use the transport-agnostic IVI API (`InstrumentManufacturer()`, `ch.SetFrequency()`, etc.).

### Transport patterns

- **LXI** (TCP/IP): Uses `lxi.NewDevice` with a VISA-style address string (`TCPIP0::<ip>::5025::SOCKET`).
- **USBTMC**: Uses `usbtmc.NewDevice` directly, requires `_ "github.com/gotmc/usbtmc/driver/google"` blank import for the USB backend.
- **VISA**: Uses `visa.NewResource` with blank imports for both the VISA driver (`_ "github.com/gotmc/visa/driver/usbtmc"`) and USB backend.
- **Prologix GPIB**: Creates a VCP serial connection, then wraps it with `prologix.NewController(vcp, gpibAddr, resetOnInit)`.
- **ASRL** (serial): Uses `asrl.NewDevice` with a serial port path.

### IVI instrument classes used

- `fgen` (Function Generator): `key33220` (Keysight 33220A), `ds345` (SRS DS345)
- `dmm` (Digital Multimeter): `key3446x` (Keysight 34461A)
- `dcpwr` (DC Power Supply): `e36xx` (Keysight E3631A)

## Go Version

Requires Go 1.25+ (see `go.mod`).
