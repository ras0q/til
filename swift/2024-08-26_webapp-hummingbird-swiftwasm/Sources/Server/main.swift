import Foundation
import Hummingbird
import MySQLKit

let env = ProcessInfo.processInfo.environment

var tlsConfiguration: TLSConfiguration = .makeClientConfiguration()
tlsConfiguration.certificateVerification = .none
let mysqlConfiguration = MySQLConfiguration(
    hostname: env["MYSQL_HOST"] ?? "localhost",
    port: env["MYSQL_PORT"].flatMap(Int.init) ?? 3306,
    username: env["MYSQL_USER"] ?? "root",
    password: env["MYSQL_PASSWORD"] ?? "password",
    database: env["MYSQL_DATABASE"] ?? "app",
    tlsConfiguration: tlsConfiguration
)
let mysqlPool = EventLoopGroupConnectionPool(
    source: MySQLConnectionSource(configuration: mysqlConfiguration),
    maxConnectionsPerEventLoop: 1,
    on: MultiThreadedEventLoopGroup(numberOfThreads: 1)
)
defer { mysqlPool.shutdown() }

var logger = Logger(label: "mysql")
logger.logLevel = .debug
let database = mysqlPool.database(logger: logger)
let router = Router()
var handler = Handler(database: database)

router.add(middleware: CORSMiddleware(allowOrigin: .all))

router.get("/cities/:cityName", use: handler.getCity)
router.post("/cities", use: handler.postCity)

let app = Application(
    router: router,
    configuration: .init(address: .hostname("0.0.0.0", port: 8080))
)

try await app.runService()
