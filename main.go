package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	apiKey      = "ваш Apikey"
	cacheFile   = "weather_cache.json"
	cacheExpiry = 10 * time.Minute
)

type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Name string    `json:"name"`
	Time time.Time `json:"time"` // Добавляем поле для хранения времени кэширования
}

func getWeather(city string) (*WeatherResponse, error) {
	// Проверяем кэш
	if cachedData, err := readCache(city); err == nil {
		return cachedData, nil
	}

	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weather WeatherResponse
	if err := json.Unmarshal(body, &weather); err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	if err := writeCache(city, &weather); err != nil {
		fmt.Println("Ошибка при сохранении в кэш:", err)
	}

	return &weather, nil
}

func readCache(city string) (*WeatherResponse, error) {
	file, err := os.Open(cacheFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cache map[string]WeatherResponse
	if err := json.NewDecoder(file).Decode(&cache); err != nil {
		return nil, err
	}

	if data, ok := cache[city]; ok {
		if time.Since(data.Time) < cacheExpiry {
			return &data, nil
		}
	}

	return nil, fmt.Errorf("данные не найдены в кэше или устарели")
}

func writeCache(city string, weather *WeatherResponse) error {
	cache := make(map[string]WeatherResponse)

	// Читаем существующий кэш
	if file, err := os.Open(cacheFile); err == nil {
		defer file.Close()
		json.NewDecoder(file).Decode(&cache)
	}

	// Устанавливаем время кэширования
	weather.Time = time.Now()

	// Обновляем кэш
	cache[city] = *weather

	// Записываем кэш в файл
	file, err := os.Create(cacheFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(cache)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: weather-cli <город>")
		return
	}

	city := os.Args[1]
	weather, err := getWeather(city)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Printf("Погода в городе %s: %.1f°C\n", weather.Name, weather.Main.Temp)
}
