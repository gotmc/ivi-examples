# Timeout Usage Examples

This document provides examples of how to use the timeout functionality in the IVI library.

## Basic Timeout Usage

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/gotmc/ivi"
    "github.com/gotmc/ivi/dmm/keysight/key3446x"
    "github.com/gotmc/lxi"  // or your preferred connection method
)

func main() {
    // Connect to instrument (example using LXI)
    conn, err := lxi.NewInstrument("192.168.1.100")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // Method 1: Create instrument with default timeout configuration
    inst := ivi.NewWithTimeout(conn, ivi.NewDefaultTimeoutConfig())

    // Create DMM driver with timeout-enabled instrument
    dmm, err := key3446x.New(inst, false)
    if err != nil {
        log.Fatal(err)
    }

    // All operations will now timeout if they take too long
    manufacturer, err := dmm.InstrumentManufacturer()
    if err != nil {
        if err == context.DeadlineExceeded {
            log.Println("Query timed out")
        } else {
            log.Printf("Query error: %v", err)
        }
        return
    }

    log.Printf("Manufacturer: %s", manufacturer)
}
```

## Custom Timeout Configuration

```go
package main

import (
    "log"
    "time"

    "github.com/gotmc/ivi"
    "github.com/gotmc/ivi/fgen/keysight/key33220"
    "github.com/gotmc/usbtmc"  // example using USB
)

func main() {
    // Connect to instrument via USB
    conn, err := usbtmc.NewDevice("USB0::0x0957::0x0407::MY44012345::INSTR")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // Create custom timeout configuration
    config := &ivi.TimeoutConfig{
        IOTimeout:      3 * time.Second,   // Shorter timeout for I/O
        QueryTimeout:   8 * time.Second,   // Reasonable timeout for queries
        CommandTimeout: 3 * time.Second,   // Quick timeout for commands
        ResetTimeout:   20 * time.Second,  // Longer timeout for reset
        ClearTimeout:   3 * time.Second,   // Quick timeout for clear
    }

    // Wrap instrument with custom timeout configuration
    inst := ivi.NewWithTimeout(conn, config)

    // Create function generator driver
    fgen, err := key33220.New(inst, true) // true = reset on initialization
    if err != nil {
        log.Fatal(err)
    }

    // Configure a sine wave with timeout protection
    channel := fgen.Channels[0]
    err = channel.ConfigureStandardWaveform(
        key33220.WaveformSine,
        1000.0,  // 1 kHz
        2.0,     // 2V amplitude
        0.0,     // 0V offset
        0.0,     // 0 degrees phase
    )
    if err != nil {
        log.Printf("Error configuring waveform: %v", err)
        return
    }

    log.Println("Waveform configured successfully")
}
```

## Context-Based Timeout Control

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/gotmc/ivi"
    "github.com/gotmc/ivi/dmm/keysight/key3446x"
)

func main() {
    // Assume conn is your instrument connection
    var conn ivi.Instrument // = your connection

    // Create instrument with timeout support
    inst := ivi.NewWithTimeout(conn, ivi.NewDefaultTimeoutConfig())

    // Use context for fine-grained timeout control
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    // Direct context-aware calls (if you need precise control)
    if instWithTimeout, ok := interface{}(inst).(ivi.InstrumentWithTimeout); ok {
        // Perform a query with a specific 2-second timeout
        response, err := instWithTimeout.QueryWithContext(ctx, "*IDN?")
        if err != nil {
            if err == context.DeadlineExceeded {
                log.Println("Query timed out after 2 seconds")
            } else {
                log.Printf("Query error: %v", err)
            }
            return
        }
        log.Printf("Instrument ID: %s", response)

        // Send a command with the same timeout
        err = instWithTimeout.CommandWithContext(ctx, "*RST")
        if err != nil {
            if err == context.DeadlineExceeded {
                log.Println("Reset command timed out")
            } else {
                log.Printf("Command error: %v", err)
            }
            return
        }
        log.Println("Reset completed")
    }
}
```

