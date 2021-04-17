package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/adshao/go-binance/v2"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func loadConfig() error {
	if err := godotenv.Load(); err != nil {
		log.WithError(err).Fatal("Error loading .env file")
		return err
	}

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/binance.prices/")
	viper.AddConfigPath("$HOME/.binance.prices")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.WithError(err).Error("Error loading config file")
		return err
	}
	return nil
}

func main() {
	if err := loadConfig(); err != nil {
		return
	}

	apiKey := os.Getenv("BINANCE_API_KEY")
	secretKey := os.Getenv("BINANCE_SECRET_KEY")

	client := binance.NewClient(apiKey, secretKey)
	aps := client.NewAveragePriceService()

	targets := viper.GetStringSlice("targets")
	fmt.Printf("Symbol, %s\n", strings.Join(targets, ", "))

	symbols := viper.GetStringSlice("symbols")
	for _, s := range symbols {
		prices := make([]string, len(targets))
		for i, t := range targets {
			symbol := s + t
			if s == t {
				prices[i] = "1"
				continue
			}
			price, err := aps.Symbol(symbol).Do(context.Background())
			if err != nil {
				log.WithField("symbol", symbol).WithError(err).Error("failed to get price")
				continue
			}
			if price != nil {
				prices[i] = price.Price
			}
		}
		fmt.Printf("%s, %s\n", s, strings.Join(prices, ", "))
	}
}
