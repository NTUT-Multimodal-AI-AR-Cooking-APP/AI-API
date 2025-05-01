import UIKit



class ViewController: UIViewController {



    @IBAction func request(_ sender: Any) {

        // 設定設備、食材和偏好

        let equipment = [Equipment(name: "平底鍋", type: "鍋具", size: "中型", material: "不鏽鋼")]

        let ingredients = [

            Ingredient(name: "油", type: "食材", amount: "2", unit: "湯匙", weight: nil),

            Ingredient(name: "牛排", type: "食材", amount: "1", unit: "塊", weight: "250g")

        ]

        let preference = Preference(cooking_method: "煎", doneness: "中等熟")

        

        // 創建請求資料

        let request = RecipeRequest(equipment: equipment, ingredients: ingredients, preference: preference)



        // 呼叫 generateRecipe 函數發送請求

        generateRecipe(request: request)

    }



    override func viewDidLoad() {

        super.viewDidLoad()

    }



    // 發送 POST 請求的函數

    func generateRecipe(request: RecipeRequest) {

        // 嘗試創建 URL "http://<你的Windows伺服器IP>:8080/generate-recipe"

        guard let url = URL(string: "http://192.168.0.103:8080/generate-recipe") else {

            print("無效的 URL")

            return

        }



        // 創建 URLRequest 物件

        var urlRequest = URLRequest(url: url)

        urlRequest.httpMethod = "POST"

        urlRequest.setValue("application/json", forHTTPHeaderField: "Content-Type")

        

        // 將資料編碼為 JSON

        let encoder = JSONEncoder()

        do {

            let jsonData = try encoder.encode(request) // 編碼 RecipeRequest

            urlRequest.httpBody = jsonData

        } catch {

            print("無法編碼請求資料：\(error)")

            return

        }

        

        // 使用 URLSession 發送請求

        let task = URLSession.shared.dataTask(with: urlRequest) { data, response, error in

            // 偵錯：檢查錯誤

            if let error = error as NSError? {

                print("錯誤：\(error.localizedDescription)")

                print("錯誤碼：\(error.code)") // 輸出錯誤碼

                return

            }



            // 偵錯：檢查回應

            if let httpResponse = response as? HTTPURLResponse {

                print("回應狀態碼：\(httpResponse.statusCode)")

                // 這裡可以根據 statusCode 判斷回應是否成功

                if httpResponse.statusCode != 200 {

                    print("伺服器回應錯誤，狀態碼：\(httpResponse.statusCode)")

                }

            }



            // 偵錯：檢查資料

            guard let data = data else {

                print("沒有收到回應資料")

                return

            }



            // 嘗試解析 JSON 回應

            do {

                let decoder = JSONDecoder()

                let recipeResponse = try decoder.decode(RecipeResponse.self, from: data)

                print("食譜名稱：\(recipeResponse.dish_name)")

                print("食譜描述：\(recipeResponse.dish_description)")

                

                for step in recipeResponse.recipe {

                    print("步驟：\(step.step), 時間：\(step.time), 溫度：\(step.temperature), 描述：\(step.description), 熟度：\(step.doneness ?? "無")")

                }

            } catch {

                print("解析 JSON 失敗：\(error)")

            }

        }

        

        // 偵錯：請求開始前

        print("發送請求：\(urlRequest)")

        

        // 啟動請求

        task.resume()

    }





}



// 設備結構

struct Equipment: Codable {

    var name: String

    var type: String

    var size: String

    var material: String

}



// 食材結構

struct Ingredient: Codable {

    var name: String

    var type: String

    var amount: String

    var unit: String

    var weight: String? // 這是選擇性的，並非所有食材都有 weight

}



// 偏好設定結構

struct Preference: Codable {

    var cooking_method: String

    var doneness: String

}



// 食譜請求資料結構

struct RecipeRequest: Codable {

    var equipment: [Equipment]

    var ingredients: [Ingredient]

    var preference: Preference

}



// 食譜回應資料結構

struct RecipeResponse: Codable {

    var dish_name: String

    var dish_description: String

    var recipe: [RecipeStep]

}



// 食譜步驟資料結構

struct RecipeStep: Codable {

    var step: String

    var time: String

    var temperature: String

    var description: String

    var doneness: String?

}



