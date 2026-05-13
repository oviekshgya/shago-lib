package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/oviekshgya/shago-lib/excel"
)

func main() {
	result, err := excel.ExcelToJSON("./FORM - SHAGO WHATSAPP AI.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonBytes))
}
