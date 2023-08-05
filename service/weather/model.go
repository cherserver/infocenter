package weather

import _ "embed"

type ConditionCode int

type currentResponse struct {
	/*Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`*/
	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string        `json:"text"`
			Icon string        `json:"icon"`
			Code ConditionCode `json:"code"`
		} `json:"condition"`
		WindKph         float64 `json:"wind_kph"`
		WindDegree      int     `json:"wind_degree"`
		WindDir         string  `json:"wind_dir"`
		PressureMb      float64 `json:"pressure_mb"`
		PrecipitationMm float64 `json:"precip_mm"`
		Humidity        int     `json:"humidity"`
		Cloud           int     `json:"cloud"`
		FeelsLikeC      float64 `json:"feelslike_c"`
		VisKm           float64 `json:"vis_km"`
		UV              float64 `json:"uv"`
		GustKph         float64 `json:"gust_kph"`
	} `json:"current"`
}

type forecastResponse struct {
	/*Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzId           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string        `json:"text"`
			Icon string        `json:"icon"`
			Code ConditionCode `json:"code"`
		} `json:"condition"`
		WindKph         float64 `json:"wind_kph"`
		WindDegree      int     `json:"wind_degree"`
		WindDir         string  `json:"wind_dir"`
		PressureMb      float64 `json:"pressure_mb"`
		PrecipitationMm float64 `json:"precip_mm"`
		Humidity        int     `json:"humidity"`
		Cloud           int     `json:"cloud"`
		FeelsLikeC      float64 `json:"feelslike_c"`
		VisKm           float64 `json:"vis_km"`
		UV              float64 `json:"uv"`
		GustKph         float64 `json:"gust_kph"`
	} `json:"current"`*/
	Forecast struct {
		Forecastday []struct {
			Date      string `json:"date"`
			DateEpoch int    `json:"date_epoch"`
			Day       struct {
				MaxTempC             float64 `json:"maxtemp_c"`
				MinTempC             float64 `json:"mintemp_c"`
				AvgTempC             float64 `json:"avgtemp_c"`
				MaxWindKph           float64 `json:"maxwind_kph"`
				TotalPrecipitationMm float64 `json:"totalprecip_mm"`
				TotalSnowCm          float64 `json:"totalsnow_cm"`
				AvgVisKm             float64 `json:"avgvis_km"`
				AvgHumidity          float64 `json:"avghumidity"`
				DailyWillItRain      int     `json:"daily_will_it_rain"`
				DailyChanceOfRain    int     `json:"daily_chance_of_rain"`
				DailyWillItSnow      int     `json:"daily_will_it_snow"`
				DailyChanceOfSnow    int     `json:"daily_chance_of_snow"`
				Condition            struct {
					Text string        `json:"text"`
					Icon string        `json:"icon"`
					Code ConditionCode `json:"code"`
				} `json:"condition"`
				UV float64 `json:"uv"`
			} `json:"day"`
			/*Astro struct {
				Sunrise          string `json:"sunrise"`
				Sunset           string `json:"sunset"`
				Moonrise         string `json:"moonrise"`
				Moonset          string `json:"moonset"`
				MoonPhase        string `json:"moon_phase"`
				MoonIllumination string `json:"moon_illumination"`
				IsMoonUp         int    `json:"is_moon_up"`
				IsSunUp          int    `json:"is_sun_up"`
			} `json:"astro"`
			Hour []struct {
				TimeEpoch int     `json:"time_epoch"`
				Time      string  `json:"time"`
				TempC     float64 `json:"temp_c"`
				IsDay     int     `json:"is_day"`
				Condition struct {
					Text string        `json:"text"`
					Icon string        `json:"icon"`
					Code ConditionCode `json:"code"`
				} `json:"condition"`
				WindKph         float64 `json:"wind_kph"`
				WindDegree      int     `json:"wind_degree"`
				WindDir         string  `json:"wind_dir"`
				PressureMb      float64 `json:"pressure_mb"`
				PrecipitationMm float64 `json:"precip_mm"`
				Humidity        int     `json:"humidity"`
				Cloud           int     `json:"cloud"`
				FeelsLikeC      float64 `json:"feelslike_c"`
				WindchillC      float64 `json:"windchill_c"`
				HeatIndexC      float64 `json:"heatindex_c"`
				DewPointC       float64 `json:"dewpoint_c"`
				WillItRain      int     `json:"will_it_rain"`
				ChanceOfRain    int     `json:"chance_of_rain"`
				WillItSnow      int     `json:"will_it_snow"`
				ChanceOfSnow    int     `json:"chance_of_snow"`
				VisKm           float64 `json:"vis_km"`
				GustKph         float64 `json:"gust_kph"`
				UV              float64 `json:"uv"`
			} `json:"hour"`*/
		} `json:"forecastday"`
	} `json:"forecast"`
}

//go:embed conditions.json
var conditionsData []byte

type Condition struct {
	Code      ConditionCode `json:"code"`
	Day       string        `json:"day"`
	Night     string        `json:"night"`
	Icon      int           `json:"icon"`
	Languages []struct {
		LangName  string `json:"lang_name"`
		LangIso   string `json:"lang_iso"`
		DayText   string `json:"day_text"`
		NightText string `json:"night_text"`
	} `json:"languages"`
}
