package weather

type CurrentWeather struct {
	IsDay              bool
	ConditionText      string
	ConditionImageCode string

	Temperature   float64
	Wind          float64
	Gust          float64
	WindDegree    int
	WindDir       string
	Pressure      uint16
	Precipitation float64
	Humidity      uint8
	CloudPercent  int
	FeelsLike     float64
	Visibility    float64
	UV            float64
}

type ForecastItem struct {
	Date string

	ConditionText      string
	ConditionImageCode string

	MaxTemp            float64
	MinTemp            float64
	AvgTemp            float64
	MaxWind            float64
	TotalPrecipitation float64
	TotalSnow          float64
	AvgVis             float64
	AvgHumidity        float64
	DailyWillItRain    bool
	DailyChanceOfRain  int
	DailyWillItSnow    bool
	DailyChanceOfSnow  int
	UV                 float64
}

type Info interface {
	CurrentWeather() CurrentWeather
	Forecast() []ForecastItem
}
