# control-light

[![Cute task runner: run with `cute`](https://raw.githubusercontent.com/ras0q/cute/refs/heads/main/badge.svg)](https://github.com/ras0q/cute)

A Swift CLI that synchronizes DDC/CI displays with the MacBook's built-in brightness,
even when an external display is configured as the macOS main display.

Brightness is matched as a percentage of each display's range: 50% on the MacBook
sets an external display to half its maximum DDC luminance value. Actual luminance
(nits) and perceived brightness can differ between panels.

## Requirements

- Apple Silicon MacBook, macOS 13 or later, and Swift 6 or later.
- [m1ddc](https://github.com/waydabber/m1ddc) on `PATH`, or built locally with the task below.
- An external display that supports DDC/CI brightness reads and writes, with DDC/CI enabled in its settings.

m1ddc does not support Intel Macs. DDC communication also depends on the port,
adapter, and dock used to connect the display.

## Tasks

Run commands from this directory. Install m1ddc, build the CLI, and list the displays.
Choose an external display from List Displays, check it with Sync Once, then start Run.

### Install m1ddc

Clone and build m1ddc locally. Run this once; it requires Git and the Xcode Command
Line Tools. The display tasks below use this local binary or an existing `m1ddc` on `PATH`.

```sh
mkdir -p tmp
git clone https://github.com/waydabber/m1ddc.git tmp/m1ddc
make -C tmp/m1ddc
```

### Build

```sh
swift build -c release
```

### List Displays

Use the output to find an external display's number or UUID.

```sh
PATH="$PWD/tmp/m1ddc:$PATH" m1ddc display list
```

### Sync Once

Read the built-in brightness with DisplayServices and synchronize once.
Set `DISPLAY_ID` to the target external display's number or UUID (for example, `2`).
Requires m1ddc; builds the CLI before running.

```sh
: "${DISPLAY_ID:?Set DISPLAY_ID to an external display number or UUID}"
PATH="$PWD/tmp/m1ddc:$PATH" swift run -c release ControlLight --display "$DISPLAY_ID" --once
```

### Run

Read the built-in brightness with DisplayServices and synchronize continuously,
waiting one second between cycles. Only the target external display uses DDC/CI.
Set `DISPLAY_ID` to its number or UUID (for example, `2`).
Requires m1ddc; builds the CLI before running. Stop with Ctrl-C.

```sh
: "${DISPLAY_ID:?Set DISPLAY_ID to an external display number or UUID}"
PATH="$PWD/tmp/m1ddc:$PATH" swift run -c release ControlLight --display "$DISPLAY_ID"
```

### Test

Tests cover brightness scaling, redundant-write suppression, invalid readings,
DDC command failures, and argument parsing. They do not require a connected monitor.

```sh
swift test
```

### Test With CLT Plugin

If Command Line Tools cannot find `TestingMacros`, specify the plugin directory.

```sh
swift test -Xswiftc -plugin-path -Xswiftc /Library/Developer/CommandLineTools/usr/lib/swift/host/plugins/testing
```

## Custom displays and intervals

For custom settings, run the CLI with m1ddc on `PATH`:

```console
# Synchronize two displays with a half-second wait between cycles.
PATH="$PWD/tmp/m1ddc:$PATH" .build/release/ControlLight --display 1 --display 2 --interval 0.5

# Use a UUID from List Displays to avoid changes in display numbering.
PATH="$PWD/tmp/m1ddc:$PATH" .build/release/ControlLight --display YOUR-DISPLAY-UUID
```

## Behavior and limitations

Each cycle reads the built-in brightness and each target's current and maximum
luminance, then writes only when the target differs. Changes from the MacBook's
brightness keys or automatic brightness adjustment are reflected on the next cycle.
Manual changes on an external display are overwritten while synchronization runs.
`--interval` is the wait after each cycle, excluding communication time; its minimum is 0.2 seconds.

The program reports an error and exits if the built-in display is unavailable
(for example, in clamshell mode), a reading is invalid, or a DDC command fails.
Restart it after reconnecting a display. Displays are updated sequentially;
if a later update fails, earlier updates remain applied.

Built-in brightness is read through the private `DisplayServicesGetBrightness` API,
which may change with macOS updates. Implementation references:
[brightness](https://github.com/nriley/brightness/blob/master/brightness.c) and
[m1ddc commands](https://github.com/waydabber/m1ddc#available-commands).

To check hardware behavior, run Sync Once with a DDC/CI display connected, change
the MacBook's brightness, and run it again.
