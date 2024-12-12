package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Структура продукта питания
type FoodProduct struct {
	Name            string
	Price           float64
	Manufacturer    string
	ManufactureDate time.Time
	ShelfLife       time.Duration
}

// Функция для генерации случайных данных о продуктах
func generateFoodProducts(size int) []FoodProduct {
	products := make([]FoodProduct, size)
	for i := 0; i < size; i++ {
		manufactureDate := time.Now().AddDate(0, 0, -rand.Intn(30)) // Продукты с датой изготовления до 30 дней назад
		products[i] = FoodProduct{
			Name:            fmt.Sprintf("Product %d", i),
			Price:           float64(rand.Intn(100)) + rand.Float64(),
			Manufacturer:    fmt.Sprintf("Manufacturer %d", rand.Intn(10)),
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
	const numGoroutines = 5 // Количество горутин

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
	ch := make(chan []FoodProduct, numGoroutines)

	// Разделяем продукты на части и запускаем горутины
	chunkSize := (size + numGoroutines - 1) / numGoroutines // Размер порции для каждой горутины
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			start := i * chunkSize
			end := start + chunkSize
			if end > size {
				end = size
			}
			expiredChunk := expiredProducts(products[start:end])
			ch <- expiredChunk // Отправляем результаты в канал
		}(i)
	}

	// Закрываем канал после завершения всех горутин
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Собираем результаты из канала
	for expiredChunk := range ch {
		expiredConcurrent = append(expiredConcurrent, expiredChunk...)
	}

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
