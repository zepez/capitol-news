package main

import (
	"bytes"
	"net/http"
	"os"
)

func SendData(data *bytes.Buffer) {
	// create new reponse, get req and err
	req, err := http.NewRequest("POST", os.Getenv("endpoint"), data)

	// set headers here
	req.Header.Set("Content-Type", "application/json")

	// error handling
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// uncomment below for debug
	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))
}
