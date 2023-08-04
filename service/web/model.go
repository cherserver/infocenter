package web

type Sensor struct {
	SID            string   `json:"sid"`
	BatteryPercent *uint8   `json:"battery_percent,omitempty"`
	Temperature    *float32 `json:"temperature,omitempty"`
	Humidity       *float32 `json:"humidity,omitempty"`
	Pressure       *float32 `json:"pressure,omitempty"`
}

type Weather struct {
	Current  CurrentWeather `json:"current"`
	Forecast []ForecastItem `json:"forecast"`
}

type CurrentWeather struct {
	ConditionText  string  `json:"condition_text"`
	ConditionImage string  `json:"condition_image"`
	Temperature    float64 `json:"temperature"`
	Wind           float64 `json:"wind"`
	Gust           float64 `json:"gust"`
	WindDegree     int     `json:"wind_degree"`
	WindDir        string  `json:"wind_dir"`
	Pressure       uint16  `json:"pressure"`
	Precipitation  float64 `json:"precipitation"`
	Humidity       uint8   `json:"humidity"`
	CloudPercent   int     `json:"cloud_percent"`
	FeelsLike      float64 `json:"feels_like"`
	Visibility     float64 `json:"visibility"`
	UV             float64 `json:"uv"`
}

type ForecastItem struct {
	ConditionText      string  `json:"condition_text"`
	ConditionImage     string  `json:"condition_image"`
	MaxTemp            float64 `json:"max_temp"`
	MinTemp            float64 `json:"min_temp"`
	AvgTemp            float64 `json:"avg_temp"`
	MaxWind            float64 `json:"max_wind"`
	TotalPrecipitation float64 `json:"total_precipitation"`
	TotalSnow          float64 `json:"total_snow"`
	AvgVis             float64 `json:"avg_vis"`
	AvgHumidity        float64 `json:"avg_humidity"`
	DailyWillItRain    bool    `json:"daily_will_it_rain"`
	DailyChanceOfRain  int     `json:"daily_chance_of_rain"`
	DailyWillItSnow    bool    `json:"daily_will_it_snow"`
	DailyChanceOfSnow  int     `json:"daily_chance_of_snow"`
	UV                 float64 `json:"uv"`
}
