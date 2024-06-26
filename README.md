# ivi-examples

Examples using the [gotmc ivi package][ivi], which proides a Go-based
implementation of the Interchangeable Virtual Instrument (IVI) standard.

[![GoDoc][godoc badge]][godoc link]
[![Go Report Card][report badge]][report card]
[![License Badge][license badge]][LICENSE.txt]

## Overview

The [IVI Specifications][ivi-specs] developed by the [IVI
Foundation][ivi-foundation] provide standardized APIs for programming test
instruments. The [ivi][] package is a partial, Go-based implementation of the
IVI Specifications, which are specified for C, COM, and .NET.

The main advantage of the [ivi][] package is not having to learn the [SCPI][]
commands for each individual piece of test equipment. For instance, by using the
[ivi][] package both the Agilent 33220A and the Stanford Research Systems DS345
function generators can be programmed using one standard API. The only
requirement for this is having an IVI driver for the desired test equipment.

## Documentation

Documentation can be found at either:

- <https://godoc.org/github.com/gotmc/ivi-examples>
- <http://localhost:6060/pkg/github.com/gotmc/ivi-examples/> after running `$
godoc -http=:6060`

## Contributing

Contributions are welcome! To contribute please:

1. Fork the repository
2. Create a feature branch
3. Code
4. Submit a [pull request][]

### Testing

Prior to submitting a [pull request][], please run:

```bash
$ make check
$ make lint
```

To update and view the test coverage report:

```bash
$ make cover
```

## License

[ivi-examples][] is released under the MIT license. Please see the
[LICENSE.txt][] file for more information.

[ivi]: https://github.com/gotmc/ivi
[ivi-examples]: https://github.com/gotmc/ivi-examples
[ivi-foundation]: http://www.ivifoundation.org/
[ivi-specs]: http://www.ivifoundation.org/specifications/
[godoc badge]: https://godoc.org/github.com/gotmc/ivi-examples?status.svg
[godoc link]: https://godoc.org/github.com/gotmc/ivi-examples
[LICENSE.txt]: https://github.com/gotmc/ivi-examples/blob/master/LICENSE.txt
[license badge]: https://img.shields.io/badge/license-MIT-blue.svg
[lxi]: https://github.com/gotmc/lxi
[prologix]: https://github.com/gotmc/prologix
[pull request]: https://help.github.com/articles/using-pull-requests
[report badge]: https://goreportcard.com/badge/github.com/gotmc/ivi-examples
[report card]: https://goreportcard.com/report/github.com/gotmc/ivi-examples
[scpi]: http://www.ivifoundation.org/scpi/
[usbtmc]: https://github.com/gotmc/usbtmc
[visa]: https://github.com/gotmc/visa
