package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type CoinbaseResult struct {
	Data CoinbaseData `json:"data"`
}

type CoinbaseData struct {
	Currency string            `json:"currency"`
	Rates    map[string]string `json:"rates"`
}

func main() {
	var investmentAmount float64
	var coinbaseResult CoinbaseResult
	url := "https://api.coinbase.com/v2/exchange-rates?currency=USD"

	// Command line argument for how many dollars the user wants to invest
	if len(os.Args) > 1 {
		var err error
		investmentAmount, err = strconv.ParseFloat(os.Args[1], 64)
		if err != nil {
			fmt.Println("Error parsing the amount to invest as a float: ", err)
			return
		}
		if investmentAmount <= 0 {
			fmt.Println("Please provide a positive investment amount")
			return
		}
	} else {
		fmt.Println("Missing argument for amount to invest.")
		fmt.Println("Usage: go run main.go [amount]")
		fmt.Println("Example: go run main.go 1623.56")
		return
	}

	// Get data from Coinbase's API
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error occurred trying to get data from the API: ", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error happened when reading the response body: ", err)
		return
	}

	// I wanted to use the variable name "err" here, but I got a compiler error
	// "no new variables on left side of :=" (NoNewVar)
	// It's weird because I used that pattern twice above
	jsonError := json.Unmarshal(body, &coinbaseResult)
	if jsonError != nil {
		fmt.Println("Error from doing json.Unmarshal: ", jsonError)
		return
	}

	// Price of BTC
	btcRate, _ := GetRate(coinbaseResult, "BTC")
	if btcRate <= 0 {
		fmt.Println("Invalid rate for BTC")
		return
	}

	// Price of ETH
	ethRate, _ := GetRate(coinbaseResult, "ETH")
	if ethRate <= 0 {
		fmt.Println("Invalid rate for ETH")
		return
	}

	// How much BTC I get with 70% of investment dollars
	btcInvestment := 0.7 * investmentAmount * btcRate

	// How much ETH I get with 30% of investment dollars
	ethInvestment := 0.3 * investmentAmount * ethRate

	// Output JSON with eth: eth_amount, btc: btc_amount
	result, err := json.Marshal(map[string]float64{
		"btc": btcInvestment,
		"eth": ethInvestment,
	})
	if err != nil {
		fmt.Println("Error trying to marshal: ", err)
		return
	}
	fmt.Println(string(result))
}

func GetRate(coinbaseResult CoinbaseResult, crypto string) (float64, error) {
	rate := coinbaseResult.Data.Rates[crypto]
	amount, err := strconv.ParseFloat(rate, 64)
	if err != nil {
		fmt.Println("Error parsing the rate as a float: ", err)
		return 0, err
	}

	return amount, nil
}
