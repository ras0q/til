// swift-tools-version: 6.0
import PackageDescription

let package = Package(
  name: "control-light",
  platforms: [.macOS(.v13)],
  targets: [
    .executableTarget(name: "ControlLight"),
    .testTarget(name: "ControlLightTests", dependencies: ["ControlLight"]),
  ]
)
