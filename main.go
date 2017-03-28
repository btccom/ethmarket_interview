package main

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"github.com/hoisie/web"
	"github.com/btccom/ethmarket_interview/orders"
	"github.com/libreoscar/utils/log"
)

var (
	logger = log.New(log.DEBUG)
)

func init() {
	rand.Seed(1234)
}

func randSide() string {
	r := rand.Intn(2)
	if r == 0 {
		return "ask"
	}
	return "bid"
}

func generateOrders(pump *orders.OrderPump) {
	for {
		durMs := 200 + 2*rand.Intn(1000)
		time.Sleep(time.Duration(durMs) * time.Millisecond)

		order := &orders.Order{
			Side:     randSide(),
			Quantity: 1 + rand.Intn(15),
			Price:    50.0 + float64(rand.Intn(180)),
		}
		pump.AddOrder(order)
	}
}

func listOrders(ctx *web.Context, pump *orders.OrderPump) {
	start, err := strconv.Atoi(ctx.Params["start"])
	if err != nil || start < 0 {
		ctx.ResponseWriter.WriteHeader(400)
		ctx.ResponseWriter.Write([]byte(`Invalid values for param "start"`))
		return
	}

	size, err := strconv.Atoi(ctx.Params["size"])
	if err != nil || size < 0 {
		ctx.ResponseWriter.WriteHeader(400)
		ctx.ResponseWriter.Write([]byte(`Invalid values for param "size"`))
		return
	}

	orders := pump.List(start, size)

	ctx.ContentType("json")
	ctx.SetHeader("Access-Control-Allow-Origin", "*", true)
	err = json.NewEncoder(ctx).Encode(orders)
	if err != nil {
		panic(err)
	}
}

func main() {
	pump := orders.NewOrderPump()

	go generateOrders(pump)

	web.Get("/reset", func(ctx *web.Context) {
		ctx.SetHeader("Access-Control-Allow-Origin", "*", true)
		logger.Debug("Handling Get /reset")
		pump.Clear()
	})

	web.Get("/listOrders", func(ctx *web.Context) {
		logger.Debug("Handling Get /listOrders")
		listOrders(ctx, pump)
	})

	web.Run("0.0.0.0:9888")
}
