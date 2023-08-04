package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"sync/atomic"
	"time"
)

const (
	currentAPI  = "https://api.weatherapi.com/v1/current.json?key=f0d704ae336f450ba55152722232907&q=59.891740,30.319351&aqi=no"
	forecastAPI = "https://api.weatherapi.com/v1/forecast.json?key=f0d704ae336f450ba55152722232907&q=59.891740,30.319351&days=3&aqi=no&alerts=no"

	currentWeatherUpdateInterval = 5 * time.Minute
	forecastUpdateInterval       = 15 * time.Minute
)

var (
	_ Info = &Weather{}
)

func New() *Weather {
	weather := &Weather{
		stopped:            make(chan struct{}),
		currentWeatherDone: make(chan struct{}),
		forecastDone:       make(chan struct{}),
	}

	weather.currentWeather.Store(&CurrentWeather{})
	forecast := make([]ForecastItem, 0)
	weather.forecast.Store(&forecast)

	return weather
}

type Weather struct {
	stopped            chan struct{}
	currentWeatherDone chan struct{}
	forecastDone       chan struct{}

	currentWeather atomic.Pointer[CurrentWeather]
	forecast       atomic.Pointer[[]ForecastItem]

	conditions map[ConditionCode]Condition
}

func (w *Weather) Init() error {
	var err error

	w.conditions, err = parseConditions()
	if err != nil {
		return err
	}

	go w.workerCurrentWeather()
	go w.workerForecast()

	return nil
}

func (w *Weather) Stop() {
	close(w.stopped)

	<-w.forecastDone
	<-w.currentWeatherDone
}

func (w *Weather) CurrentWeather() CurrentWeather {
	return *w.currentWeather.Load()
}

func (w *Weather) Forecast() []ForecastItem {
	return *w.forecast.Load()
}

func (w *Weather) workerCurrentWeather() {
	w.getCurrentWeather()

	ticker := time.NewTicker(currentWeatherUpdateInterval)

workerLoop:
	for {
		select {
		case <-w.stopped:
			break workerLoop
		case <-ticker.C:
			w.getCurrentWeather()
		}
	}

	ticker.Stop()
	fmt.Println("workerCurrentWeather stopped")
	close(w.currentWeatherDone)
}

func (w *Weather) workerForecast() {
	w.getForecast()

	ticker := time.NewTicker(forecastUpdateInterval)

workerLoop:
	for {
		select {
		case <-w.stopped:
			break workerLoop
		case <-ticker.C:
			w.getForecast()
		}
	}

	ticker.Stop()
	fmt.Println("workerForecast stopped")
	close(w.forecastDone)
}

func (w *Weather) getCurrentWeather() {
	data, err := w.performGet(currentAPI)
	if err != nil {
		log.Printf("failed to get current weather: %v", err)
		return
	}

	curParsed := currentResponse{}
	err = json.Unmarshal(data, &curParsed)
	if err != nil {
		log.Printf("failed to parse current weather response '%v': %v", string(data), err)
		return
	}

	w.currentWeather.Store(w.fillCurrentWeather(&curParsed))
	log.Printf("Current weather received successfully")
}

func (w *Weather) getForecast() {
	data, err := w.performGet(forecastAPI)
	if err != nil {
		log.Printf("failed to get current weather: %v", err)
		return
	}

	forecastParsed := forecastResponse{}
	err = json.Unmarshal(data, &forecastParsed)
	if err != nil {
		log.Printf("failed to parse weather forecast response '%v': %v", string(data), err)
		return
	}

	w.forecast.Store(w.fillForecast(&forecastParsed))
	log.Printf("Weather forecast received successfully")
}

func (w *Weather) performGet(url string) ([]byte, error) {
	log.Printf("Performing GET on '%s'", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to perform GET: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read GET response body: %w", err)
	}

	return body, nil
}

type conditionDesc struct {
	text  string
	image string
}

func (w *Weather) getCondition(code ConditionCode, isDay bool) conditionDesc {
	cond := w.conditions[code]

	var condDesc conditionDesc
	if isDay {
		condDesc.text = cond.Day
	} else {
		condDesc.text = cond.Night
	}

	condDesc.image = fmt.Sprintf("%d", cond.Icon)

	return condDesc
}

func (w *Weather) fillCurrentWeather(response *currentResponse) *CurrentWeather {
	isDay := response.Current.IsDay > 0
	cond := w.getCondition(response.Current.Condition.Code, isDay)

	return &CurrentWeather{
		IsDay:              isDay,
		ConditionText:      cond.text,
		ConditionImageCode: cond.image,

		Temperature:   response.Current.TempC,
		Wind:          kphToMps(response.Current.WindKph),
		Gust:          kphToMps(response.Current.GustKph),
		WindDegree:    response.Current.WindDegree,
		WindDir:       response.Current.WindDir,
		Pressure:      mBarToMmHg(response.Current.PressureMb),
		Precipitation: response.Current.PrecipitationMm,
		Humidity:      uint8(response.Current.Humidity),
		CloudPercent:  response.Current.Cloud,
		FeelsLike:     response.Current.FeelsLikeC,
		Visibility:    response.Current.VisKm,
		UV:            response.Current.UV,
	}
}

func (w *Weather) fillForecast(response *forecastResponse) *[]ForecastItem {
	forecast := make([]ForecastItem, 0, len(response.Forecast.Forecastday))

	for _, dayForecast := range response.Forecast.Forecastday {
		cond := w.getCondition(dayForecast.Day.Condition.Code, true)
		forecast = append(forecast, ForecastItem{
			Date:               dayForecast.Date,
			ConditionText:      cond.text,
			ConditionImageCode: cond.image,
			MaxTemp:            dayForecast.Day.MaxTempC,
			MinTemp:            dayForecast.Day.MinTempC,
			AvgTemp:            dayForecast.Day.AvgTempC,
			MaxWind:            kphToMps(dayForecast.Day.MaxWindKph),
			TotalPrecipitation: dayForecast.Day.TotalPrecipitationMm,
			TotalSnow:          dayForecast.Day.TotalSnowCm * 10,
			AvgVis:             dayForecast.Day.AvgVisKm,
			AvgHumidity:        dayForecast.Day.AvgHumidity,
			DailyWillItRain:    dayForecast.Day.DailyWillItRain > 0,
			DailyChanceOfRain:  dayForecast.Day.DailyChanceOfRain,
			DailyWillItSnow:    dayForecast.Day.DailyWillItSnow > 0,
			DailyChanceOfSnow:  dayForecast.Day.DailyChanceOfSnow,
			UV:                 dayForecast.Day.UV,
		})
	}

	return &forecast
}

func parseConditions() (map[ConditionCode]Condition, error) {
	conditions := make([]Condition, 0)
	err := json.Unmarshal(conditionsData, &conditions)
	if err != nil {
		return nil, fmt.Errorf("failed to parse conditions: %w", err)
	}

	condMap := make(map[ConditionCode]Condition, len(conditions))
	for _, cond := range conditions {
		condMap[cond.Code] = cond
	}

	return condMap, nil
}

func mBarToMmHg(val float64) uint16 {
	return uint16(math.Round(val / 1.33322387415))
}

func kphToMps(val float64) float64 {
	return (val * 1000) / 3600
}
