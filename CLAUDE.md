# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This repository contains example applications demonstrating the [gotmc/ivi](https://github.com/gotmc/ivi) ecosystem for controlling test and measurement instruments from Go. Each example pairs an IVI instrument driver with a specific transport (LXI, USBTMC, VISA, Prologix GPIB, or serial/ASRL) and lives under `cmd/<transport>/<instrument>/main.go`.

Example directories are named for the *instrument*, not the driver package. A single driver often covers many models, so `cmd/lxi/kt33220` and `cmd/lxi/kt33512` both import `fgen/keysight/kt33000`.

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
just k33512lxi 192.168.1.100    # LXI Keysight 33512B function generator
just k34461lxi 10.12.100.56     # LXI Keysight 34461A DMM
just k3024lxi 192.168.1.100     # LXI Keysight MSO-X 3024A oscilloscope
just k36102lxi 192.168.1.100    # LXI Keysight E36102B DC power supply
just pmxlxi 192.168.1.100       # LXI Kikusui PMX DC power supply
just k33220usb MY44035849       # USBTMC Keysight 33220A (by serial number)
just ku2751usb                  # USBTMC Keysight U2751A switch matrix
just k33220visa                 # VISA USBTMC Keysight 33220A
just k33220gpib /dev/tty.usbserial-PX8X3YR6  # Prologix GPIB Keysight 33220A
just k3631gpib /dev/tty.usbserial             # Prologix GPIB Keysight E3631A
just f45gpib /dev/tty.usbserial               # Prologix GPIB Fluke 45 DMM
just k3631asrl /dev/tty.usbserial             # ASRL Keysight E3631A
just ds345 /dev/tty.usbserial   # ASRL SRS DS345
```

Build a single example manually. Each recipe builds into the example's own
directory using the directory name as the binary name, and `.gitignore` lists
those paths:

```bash
cd cmd/lxi/kt33220 && go build -o kt33220
```

## Depending on gotmc/ivi

`go.mod` carries a commented-out replace directive:

```
// replace github.com/gotmc/ivi => ../ivi
```

Uncomment it to build against an unreleased local ivi checkout, which is
necessary whenever an example uses driver packages or API that are not in a
tagged release yet. Re-comment it and bump the `require` line once ivi is
tagged, so the repository builds without a sibling checkout.

## Architecture

### Example structure

Every example follows the same pattern:

1. Parse CLI flags for the connection target (IP address, serial port, or VISA address).
2. Create a transport-specific device/resource (e.g., `lxi.NewDevice`, `visa.NewResource`, `prologix.NewController`).
3. Wrap it in an IVI driver (e.g., `kt33000.New(dev, ivi.WithReset())`).
4. Use the transport-agnostic IVI API (`InstrumentManufacturer()`, `ch.SetFrequency()`, etc.).

### Transport patterns

- **LXI** (TCP/IP): Uses `lxi.NewDevice` with a VISA-style address string (`TCPIP0::<ip>::5025::SOCKET`).
- **USBTMC**: Uses `usbtmc.NewDevice` directly, requires `_ "github.com/gotmc/usbtmc/driver/google"` blank import for the USB backend. Devices are opened by VID/PID via `usbCtx.NewDeviceByVIDPID()`.
- **VISA**: Uses `visa.NewResource` with blank imports for both the VISA driver (`_ "github.com/gotmc/visa/driver/usbtmc"`) and USB backend.
- **Prologix GPIB**: Creates a VCP serial connection via `vcp.NewVCP(serialPort)`, then wraps it with `prologix.NewController(vcp, gpibAddr, resetOnInit)`. Requires cleanup: `vcp.Flush()` then `vcp.Close()`.
- **ASRL** (serial): Uses `asrl.NewDevice` with a serial port path.

### IVI driver constructors

All IVI drivers use functional options: `driver.New(dev, opts ...ivi.DriverOption)`.

The `*IDN?` query and model validation happen by default, so there is no
`WithIDQuery` option; pass `ivi.WithoutIDQuery()` to skip the check. The other
options are `ivi.WithReset()` to send `*RST` on creation,
`ivi.WithTimeout(d)` to override the default per-call I/O timeout, and
`ivi.WithStandalone()`.

The `cmd/asrl/e3631a` and `cmd/lxi/kt34461a` examples demonstrate the
`-timeout` flag paired with `ivi.WithTimeout`.

### IVI instrument classes used

- `fgen` (Function Generator): `kt33000` (Keysight 33220A and 33512B), `ds345` (SRS DS345)
- `dmm` (Digital Multimeter): `kt34400` (Keysight 34461A), `fluke45` (Fluke 45)
- `dcpwr` (DC Power Supply): `e36000` (Keysight E3631A and E36102B), `pmx`
  (Kikusui PMX)
- `scope` (Oscilloscope): `infiniivision` (Keysight MSO-X 3024A)
- `swtch` (Switch): `u2751a` (Keysight U2751A)

Driver packages are renamed upstream from time to time. When an import stops
resolving, check the package directories under the ivi module rather than
assuming the example is wrong.

## Go Version

Requires Go 1.26+ (see `go.mod`).
