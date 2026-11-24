import Testing

@testable import ControlLight

struct SyncTests {
  @Test func scalesToMonitorMaximumAndAvoidsRedundantWrites() throws {
    var calls: [[String]] = []
    let brightness = try Brightness(0.425)
    let result = try synchronize(brightness, display: "UUID") { args in
      calls.append(args)
      return args[2] == "max" ? "200\n" : "85\n"
    }
    #expect(result == nil)
    #expect(calls.count == 2)

    calls = []
    let updated = try synchronize(brightness, display: "UUID") { args in
      calls.append(args)
      return args[2] == "max" ? "200\n" : "10\n"
    }
    #expect(updated == 85)
    #expect(calls.last == ["display", "UUID", "set", "luminance", "85"])
  }

  @Test func rejectsInvalidReadingsWithoutWriting() throws {
    for response in ["-1", "0", "65536", "error", ""] {
      var writes = 0
      #expect(throws: ControlError.self) {
        try synchronize(Brightness(0.5), display: "1") { args in
          if args[2] == "set" { writes += 1 }
          return response
        }
      }
      #expect(writes == 0)
    }
    #expect(throws: ControlError.self) { try Brightness(.nan) }
    #expect(throws: ControlError.self) { try Brightness(1.1) }
    #expect(try Brightness(0).target(maximum: 100) == 0)
    #expect(try Brightness(1).target(maximum: 255) == 255)
  }

  @Test func parsesOptions() throws {
    let options = try Options([
      "--display", "UUID", "--display", "2", "--once", "--interval", "0.5",
    ])
    #expect(options.displays == ["UUID", "2"])
    #expect(options.once)
    #expect(options.interval == 0.5)
    #expect(try Options(["--help"]).help)
    for args in [[], ["--interval", "nan"], ["--interval", "0"], ["--display"], ["--typo"]] {
      #expect(throws: ControlError.self) { try Options(args) }
    }
  }

  @Test func rejectsInvalidCurrentValuesAndPropagatesWriteFailures() throws {
    for response in ["-1", "101", "unsupported"] {
      #expect(throws: ControlError.self) {
        try synchronize(Brightness(0.5), display: "1") { args in
          #expect(args[2] != "set")
          return args[2] == "max" ? "100" : response
        }
      }
    }
    #expect(throws: ControlError.self) {
      try synchronize(Brightness(0.5), display: "1") { args in
        switch args[2] {
        case "max": return "100"
        case "get": return "0"
        default: throw ControlError("DDC write failed")
        }
      }
    }
  }
}
