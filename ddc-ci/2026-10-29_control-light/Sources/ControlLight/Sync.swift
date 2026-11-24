import Foundation

struct ControlError: Error, CustomStringConvertible {
  let description: String
  init(_ description: String) { self.description = description }
}

struct Options {
  var displays: [String] = []
  var interval = 1.0
  var once = false
  var help = false

  init(_ arguments: [String]) throws {
    var args = arguments.makeIterator()
    while let arg = args.next() {
      switch arg {
      case "--display":
        guard let value = args.next(), !value.isEmpty, !value.hasPrefix("-") else {
          throw ControlError("--display requires an m1ddc display number or identifier")
        }
        if !displays.contains(value) { displays.append(value) }
      case "--interval":
        guard let value = args.next(), let seconds = Double(value),
          seconds.isFinite, seconds >= 0.2
        else { throw ControlError("--interval must be at least 0.2 seconds") }
        interval = seconds
      case "--once": once = true
      case "--help", "-h": help = true
      default: throw ControlError("Unknown argument: \(arg)")
      }
    }
    guard help || !displays.isEmpty else {
      throw ControlError(
        "Specify an external display with --display ID. Use m1ddc display list to find it.")
    }
  }
}

struct Brightness {
  let fraction: Double

  init(_ fraction: Double) throws {
    guard fraction.isFinite, (0...1).contains(fraction) else {
      throw ControlError("Invalid built-in brightness: \(fraction)")
    }
    self.fraction = fraction
  }

  func target(maximum: Int) throws -> Int {
    guard (1...65535).contains(maximum) else {
      throw ControlError("Invalid DDC luminance maximum: \(maximum)")
    }
    return Int((fraction * Double(maximum)).rounded())
  }
}

func synchronize(
  _ brightness: Brightness,
  display: String,
  command: ([String]) throws -> String
) throws -> Int? {
  func read(_ operation: String) throws -> Int {
    let output = try command(["display", display, operation, "luminance"])
    guard let value = Int(output.trimmingCharacters(in: .whitespacesAndNewlines)) else {
      throw ControlError("Display \(display): invalid \(operation) luminance response: \(output)")
    }
    return value
  }
  let maximum = try read("max")
  let target = try brightness.target(maximum: maximum)
  let current = try read("get")
  guard (0...maximum).contains(current) else {
    throw ControlError("Display \(display): invalid current luminance: \(current)")
  }
  guard current != target else { return nil }
  _ = try command(["display", display, "set", "luminance", String(target)])
  return target
}
