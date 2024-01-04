package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/cathieyun/libfss/libfss"
)

func main() {
	
	startTime := time.Now()


	var size uint = 11 // number of bits in input domain
	var lookup uint = 1993 //looking for ...
	client := libfss.ClientInitialize(size)
	fssKeys := client.GenerateTreePF(lookup, 1)

	endTime := time.Now()

	//servers  configurations
	server1ToWrite := []string{packageKeys(fssKeys[0]), packageKeys(client.PrfKeys),strconv.FormatUint(uint64(client.NumBits), 10)}
	server2ToWrite := []string{packageKeys(fssKeys[1]), packageKeys(client.PrfKeys),strconv.FormatUint(uint64(client.NumBits), 10)}
	elapsedTime := endTime.Sub(startTime).Nanoseconds()
	print("client time in milisecond:", elapsedTime/1000)

	err  := writeCSV("../config_server/server1_config_keys.csv", server1ToWrite)
	err1 := writeCSV("../config_server/server2_config_keys.csv", server2ToWrite)

	if err != nil || err1 != nil {
		fmt.Println("Error writing to CSV:", err)
	}
	// print(t)
}

func stringToInt(s string) uint {
	h := sha256.New()
	h.Write([]byte(s))
	num := binary.LittleEndian.Uint32(h.Sum(nil))
	return uint(num)
}

func packageKeys(key interface{}) string {
	marshalledKey, _ := json.Marshal(key)
	return base64.StdEncoding.EncodeToString(marshalledKey)
}

func writeCSV(filename string, data []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, num_ := range data {
		row := []string{num_}
		err := writer.Write(row)
		if err != nil {
			panic(err)
		}
	}

	return nil
}
