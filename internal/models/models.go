package models

// Weather Models

type ForecastInput struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type AlertsInput struct {
	State string `json:"state"`
}

type PointsResponse struct {
	Properties struct {
		Forecast string `json:"forecast"`
	} `json:"properties"`
}

type ForecastPeriod struct {
	Name             string `json:"name"`
	Temperature      int    `json:"temperature"`
	TemperatureUnit  string `json:"temperatureUnit"`
	WindSpeed        string `json:"windSpeed"`
	WindDirection    string `json:"windDirection"`
	ShortForecast    string `json:"shortForecast"`
	DetailedForecast string `json:"detailedForecast"`
}

type ForecastResponse struct {
	Properties struct {
		Periods []ForecastPeriod `json:"periods"`
	} `json:"properties"`
}

type AlertFeature struct {
	Properties struct {
		Event       string `json:"event"`
		AreaDesc    string `json:"areaDesc"`
		Severity    string `json:"severity"`
		Description string `json:"description"`
		Instruction string `json:"instruction"`
	} `json:"properties"`
}

type AlertsResponse struct {
	Features []AlertFeature `json:"features"`
}

// Task Models

type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

type AddTaskInput struct {
	Title  string `json:"title"`
	Status string `json:"status"` // Optional
}

type UpdateTaskInput struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`  // Optional
	Status string `json:"status"` // Optional
}

type DeleteTaskInput struct {
	ID int `json:"id"`
}

type GetTaskInput struct {
	ID int `json:"id"`
}

type ListTasksInput struct{}
