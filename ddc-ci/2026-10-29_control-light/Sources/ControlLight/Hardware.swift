import CoreGraphics
import Darwin
import Foundation

final class BuiltinDisplay {
  private typealias GetBrightness = @convention(c) (UInt32, UnsafeMutablePointer<Float>) -> Int32
  private let framework: UnsafeMutableRawPointer
  private let getBrightness: GetBrightness

  init() throws {
    let path = "/System/Library/PrivateFrameworks/DisplayServices.framework/DisplayServices"
    guard let handle = dlopen(path, RTLD_LAZY) else {
      throw ControlError("Cannot load DisplayServices")
    }
    guard let symbol = dlsym(handle, "DisplayServicesGetBrightness") else {
      dlclose(handle)
      throw ControlError("DisplayServicesGetBrightness is unavailable")
    }
    framework = handle
    getBrightness = unsafeBitCast(symbol, to: GetBrightness.self)
  }

  deinit { dlclose(framework) }

  func read() throws -> Brightness {
    var count: UInt32 = 0
    var status = CGGetOnlineDisplayList(0, nil, &count)
    guard status == .success else { throw ControlError("Cannot enumerate displays: \(status)") }
    var displays = [CGDirectDisplayID](repeating: 0, count: Int(count))
    status = CGGetOnlineDisplayList(count, &displays, &count)
    guard status == .success else { throw ControlError("Cannot enumerate displays: \(status)") }
    guard let display = displays.prefix(Int(count)).first(where: { CGDisplayIsBuiltin($0) != 0 })
    else {
      throw ControlError("No online built-in display. Open the MacBook lid.")
    }
    var value: Float = 0
    let result = getBrightness(display, &value)
    guard result == 0 else { throw ControlError("Cannot read built-in brightness: \(result)") }
    return try Brightness(Double(value))
  }
}

struct M1DDC {
  private let executable: URL

  init() throws {
    let directories = (ProcessInfo.processInfo.environment["PATH"] ?? "").split(separator: ":")
    guard
      let path = directories.map({
        URL(fileURLWithPath: String($0)).appendingPathComponent("m1ddc")
      })
      .first(where: { FileManager.default.isExecutableFile(atPath: $0.path) })
    else { throw ControlError("m1ddc is missing from PATH. See README.md for setup.") }
    executable = path
  }

  func run(_ arguments: [String]) throws -> String {
    let process = Process()
    let output = Pipe()
    process.executableURL = executable
    process.arguments = arguments
    process.standardOutput = output
    process.standardError = FileHandle.standardError
    try process.run()
    let data = output.fileHandleForReading.readDataToEndOfFile()
    process.waitUntilExit()
    guard process.terminationReason == .exit, process.terminationStatus == 0 else {
      throw ControlError(
        "m1ddc \(arguments.joined(separator: " ")) failed (\(process.terminationStatus))")
    }
    guard let text = String(data: data, encoding: .utf8) else {
      throw ControlError("m1ddc returned invalid UTF-8")
    }
    return text
  }
}
