# Weather CLI

Простая утилита командной строки на Go для получения текущей погоды по названию города с использованием API OpenWeatherMap.

## Установка

1. Убедитесь, что у вас установлен Go (https://golang.org/dl/).
2. Клонируйте репозиторий:

   ```bash
   git clone https://github.com/ваш-username/weather-cli.git
   cd weather-cli
   ```

3. Установите зависимости:

   ```bash
   go mod download
   ```

4. Соберите проект:

   ```bash
   go build -o weather-cli
   ```

## Использование

Запустите программу, указав название города:

```bash
./weather-cli Moscow
