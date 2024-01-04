package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	// "github.com/frankw2/libfss/tree/master/go/libfss"
	// "github.com/frankw2/go/libfss"
	// "github.com/frankw2/libfss"
	// "github.com/frankw2/libfss/go"
	"github.com/cathieyun/libfss/libfss"


	// "./libfss/client"
	// "./libfss/server"

)

// Item represents the structure of your JSON objects
type Item struct {
	Category string `json:"categories"`
	// Add other fields as needed
}

func main() {
	// Specify the JSON file path
	jsonFilePath := "../yelp_dataset 2/yelp_academic_dataset_business.json"

	// Read the JSON file
	jsonData, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	// Unmarshal JSON data into a slice of maps
	var rawItems []map[string]interface{}
	err = json.Unmarshal(jsonData, &rawItems)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// Create an array to hold the extracted data
	var items []Item

	// Iterate through each map and extract the desired values
	for _, rawItem := range rawItems {
		category, categoryExists := rawItem["categories"].(string)

		// If both "name" and "category" values exist, create a struct and add to the array
		if categoryExists {
			item := Item{Category: category}
			items = append(items, item)
		}
		// Add other fields as needed

		// You can include additional logic to handle other fields if necessary
	}
	csvFilePath := "output.csv"

	csvFile, err := os.Create(csvFilePath)
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer csvFile.Close()

	// Create a CSV writer
	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	// Write header to CSV file
	header := []string{"ID", "category"}
	err = csvWriter.Write(header)
	if err != nil {
		fmt.Println("Error writing CSV header:", err)
		return
	}

	// Write data to CSV file with an incremental ID
	for i, item := range items {
		row := []string{fmt.Sprint(i + 1), item.Category}
		err := csvWriter.Write(row)
		if err != nil {
			fmt.Println("Error writing CSV row:", err)
			return
		}
	}

	// Print or process the extracted data
	// for _, item := range items {
	// 	fmt.Printf("Name: %s, Category: %s\n", item.Category)
	// 	// Add other fields as needed
	// }
	fClient := libfss.ClientInitialize(11)
	// Test with if x = 10, evaluate to 2
	fssKeys := fClient.GenerateTreePF(10, 2)

	// Simulate server
	fServer := libfss.ServerInitialize(fClient.PrfKeys, fClient.NumBits)

	// Test 2-party Equality Function
	var ans0, ans1 int = 0, 0
	ans0 = fServer.EvaluatePF(0, fssKeys[0], 10)
	ans1 = fServer.EvaluatePF(1, fssKeys[1], 10)
	fmt.Println("Match (should be non-zero):", ans0+ans1)

	ans0 = fServer.EvaluatePF(0, fssKeys[0], 11)
	ans1 = fServer.EvaluatePF(1, fssKeys[1], 11)
	fmt.Println("No Match (should be 0):", ans0+ans1)

	ans0 = fServer.EvaluatePF(0, fssKeys[0], 9)
	ans1 = fServer.EvaluatePF(1, fssKeys[1], 9)
	fmt.Println("No Match (should be 0):", ans0+ans1)

	// Test 2-party Less than Function
	// Test if x < 10, evaluate to 2
	fssKeysLt := fClient.GenerateTreeLt(10, 2)

	var anslt0, anslt1 uint = 0, 0
	anslt0 = fServer.EvaluateLt(fssKeysLt[0], 8)
	anslt1 = fServer.EvaluateLt(fssKeysLt[1], 8)
	fmt.Println("Less than (should be non-zero):", anslt0-anslt1)
	anslt0 = fServer.EvaluateLt(fssKeysLt[0], 11)
	anslt1 = fServer.EvaluateLt(fssKeysLt[1], 11)
	fmt.Println("Greater than (should be zero):", anslt0-anslt1)

	// Test multiparty equal function case
	fssKeysEqMP := fServer.GenerateTreeEqMP(10, 2, 3)

	var ansEqMP0, ansEqMP1, ansEqMP2 uint32 = 0, 0, 0
	ansEqMP0 = fServer.EvaluateEqMP(fssKeysEqMP[0], 10)
	ansEqMP1 = fServer.EvaluateEqMP(fssKeysEqMP[1], 10)
	ansEqMP2 = fServer.EvaluateEqMP(fssKeysEqMP[2], 10)

	fmt.Println("Multi-party Equal Match (should be non-zero):", ansEqMP0^ansEqMP1^ansEqMP2)

	ansEqMP0 = fServer.EvaluateEqMP(fssKeysEqMP[0], 12)
	ansEqMP1 = fServer.EvaluateEqMP(fssKeysEqMP[1], 12)
	ansEqMP2 = fServer.EvaluateEqMP(fssKeysEqMP[2], 12)
	fmt.Println("Multi-party Equal not Match (should be zero):", ansEqMP0^ansEqMP1^ansEqMP2)

}
