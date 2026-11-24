import Darwin
import Foundation

@main
enum ControlLight {
  static func main() {
    do {
      let options = try Options(Array(CommandLine.arguments.dropFirst()))
      if options.help {
        print(
          """
          Usage: ControlLight --display ID [--display ID ...] [--interval SECONDS] [--once]
            --display ID        Target external display number or UUID; required, repeatable
            --interval SECONDS  Poll interval, at least 0.2 seconds (default: 1)
            --once              Synchronize once and exit
            --help, -h          Show this help
          Stop continuous synchronization with Ctrl-C.
          """)
        return
      }
      #if !arch(arm64)
        throw ControlError("This program requires an Apple Silicon MacBook (m1ddc).")
      #endif
      let ddc = try M1DDC()
      let builtin = try BuiltinDisplay()
      repeat {
        let brightness = try builtin.read()
        for display in options.displays {
          if let target = try synchronize(brightness, display: display, command: ddc.run) {
            print(
              "Built-in \(Int((brightness.fraction * 100).rounded()))% → display \(display): DDC \(target)"
            )
            fflush(stdout)
          }
        }
        if options.once { break }
        Thread.sleep(forTimeInterval: options.interval)
      } while true
    } catch {
      FileHandle.standardError.write(Data("control-light: \(error)\n".utf8))
      exit(1)
    }
  }
}
