package services

import "github.com/prometheus/client_golang/prometheus"

// Represents prometheus available sessions.
var (
	availableSessions = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "available_sessions",
			Help: "The current number of available sessions",
		},
	)

	availableLobbies = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "available_lobbies",
			Help: "The current number of available lobbies",
		},
	)
)

// IncAvailableSession performs available session value incrementation.
func IncAvailableSession() {
	availableSessions.Inc()
}

// DecAvailableSession performs available session value decremention.
func DecAvailableSession() {
	availableSessions.Dec()
}

// SetAvailableSession performs available session value setup.
func SetAvailableSession(value int64) {
	availableSessions.Set(float64(value))
}

// IncAvailableLobby performs available lobby value incrementation.
func IncAvailableLobby() {
	availableLobbies.Inc()
}

// DecAvailableLobby performs available lobby value decremention.
func DecAvailableLobby() {
	availableLobbies.Dec()
}

// SetAvailableLobby performs available lobby value setup.
func SetAvailableLobby(value int64) {
	availableLobbies.Set(float64(value))
}

// Init performs registers initialization.
func Init() {
	prometheus.MustRegister(
		availableSessions, availableLobbies)
}
