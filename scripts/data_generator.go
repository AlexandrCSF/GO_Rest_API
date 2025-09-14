package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"wb_cource/internal/app/model"

	"github.com/brianvoe/gofakeit/v6"
)

func main() {
	gofakeit.Seed(time.Now().UnixNano())

	orders := generateOrders(10)

	for i, order := range orders {
		fmt.Printf("Отправка заказа %d: %s\n", i+1, order.OrderUID)

		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Printf("Ошибка маршалинга заказа %s: %v", order.OrderUID, err)
			continue
		}

		resp, err := http.Post("http://localhost:8080/order", "application/json", bytes.NewBuffer(orderJSON))
		if err != nil {
			log.Printf("Ошибка отправки заказа %s: %v", order.OrderUID, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			fmt.Printf("✅ Заказ %s успешно создан\n", order.OrderUID)
		} else {
			fmt.Printf("❌ Ошибка создания заказа %s: статус %d\n", order.OrderUID, resp.StatusCode)
		}

		time.Sleep(1 * time.Second)
	}

	fmt.Println("\n🎯 Генерация данных завершена!")
}

func generateOrders(count int) []*model.Order {
	orders := make([]*model.Order, count)

	for i := 0; i < count; i++ {
		orders[i] = generateOrder()
	}

	return orders
}

func generateOrder() *model.Order {
	orderUID := gofakeit.UUID()

	return &model.Order{
		OrderUID:    orderUID,
		TrackNumber: gofakeit.Regex("^[A-Z]{2}[0-9]{9}[A-Z]{2}$"),
		Entry:       gofakeit.RandomString([]string{"WBIL", "WBILM", "WBILZ"}),
		Delivery: model.Delivery{
			Name:    gofakeit.Name(),
			Phone:   gofakeit.Phone(),
			Zip:     gofakeit.Zip(),
			City:    gofakeit.City(),
			Address: gofakeit.Street() + ", " + gofakeit.Address().Street,
			Region:  gofakeit.State(),
			Email:   gofakeit.Email(),
		},
		Payment: model.Payment{
			Transaction:  gofakeit.UUID(),
			RequestID:    gofakeit.UUID(),
			Currency:     gofakeit.RandomString([]string{"RUB", "USD", "EUR"}),
			Provider:     gofakeit.RandomString([]string{"wbpay", "sberpay", "yoomoney"}),
			Amount:       gofakeit.IntRange(100, 10000),
			PaymentDt:    gofakeit.DateRange(time.Now().AddDate(0, -1, 0), time.Now()).Unix(),
			Bank:         gofakeit.RandomString([]string{"alpha", "sber", "tinkoff"}),
			DeliveryCost: gofakeit.IntRange(0, 500),
			GoodsTotal:   gofakeit.IntRange(100, 5000),
			CustomFee:    gofakeit.IntRange(0, 100),
		},
		Items:             generateItems(),
		Locale:            gofakeit.RandomString([]string{"ru", "en", "de"}),
		InternalSignature: "",
		CustomerID:        gofakeit.UUID(),
		DeliveryService:   gofakeit.RandomString([]string{"meest", "novaposhta", "ukrposhta"}),
		ShardKey:          gofakeit.DigitN(1),
		SmID:              gofakeit.IntRange(1, 100),
		DateCreated:       gofakeit.DateRange(time.Now().AddDate(0, -1, 0), time.Now()),
		OofShard:          gofakeit.DigitN(1),
	}
}

func generateItems() []model.Item {
	itemCount := gofakeit.IntRange(1, 5)
	items := make([]model.Item, itemCount)

	for i := 0; i < itemCount; i++ {
		price := gofakeit.IntRange(100, 2000)
		sale := gofakeit.IntRange(0, 50)
		totalPrice := price - (price * sale / 100)

		items[i] = model.Item{
			ChrtID:      gofakeit.IntRange(100000, 999999),
			TrackNumber: gofakeit.Regex("^[A-Z]{2}[0-9]{9}[A-Z]{2}$"),
			Price:       price,
			Rid:         gofakeit.UUID(),
			Name:        gofakeit.ProductName(),
			Sale:        sale,
			Size:        gofakeit.RandomString([]string{"XS", "S", "M", "L", "XL", "XXL"}),
			TotalPrice:  totalPrice,
			NmID:        gofakeit.IntRange(100000, 999999),
			Brand:       gofakeit.Company(),
			Status:      gofakeit.IntRange(200, 210),
		}
	}

	return items
}