## Setting Timeouts on Existing Drivers

```go
package main

import (
    "log"
    "time"

    "github.com/gotmc/ivi/dmm/keysight/key3446x"
)

func main() {
    // Assume you have an existing instrument connection
    var conn ivi.Instrument // = your connection

    // Create driver with existing connection
    dmm, err := key3446x.New(conn, false)
    if err != nil {
        log.Fatal(err)
    }

    // Method 1: Set a simple timeout for all operations
    dmm.SetTimeout(15 * time.Second)

    // Method 2: Set a detailed timeout configuration
    dmm.SetTimeoutConfig(&ivi.TimeoutConfig{
        IOTimeout:      5 * time.Second,
        QueryTimeout:   12 * time.Second,
        CommandTimeout: 5 * time.Second,
        ResetTimeout:   25 * time.Second,
        ClearTimeout:   4 * time.Second,
    })

    // Now all DMM operations use the configured timeouts
    voltage, err := dmm.ReadMeasurement(10 * time.Second)
    if err != nil {
        log.Printf("Measurement error: %v", err)
        return
    }

    log.Printf("Measured voltage: %.6f V", voltage)
}
```

## Error Handling with Timeouts

```go
package main

import (
    "context"
    "errors"
    "log"
    "time"

    "github.com/gotmc/ivi"
)

func handleInstrumentOperation(inst ivi.Instrument) {
    // Wrap with timeout
    timeoutInst := ivi.NewWithTimeout(inst, &ivi.TimeoutConfig{
        QueryTimeout: 5 * time.Second,
    })

    response, err := timeoutInst.Query("*IDN?")
    
    switch {
    case err == nil:
        log.Printf("Success: %s", response)
    case errors.Is(err, context.DeadlineExceeded):
        log.Println("Operation timed out - instrument may be unresponsive")
        // Could implement retry logic here
    case errors.Is(err, context.Canceled):
        log.Println("Operation was canceled")
    default:
        log.Printf("Other error: %v", err)
    }
}
```

## Retry Logic with Timeouts

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/gotmc/ivi"
)

func queryWithRetry(inst ivi.InstrumentWithTimeout, query string, maxRetries int) (string, error) {
    for attempt := 0; attempt <= maxRetries; attempt++ {
        ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
        
        result, err := inst.QueryWithContext(ctx, query)
        cancel()
        
        if err == nil {
            return result, nil
        }
        
        if err == context.DeadlineExceeded && attempt < maxRetries {
            log.Printf("Attempt %d timed out, retrying...", attempt+1)
            time.Sleep(1 * time.Second) // Brief delay before retry
            continue
        }
        
        return "", err
    }
    
    return "", context.DeadlineExceeded
}

func main() {
    // Example usage
    var conn ivi.Instrument // = your connection
    inst := ivi.NewWithTimeout(conn, ivi.NewDefaultTimeoutConfig())
    
    if instWithTimeout, ok := interface{}(inst).(ivi.InstrumentWithTimeout); ok {
        result, err := queryWithRetry(instWithTimeout, "*IDN?", 3)
        if err != nil {
            log.Printf("Failed after retries: %v", err)
        } else {
            log.Printf("Result: %s", result)
        }
    }
}
```

## Notes

- The timeout feature is backward compatible. Existing code continues to work without modification.
- To enable timeouts, wrap your instrument connection with `ivi.NewWithTimeout()` before creating drivers.
- Different timeout values can be set for different operation types (IO, Query, Command, Reset, Clear).
- Context-based methods provide the most control for specific operations.
- Always check for `context.DeadlineExceeded` errors to handle timeout scenarios appropriately.
- Consider implementing retry logic for critical operations that might occasionally timeout due to network issues.

For more complete working examples, see the [ivi-examples repository](https://github.com/gotmc/ivi-examples).