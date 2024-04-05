module github.com/gotmc/ivi-examples

go 1.21

require (
	github.com/gotmc/ivi v0.7.1
	github.com/gotmc/lxi v0.3.1
	github.com/gotmc/prologix v0.5.0
	github.com/gotmc/usbtmc v0.8.0
	github.com/gotmc/visa v0.7.0
)

require (
	github.com/creack/goselect v0.1.2 // indirect
	github.com/google/gousb v1.1.3 // indirect
	github.com/gotmc/convert v0.2.0 // indirect
	github.com/gotmc/query v0.4.0 // indirect
	go.bug.st/serial v1.6.0 // indirect
	golang.org/x/sys v0.0.0-20220829200755-d48e67d00261 // indirect
)

replace github.com/gotmc/ivi => ../ivi

replace github.com/gotmc/prologix => ../prologix
