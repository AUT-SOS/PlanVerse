package models

import "github.com/prometheus/client_golang/prometheus"

var (
	SuccessRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "successful_requests",
			Help: "Total number of successful requests.",
		},
		[]string{"method", "endpoint"},
	)
	FailedRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "failed_requests",
			Help: "Total number of failed requests.",
		},
		[]string{"method", "endpoint"},
	)
	SuccessDBRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "successful_db_requests",
			Help: "Total number of successful database requests.",
		},
		[]string{"method", "endpoint"},
	)
	FailedDBRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "failed_db_requests",
			Help: "Total number of failed database requests.",
		},
		[]string{"method", "endpoint"},
	)
	ResponseTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "response_time",
			Help:    "Response time in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)
