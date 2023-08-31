module github.com/gotmc/ivi-examples

go 1.21

require (
	github.com/gotmc/ivi v0.7.1
	github.com/gotmc/lxi v0.3.1
	github.com/gotmc/prologix v0.4.0
	github.com/gotmc/usbtmc v0.5.1
	github.com/gotmc/visa v0.5.0
	github.com/tarm/serial v0.0.0-20180830185346-98f6abe2eb07
)

require (
	github.com/google/gousb v0.0.0-20190812193832-18f4c1d8a750 // indirect
	github.com/gotmc/convert v0.2.0 // indirect
	github.com/gotmc/query v0.4.0 // indirect
	golang.org/x/sys v0.0.0-20210823070655-63515b42dcdf // indirect
)

replace github.com/gotmc/ivi => ../ivi

replace github.com/gotmc/prologix => ../prologix
