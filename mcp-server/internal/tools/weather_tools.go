package tools

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"mcp-go-demo/internal/models"
	"mcp-go-demo/internal/service"
)

type WeatherTools struct {
	Logger *slog.Logger
}

func (w *WeatherTools) Register(server *mcp.Server) {
	server.AddTool(mcp.NewTool("get_forecast",
		mcp.WithDescription("Get weather forecast for a given latitude and longitude"),
		mcp.WithNumber("latitude", mcp.Required(), mcp.Description("Latitude of the location")),
		mcp.WithNumber("longitude", mcp.Required(), mcp.Description("Longitude of the location")),
	), w.getForecast)

	server.AddTool(mcp.NewTool("get_alerts",
		mcp.WithDescription("Get active weather alerts for a given US state"),
		mcp.WithString("state", mcp.Required(), mcp.Description("Two-letter US state abbreviation (e.g., CA, NY)")),
	), w.getAlerts)
}

func (w *WeatherTools) getForecast(ctx context.Context, req *mcp.CallToolRequest, input models.ForecastInput) (
	*mcp.CallToolResult, any, error,
) {
	traceID := uuid.New().String()
	start := time.Now()
	reqLogger := w.Logger.With(slog.String("trace_id", traceID), slog.String("tool", "get_forecast"))
	reqLogger.Info("Tool execution started", slog.Any("input_args", input))

	// Get points data
	pointsURL := fmt.Sprintf("%s/points/%f,%f", service.NWSAPIBase, input.Latitude, input.Longitude)
	reqLogger.Debug("Fetching points data", slog.String("url", pointsURL))
	pointsData, err := service.MakeNWSRequest[models.PointsResponse](ctx, pointsURL)
	if err != nil {
		reqLogger.Error("Unable to fetch points data", slog.String("error", err.Error()))
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Unable to fetch forecast data for this location."},
			},
		}, nil, nil
	}

	// Get forecast data
	forecastURL := pointsData.Properties.Forecast
	if forecastURL == "" {
		reqLogger.Error("Forecast URL is empty")
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Unable to fetch forecast URL."},
			},
		}, nil, nil
	}

	reqLogger.Debug("Fetching detailed forecast data", slog.String("url", forecastURL))
	forecastData, err := service.MakeNWSRequest[models.ForecastResponse](ctx, forecastURL)
	if err != nil {
		reqLogger.Error("Unable to fetch detailed forecast data", slog.String("error", err.Error()))
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Unable to fetch detailed forecast."},
			},
		}, nil, nil
	}

	// Format the periods
	periods := forecastData.Properties.Periods
	if len(periods) == 0 {
		reqLogger.Warn("No forecast periods available")
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "No forecast periods available."},
			},
		}, nil, nil
	}

	// Show next 5 periods
	var forecasts []string
	for i := range min(5, len(periods)) {
		forecasts = append(forecasts, service.FormatPeriod(periods[i]))
	}

	result := strings.Join(forecasts, "\n---\n")

	reqLogger.Info("Tool execution completed", slog.Duration("latency", time.Since(start)))
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, nil, nil
}

func (w *WeatherTools) getAlerts(ctx context.Context, req *mcp.CallToolRequest, input models.AlertsInput) (
	*mcp.CallToolResult, any, error,
) {
	traceID := uuid.New().String()
	start := time.Now()
	reqLogger := w.Logger.With(slog.String("trace_id", traceID), slog.String("tool", "get_alerts"))
	reqLogger.Info("Tool execution started", slog.Any("input_args", input))

	// Build alerts URL
	stateCode := strings.ToUpper(input.State)
	alertsURL := fmt.Sprintf("%s/alerts/active/area/%s", service.NWSAPIBase, stateCode)

	reqLogger.Debug("Fetching alerts data", slog.String("url", alertsURL))
	alertsData, err := service.MakeNWSRequest[models.AlertsResponse](ctx, alertsURL)
	if err != nil {
		reqLogger.Error("Unable to fetch alerts data", slog.String("error", err.Error()))
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Unable to fetch alerts or no alerts found."},
			},
		}, nil, nil
	}

	// Check if there are any alerts
	if len(alertsData.Features) == 0 {
		reqLogger.Info("Tool execution completed (no alerts)")
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "No active alerts for this state."},
			},
		}, nil, nil
	}

	// Format alerts
	var alerts []string
	for _, feature := range alertsData.Features {
		alerts = append(alerts, service.FormatAlert(feature))
	}

	result := strings.Join(alerts, "\n---\n")

	reqLogger.Info("Tool execution completed", slog.Int("alert_count", len(alerts)), slog.Duration("latency", time.Since(start)))
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, nil, nil
}
