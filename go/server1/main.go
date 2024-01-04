package main

import (
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	// "time"

	"github.com/cathieyun/libfss/libfss"
)

func main() {

	//read csv file and get fsskey and prfkey and numBit

	file, err := os.Open("../config_server/server1_config_keys.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)
	record1, err := reader.Read()
	if err != nil {
		fmt.Println("Error reading first line:", err)
		return
	}

	record2, err := reader.Read()
	if err != nil {
		fmt.Println("Error reading second line:", err)
		return
	}

	getdata(record1, record2)


}
func getdata(record1, record2 []string) {


	fssKey := record1[0]
	prfKeys := record2[0]
	// print("server 2 fss key",fssKey)

	ans := evalQuery(fssKey, prfKeys)
	print("Answer from server 2 is : ", ans)
	// print(fssKey, "\n", prfKeys)
}


func evalQuery(fssKey, prfKeys string) string {
	var parsedPrfKeys [][]byte
	_ = json.Unmarshal(decodeKey(prfKeys), &parsedPrfKeys)
	Server := libfss.ServerInitialize(parsedPrfKeys, uint(11))
	var parsedFssKey libfss.FssKeyEq2P
	_ = json.Unmarshal(decodeKey(fssKey), &parsedFssKey)

	ans := ""
	// Evaluate PF over all values of DB
	serverNum := byte(2 - 1)
	ans = readIntFetchSmall(Server, serverNum, parsedFssKey, "")
	return ans
}



func decodeKey(str string) []byte {
	dec, _ := base64.StdEncoding.DecodeString(str)
	return dec
}
