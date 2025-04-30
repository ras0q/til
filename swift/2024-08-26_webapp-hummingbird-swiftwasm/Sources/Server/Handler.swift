import Foundation
import Hummingbird
import MySQLKit
import Schema

actor Handler {
    let database: MySQLDatabase

    init(database: MySQLDatabase) {
        self.database = database
    }

    @Sendable
    func getCity(request: Request, context: BasicRequestContext) async throws -> some ResponseGenerator {
        guard let cityName = context.parameters.get("cityName") else {
            throw HTTPError(.badRequest, message: "City name not provided")
        }

        let query = database.query("SELECT * FROM city WHERE Name = ? LIMIT 1", [MySQLData(string: cityName)])
        guard let row = try await query.get().first else {
            throw HTTPError(.notFound, message: "City not found")
        }

        guard let city = City(row: row) else {
            throw HTTPError(.internalServerError, message: "Failed to decode city")
        }

        return Response(status: .ok, from: city)
    }

    struct PostCityRequest: Decodable {
        let name: String?
        let countryCode: String?
        let district: String?
        let population: String?
    }

    @Sendable
    func postCity(request: Request, context: BasicRequestContext) async throws -> some ResponseGenerator {
        let cityRequest = try await request.decode(as: City.self, context: context)

        let query = database.query("INSERT INTO city (Name, CountryCode, District, Population) VALUES (?, ?, ?, ?)", [
            MySQLData(string: cityRequest.name ?? ""),
            MySQLData(string: cityRequest.countryCode ?? ""),
            MySQLData(string: cityRequest.district ?? ""),
            MySQLData(string: cityRequest.population ?? "")
        ])
        guard let row = try await query.get().first else {
            throw HTTPError(.internalServerError, message: "Failed to insert city")
        }

        guard let city = City(row: row) else {
            throw HTTPError(.internalServerError, message: "Failed to decode city")
        }

        return Response(status: .ok, from: city)
    }
}

extension Response {
    init(status: HTTPTypes.HTTPResponse.Status, headers: HTTPFields = .init(), from value: some Encodable) {
        let json = try! JSONEncoder().encode(value)
        self.init(status: status, headers: headers, body: ResponseBody(byteBuffer: ByteBuffer(bytes: json)))
    }
}
