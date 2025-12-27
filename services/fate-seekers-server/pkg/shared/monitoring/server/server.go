package server

import (
	"fmt"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/logging"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/monitoring/services"
	"github.com/fasthttp/router"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// MonitoringServer represents monitoring server, which serves metrics endpoint for prometheus.
type MonitoringServer struct {
}

// Start starts server, which serves metrics for monitoring components.
func (ms *MonitoringServer) Start(failover func()) {
	go func() {
		services.Init()

		r := router.New()

		r.GET("/metrics", fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler()))

		err := fasthttp.ListenAndServe(
			fmt.Sprintf("0.0.0.0:%v", config.GetSettingsMonitoringPrometheusPort()), r.Handler)
		if err != nil {
			logging.GetInstance().Error(err.Error())

			failover()

			return
		}
	}()
}

// NewMonitoringServer initializes MonitoringServer.
func NewMonitoringServer() *MonitoringServer {
	return new(MonitoringServer)
}
