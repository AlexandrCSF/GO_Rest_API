package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"wb_cource/internal/app/model"
)

func main() {
	// Тестовые заказы
	orders := []*model.Order{
		createTestOrder("order1", "Test User 1"),
		createTestOrder("order2", "Test User 2"),
		createTestOrder("order3", "Test User 3"),
	}

	// Отправляем заказы
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

		time.Sleep(1 * time.Second) // Пауза между заказами
	}

	fmt.Println("\n🎯 Эмуляция завершена!")
}

func createTestOrder(orderUID, customerName string) *model.Order {
	return &model.Order{
		OrderUID:    orderUID,
		TrackNumber: fmt.Sprintf("TRACK_%s", orderUID),
		Entry:       "WBIL",
		Delivery: model.Delivery{
			Name:    customerName,
			Phone:   "+79001234567",
			Zip:     "123456",
			City:    "Москва",
			Address: "ул. Тестовая, д. 1",
			Region:  "Московская область",
			Email:   fmt.Sprintf("%s@test.com", orderUID),
		},
		Payment: model.Payment{
			Transaction:  fmt.Sprintf("txn_%s", orderUID),
			RequestID:    fmt.Sprintf("req_%s", orderUID),
			Currency:     "RUB",
			Provider:     "test_provider",
			Amount:       1500,
			PaymentDt:    time.Now().Unix(),
			Bank:         "test_bank",
			DeliveryCost: 300,
			GoodsTotal:   1200,
			CustomFee:    0,
		},
		Items: []model.Item{
			{
				ChrtID:      12345,
				TrackNumber: fmt.Sprintf("TRACK_%s", orderUID),
				Price:       600,
				Rid:         fmt.Sprintf("rid_%s", orderUID),
				Name:        "Тестовый товар 1",
				Sale:        20,
				Size:        "M",
				TotalPrice:  480,
				NmID:        67890,
				Brand:       "Test Brand",
				Status:      202,
			},
			{
				ChrtID:      12346,
				TrackNumber: fmt.Sprintf("TRACK_%s", orderUID),
				Price:       600,
				Rid:         fmt.Sprintf("rid_%s_2", orderUID),
				Name:        "Тестовый товар 2",
				Sale:        0,
				Size:        "L",
				TotalPrice:  600,
				NmID:        67891,
				Brand:       "Test Brand",
				Status:      202,
			},
		},
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        orderUID,
		DeliveryService:   "test_service",
		ShardKey:          "1",
		SmID:              99,
		DateCreated:       time.Now(),
		OofShard:          "1",
	}
}
