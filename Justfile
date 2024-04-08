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

# Build and run the LXI Keysight/Agilent 33220A example application.
ex1 ip:
  #!/usr/bin/env bash
  echo '# LXI Keysight 33220A Example Application'
  cd {{justfile_directory()}}/cmd/lxi/key33220
  env go build -o key33220
  ./key33220 -ip={{ip}}

# Build and run the USBTMC Keysight/Agilent 33220A example application.
ex2:
  #!/usr/bin/env bash
  echo '# USBTMC Keysight 33220A Example Application'
  cd {{justfile_directory()}}/cmd/usbtmc/key33220
  env go build -o key33220
  ./key33220 -visa="USB0::2391::1031::MY44035849::INSTR"

# Build and run the VISA USBTMC Keysight/Agilent 33220A example application.
ex3:
  #!/usr/bin/env bash
  echo '# VISA USBTMC Keysight 33220A Example Application'
  cd {{justfile_directory()}}/cmd/visa/usbtmc/key33220
  env go build -o key33220
  ./key33220 -visa="USB0::2391::1031::MY44035849::INSTR"
