package main

import (
	"log"
	"time"
	"wb_cource/internal/app/kafka"
	"wb_cource/internal/app/model"

	"github.com/brianvoe/gofakeit/v6"
)

func main() {
	gofakeit.Seed(time.Now().UnixNano())

	producer, err := kafka.NewProducer("localhost:9092", "orders")
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	log.Println("Kafka Producer эмулятор запущен")

	for i := 0; i < 20; i++ {
		order := generateOrder()

		err := producer.SendOrder(order)
		if err != nil {
			log.Printf("Ошибка отправки заказа %s: %v", order.OrderUID, err)
		} else {
			log.Printf("Заказ %s отправлен в Kafka", order.OrderUID)
		}

		time.Sleep(2 * time.Second)
	}

	log.Println("Эмуляция producer завершена")
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
