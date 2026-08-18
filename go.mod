module github.com/gotmc/ivi-examples

go 1.26

// Uncomment when developing against an unreleased local copy of gotmc/ivi.
// The E36102B example needs the e36000 single-output command set, which is
// not in a tagged release yet. Re-comment this and bump the require line once
// ivi is tagged.
replace github.com/gotmc/ivi => ../ivi

require (
	github.com/gotmc/asrl v0.14.0
	github.com/gotmc/ivi v0.30.0
	github.com/gotmc/lxi v0.17.0
	github.com/gotmc/prologix v0.11.0
	github.com/gotmc/usbtmc v0.15.1
	github.com/gotmc/visa v0.16.0
)

require (
	github.com/creack/goselect v0.1.3 // indirect
	github.com/google/gousb v1.1.3 // indirect
	github.com/gotmc/convert v0.5.1 // indirect
	github.com/gotmc/query v0.7.1 // indirect
	go.bug.st/serial v1.6.4 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
)
