import Foundation
import JavaScriptKit
import JavaScriptEventLoop
import Schema
import TokamakShim

@main
struct TokamakApp: App {
    init() {
        JavaScriptEventLoop.installGlobalExecutor()
    }

    var body: some Scene {
        WindowGroup("Tokamak App") {
            ContentView()
        }
    }
}

struct ContentView: View {
    @State private var cityName: String = ""
    @State private var city: City?

    var body: some View {
        VStack {
            TextField("Enter a city name...", text: $cityName)
                .padding()
                .textFieldStyle(RoundedBorderTextFieldStyle())

            Button("Fetch the city") {
                Task {
                    do {
                        let response = try await fetch("http://localhost:8080/cities/\(cityName)").value
                        let json = try await JSPromise(response.json().object!)!.value
                        city = try JSValueDecoder().decode(City.self, from: json)
                    } catch {
                        print(error)
                    }
                }
            }

            if let city {
                Text("City: \(city)")
            }
        }
    }

    private let jsFetch = JSObject.global.fetch.function!
    private func fetch(_ url: String) -> JSPromise {
        JSPromise(jsFetch(url).object!)!
    }
}
