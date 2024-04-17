# List the available justfile recipes.
@default:
  just --list

# Format, vet, and test Go code.
check:
	go fmt ./...
	go vet ./...
	GOEXPERIMENT=loopvar go test ./... -cover

# Verbosely format, vet, and test Go code.
checkv:
	go fmt ./...
	go vet ./...
	GOEXPERIMENT=loopvar go test -v ./... -cover

# Lint code using staticcheck.
lint:
	staticcheck -f stylish ./...

# Test and provide HTML coverage report.
cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

# List the outdated go modules.
outdated:
  go list -u -m all

# Run the LXI Keysight 33220A fcn gen example.
k33220lxi ip:
  #!/usr/bin/env bash
  echo '# IVI LXI Keysight 33220A Example Application'
  cd {{justfile_directory()}}/cmd/lxi/key33220
  env go build -o key33220
  ./key33220 -ip={{ip}}

# Run the USBTMC Keysight 33220A fcn gen example.
k33220usb:
  #!/usr/bin/env bash
  echo '# IVI USBTMC Keysight 33220A Example Application'
  cd {{justfile_directory()}}/cmd/usbtmc/key33220
  env go build -o key33220
  ./key33220 -visa="USB0::2391::1031::MY44035849::INSTR"

# Run the VISA USBTMC Keysight 33220A fcn gen example.
k33220visa:
  #!/usr/bin/env bash
  echo '# IVI VISA USBTMC Keysight 33220A Example Application'
  cd {{justfile_directory()}}/cmd/visa/usbtmc/key33220
  env go build -o key33220
  ./key33220 -visa="USB0::2391::1031::MY44035849::INSTR"

# Run the Prologix VCP GPIB Keysight 33220A fcn gen example.
k33220gpib port:
  #!/usr/bin/env bash
  echo '# IVI Prologix VCP GPIB Keysight 33220A Example Application'
  cd {{justfile_directory()}}/cmd/prologix/vcp/key33220
  env go build -o key33220
  ./key33220 -port={{port}}

# Run the LXI Keysight 34461A DMM example.
k34461lxi ip:
  #!/usr/bin/env bash
  echo '# IVI LXI Keysight 34461A Example Application'
  cd {{justfile_directory()}}/cmd/lxi/key34461a
  env go build -o key34461a
  ./key34461a -ip={{ip}}

# Run the Prologix VCP GPIB Keysight E3631A pwr supply example.
k3631gpib port:
  #!/usr/bin/env bash
  echo '# IVI Prologix VCP GPIB Keysight E3631A Example Application'
  cd {{justfile_directory()}}/cmd/prologix/vcp/e3631a
  env go build -o e3631a
  ./e3631a -port={{port}}

# Run the ASRL Keysight E3631A pwr supply example.
k3631asrl port:
  #!/usr/bin/env bash
  echo '# IVI ASRL Keysight E3631A Example Application'
  cd {{justfile_directory()}}/cmd/asrl/e3631a
  env go build -o e3631a
  ./e3631a -port={{port}}

