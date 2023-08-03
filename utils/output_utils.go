package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func ToJSON(data interface{}, indent string) (string, error) {
	var jsonbyte []byte
	var err error
	if indent == "" {
		jsonbyte, err = json.Marshal(data)
	} else {
		jsonbyte, err = json.MarshalIndent(data, "", indent)
	}
	return string(jsonbyte), err
}

func WriteFileLine(filename string, data string) error {
	if filename == "-" {
		fmt.Println(data)
	} else {
		err := os.WriteFile(filename, []byte(data), 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

func OutputJSON(data interface{}, filename string, indent string) error {
	str, err := ToJSON(data, indent)
	if err != nil {
		fmt.Printf("Error encoding JSON: %s", err)
		return err
	}
	err = WriteFileLine(filename, str)
	if err != nil {
		fmt.Printf("Error writing to file: %s", err)
		return err
	}
	return nil
}
