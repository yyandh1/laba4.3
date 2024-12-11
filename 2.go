package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Структура продукта питания
type FoodProduct struct {
	Name           string
	Price          float64
	Manufacturer   string
	ManufactureDate time.Time
	ShelfLife      time.Duration
}

// Функция для генерации случайных данных о продуктах
func generateFoodProducts(size int) []FoodProduct {
	products := make([]FoodProduct, size)
	for i := 0; i < size; i++ {
		manufactureDate := time.Now().AddDate(0, 0, -rand.Intn(30)) // Продукты с датой изготовления до 30 дней назад
		products[i] = FoodProduct{
			Name:           fmt.Sprintf("Product %d", i),
			Price:          float64(rand.Intn(100)) + rand.Float64(),
			Manufacturer:   fmt.Sprintf("Manufacturer %d", rand.Intn(10)),
			ManufactureDate: manufactureDate,
			ShelfLife:      time.Hour * 24 * time.Duration(rand.Intn(30)), // Срок годности до 30 дней
		}
	}
	return products
}

// Функция для получения списка продуктов с истекшим сроком годности
func expiredProducts(products []FoodProduct) []FoodProduct {
	var expired []FoodProduct
	for _, product := range products {
		if time.Now().After(product.ManufactureDate.Add(product.ShelfLife)) {
			expired = append(expired, product)
		}
	}
	return expired
}

func main() {
	// Параметры задачи
	const size = 100 // Размер массива данных

	// Генерация случайных данных о продуктах
	products := generateFoodProducts(size)

	// Измерение времени без многопоточности
	startTime := time.Now()
	expired := expiredProducts(products)
	elapsedNoConcurrency := time.Since(startTime)

	// Вывод результатов без многопоточности
	fmt.Println("Expired Food Products (No Concurrency):")
	for _, product := range expired {
		fmt.Printf("Name: %s, Price: %.2f, Manufacturer: %s, Manufacture Date: %s, Shelf Life: %d days\n",
			product.Name, product.Price, product.Manufacturer, product.ManufactureDate.Format("2006-01-02"), int(product.ShelfLife.Hours()/24))
	}

	// Измерение времени с многопоточностью
	startTime = time.Now()
	var wg sync.WaitGroup
	var expiredConcurrent []FoodProduct

	wg.Add(1) // Указываем, что ожидаем завершения одной горутины
	go func() {
		defer wg.Done()
		expiredConcurrent = expiredProducts(products)
	}()
	wg.Wait()
	elapsedWithConcurrency := time.Since(startTime)

	// Вывод результатов с многопоточностью
	fmt.Println("\nExpired Food Products (With Concurrency):")
	for _, product := range expiredConcurrent {
		fmt.Printf("Name: %s, Price: %.2f, Manufacturer: %s, Manufacture Date: %s, Shelf Life: %d days\n",
			product.Name, product.Price, product.Manufacturer, product.ManufactureDate.Format("2006-01-02"), int(product.ShelfLife.Hours()/24))
	}

	// Вывод времени обработки
	fmt.Printf("\nTime without concurrency: %d microseconds\n", elapsedNoConcurrency.Microseconds())
	fmt.Printf("Time with concurrency: %d microseconds\n", elapsedWithConcurrency.Microseconds())
}
