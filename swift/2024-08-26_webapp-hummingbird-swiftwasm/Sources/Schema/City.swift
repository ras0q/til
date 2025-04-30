public struct City: Codable {
    public let id: Int
    public let name: String?
    public let countryCode: String?
    public let district: String?
    public let population: String?

    public init(id: Int = 0, name: String? = nil, countryCode: String? = nil, district: String? = nil, population: String? = nil) {
        self.id = id
        self.name = name
        self.countryCode = countryCode
        self.district = district
        self.population = population
    }
}
