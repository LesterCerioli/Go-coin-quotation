package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Quote struct {
	USD float64 `json:"usd"`
	EUR float64 `json:"eur"`
	ETH float64 `json:"eth"`
}

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

const logFile = "logs/log.json"

func fetchQuote(url string) (float64, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var result map[string]float64
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	return result["price"], nil
}

func writeLog(status, message string) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	logEntry := LogEntry{
		Timestamp: currentTime,
		Status:    status,
		Message:   message,
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer file.Close()

	logData, _ := json.Marshal(logEntry)
	file.Write(logData)
	file.Write([]byte("\n"))
}

func main() {

	usdAPI := "https://api.example.com/usd"
	eurAPI := "https://api.example.com/eur"
	ethAPI := "https://api.example.com/eth"

	for {
		usd, err := fetchQuote(usdAPI)
		if err != nil {
			writeLog("Failure", fmt.Sprintf("Error fetching USD quote: %v", err))
		} else {
			writeLog("Success", fmt.Sprintf("USD quote retrieved: %.2f", usd))
		}

		eur, err := fetchQuote(eurAPI)
		if err != nil {
			writeLog("Failure", fmt.Sprintf("Error fetching EUR quote: %v", err))
		} else {
			writeLog("Success", fmt.Sprintf("EUR quote retrieved: %.2f", eur))
		}

		eth, err := fetchQuote(ethAPI)
		if err != nil {
			writeLog("Failure", fmt.Sprintf("Error fetching ETH quote: %v", err))
		} else {
			writeLog("Success", fmt.Sprintf("ETH quote retrieved: %.2f", eth))
		}

		if usd > eur {
			writeLog("Info", "USD is greater than EUR")
		} else {
			writeLog("Info", "EUR is greater than USD")
		}

		if eth > usd {
			writeLog("Info", "ETH is greater than USD")
		} else {
			writeLog("Info", "USD is greater than ETH")
		}

		time.Sleep(2 * time.Minute)
	}
}
