# -*- Justfile -*-

coverage_file := "coverage.out"

# List the available justfile recipes.
[group('general')]
@default:
  just --list --unsorted

# List the lines of code in the project.
[group('general')]
loc:
  scc --remap-unknown "-*- Justfile -*-":"justfile"

# View documentation in web browser using pkgsite.
[group('general')]
docs:
  pkgsite -open .

# Format and vet Go code. Runs before tests.
[group('test')]
check:
	go fix ./...
	go fmt ./...
	go vet ./...

# Lint using golangci-lint
[group('test')]
lint:
  golangci-lint run --config .golangci.yaml

# Run the unit tests.
[group('test')]
unit *FLAGS: check
  go test ./... -cover -vet=off -race {{FLAGS}} -short

# Run the integration tests.
[group('test')]
int *FLAGS: check
  go test ./... -cover -vet=off -race {{FLAGS}} -run Integration

# Run the end-to-end tests.
[group('test')]
e2e *FLAGS: check
  go test ./... -cover -vet=off -race {{FLAGS}} -run E2E

# HTML report for unit (default), int, e2e, or all tests.
[group('test')]
cover test='unit': check
  go test ./... -vet=off -coverprofile={{coverage_file}} \
  {{ if test == 'all' { '' } \
    else if test == 'int' { '-run Integration' } \
    else if test == 'e2e' { '-run E2E' } \
    else { '-short' } }}
  go tool cover -html={{coverage_file}}

# List the outdated direct dependencies (slow to run).
[group('dependencies')]
outdated:
  # (requires https://github.com/psampaz/go-mod-outdated).
  go list -u -m -json all | go-mod-outdated -update -direct

# Update the given module to the latest version.
[group('dependencies')]
update mod:
  go get -u {{mod}}
  go mod tidy

# Update all modules.
[group('dependencies')]
updateall:
  go get -u ./...
  go mod tidy

# Run go mod tidy and verify.
[group('dependencies')]
tidy:
  go mod tidy
 
# Run the ASRL SRS DS345 function generator example.
[group('examples')]
ds345 port:
  #!/usr/bin/env bash
  echo '# IVI ASRL SRS DS345 Example Application'
  cd {{justfile_directory()}}/cmd/asrl/ds345
  env go build -o ds345
  ./ds345 -port={{port}}

# Run the LXI Keysight 33512B fcn gen example.
[group('examples')]
k33512lxi ip:
  #!/usr/bin/env bash
  echo '# IVI LXI Keysight 33512B Example Application'
  cd {{justfile_directory()}}/cmd/lxi/key33512
  env go build -o key33512
  ./key33512 -ip={{ip}}

# Run the LXI Keysight 33220A fcn gen example.
[group('examples')]
k33220lxi ip:
  #!/usr/bin/env bash
  echo '# IVI LXI Keysight 33220A Example Application'
  cd {{justfile_directory()}}/cmd/lxi/key33220
  env go build -o key33220
  ./key33220 -ip={{ip}}

# Run the USBTMC Keysight 33220A fcn gen example.
[group('examples')]
k33220usb:
  #!/usr/bin/env bash
  echo '# IVI USBTMC Keysight 33220A Example Application'
  cd {{justfile_directory()}}/cmd/usbtmc/key33220
  env go build -o key33220
  ./key33220 -visa="USB0::2391::1031::MY44035849::INSTR"

# Run the VISA USBTMC Keysight 33220A fcn gen example.
[group('examples')]
k33220visa:
  #!/usr/bin/env bash
  echo '# IVI VISA USBTMC Keysight 33220A Example Application'
  cd {{justfile_directory()}}/cmd/visa/usbtmc/key33220
  env go build -o key33220
  ./key33220 -visa="USB0::2391::1031::MY44035849::INSTR"

# Run the Prologix VCP GPIB Keysight 33220A fcn gen example.
[group('examples')]
k33220gpib port:
  #!/usr/bin/env bash
  echo '# IVI Prologix VCP GPIB Keysight 33220A Example Application'
  cd {{justfile_directory()}}/cmd/prologix/vcp/key33220
  env go build -o key33220
  ./key33220 -port={{port}}

# Run the LXI Keysight 34461A DMM example.
[group('examples')]
k34461lxi ip:
  #!/usr/bin/env bash
  echo '# IVI LXI Keysight 34461A Example Application'
  cd {{justfile_directory()}}/cmd/lxi/key34461a
  env go build -o key34461a
  ./key34461a -ip={{ip}}

# Run the Prologix VCP GPIB Keysight E3631A pwr supply example.
[group('examples')]
k3631gpib port:
  #!/usr/bin/env bash
  echo '# IVI Prologix VCP GPIB Keysight E3631A Example Application'
  cd {{justfile_directory()}}/cmd/prologix/vcp/e3631a
  env go build -o e3631a
  ./e3631a -port={{port}}

# Run the ASRL Keysight E3631A pwr supply example.
[group('examples')]
k3631asrl port:
  #!/usr/bin/env bash
  echo '# IVI ASRL Keysight E3631A Example Application'
  cd {{justfile_directory()}}/cmd/asrl/e3631a
  env go build -o e3631a
  ./e3631a -port={{port}}

# Run the LXI Keysight InfiniiVision MSO-X 3024A oscilloscope example.
[group('examples')]
k3024lxi ip:
  #!/usr/bin/env bash
  echo '# IVI LXI Keysight InfiniiVision MSO-X 3024A Example Application'
  cd {{justfile_directory()}}/cmd/lxi/key3024
  env go build -o key3024
  ./key3024 -ip={{ip}}
