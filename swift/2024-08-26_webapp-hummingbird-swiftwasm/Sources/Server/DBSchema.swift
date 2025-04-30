import MySQLKit
import Schema

extension City {
    public init?(row: MySQLRow) {
        guard let cityID = row.column("ID")?.int else {
            return nil
        }

        self.init(
            id: cityID,
            name: row.column("Name")?.string,
            countryCode: row.column("CountryCode")?.string,
            district: row.column("District")?.string,
            population: row.column("Population")?.string
        )
    }
}
