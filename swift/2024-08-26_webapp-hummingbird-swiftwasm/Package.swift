// swift-tools-version: 5.9
// The swift-tools-version declares the minimum version of Swift required to build this package.

import PackageDescription

let package = Package(
    name: "SwiftWebAppPlayground",
    dependencies: [
        .package(url: "https://github.com/hummingbird-project/hummingbird", from: "2.0.0-rc.4"),
        .package(url: "https://github.com/swiftwasm/carton", from: "1.0.0"),
        .package(url: "https://github.com/swiftwasm/JavaScriptKit.git", from: "0.20.1"), // follow the Tokamak's JavaScriptKit version
        .package(url: "https://github.com/TokamakUI/Tokamak", from: "0.11.0"),
        .package(url: "https://github.com/vapor/mysql-kit.git", from: "4.0.0"),
    ],
    targets: [
        .executableTarget(name: "Client", dependencies: [
            .target(name: "Schema"),
            .product(name: "JavaScriptKit", package: "JavaScriptKit"),
            .product(name: "JavaScriptEventLoop", package: "JavaScriptKit"),
            .product(name: "TokamakShim", package: "Tokamak")
        ]),
        .target(name: "Schema", dependencies: []),
        .executableTarget(name: "Server", dependencies: [
            .target(name: "Schema"),
            .product(name: "Hummingbird", package: "hummingbird"),
            .product(name: "MySQLKit", package: "mysql-kit")
        ])
    ]
)
