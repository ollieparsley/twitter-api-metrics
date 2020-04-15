package main

import (
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var l *log.Logger = log.New(os.Stdout, "[twitter-api-metrics] ", 2)

// getEnv Get an env vairable and set a default
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv("TWITTER_API_METRICS_" + key); ok {
		return value
	}
	return fallback
}

func main() {
	l.Println("Starting service")

	// Fetch env variables with defaults
	name := getEnv("NAME", "default")
	httpPort := getEnv("HTTP_PORT", "9100")
	httpPath := getEnv("HTTP_PATH", "metrics")
	metricPrefix := getEnv("METRICS_PREFIX", "twitter_api_ratelimits")
	apiKey := getEnv("API_KEY", "")
	apiSecretKey := getEnv("API_SECRET_KEY", "")
	intervalSeconds := getEnv("INTERVAL_SECONDS", "10")
	intervalSecondsInt, err := strconv.Atoi(intervalSeconds)
	if err != nil {
		l.Printf("Error converting interval: %s", err.Error())
		os.Exit(1)
	}

	// oauth2 configures a client that uses app credentials to keep a fresh token
	config := &clientcredentials.Config{
		ClientID:     apiKey,
		ClientSecret: apiSecretKey,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth2.NoContext)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// All resources
	resources := []string{
		"application",
		"favorites",
		"followers",
		"friends",
		"friendships",
		"geo",
		"help",
		"lists",
		"search",
		"statuses",
		"trends",
		"users",
	}

	// API request counter counter
	counter := promauto.NewCounter(prometheus.CounterOpts{
		Name: metricPrefix + "_counter",
		Help: "Increment each time we call the API to get rate limits",
		ConstLabels: prometheus.Labels{
			"type": "",
		},
	})

	// Set up metric guages
	limitGauges := map[string]map[string]prometheus.Gauge{}
	remainingGauges := map[string]map[string]prometheus.Gauge{}
	for _, resource := range resources {
		limitGauges[resource] = map[string]prometheus.Gauge{}
		remainingGauges[resource] = map[string]prometheus.Gauge{}
	}

	// Increment counter every second
	go func() {
		for {
			// Show that we're making a request
			counter.Inc()

			l.Println("Fetching rate limit information")

			// Fetch rate limits
			rateLimits, _, err := client.RateLimits.Status(&twitter.RateLimitParams{Resources: []string{}})
			if err != nil {
				l.Printf("Error getting rate limits: %s", err.Error())
				os.Exit(1)
			}

			l.Println("Fetchied rate limit information")

			resourcesReflected := reflect.ValueOf(rateLimits.Resources)
			resourcesReflectedIndirect := reflect.Indirect(resourcesReflected)

			// Each resource group
			for _, resource := range resources {
				if _, ok := limitGauges[resource]; ok == false {
					limitGauges[resource] = map[string]prometheus.Gauge{}
					remainingGauges[resource] = map[string]prometheus.Gauge{}
				}
				resourceLimitEndpoints := resourcesReflectedIndirect.FieldByName(strings.Title(resource)).Interface().(map[string]*twitter.RateLimitResource)
				for endpoint, data := range resourceLimitEndpoints {
					if _, ok := limitGauges[resource][endpoint]; ok == false {
						limitGauges[resource][endpoint] = promauto.NewGauge(prometheus.GaugeOpts{
							Name: metricPrefix + "_limit",
							Help: "The request limit for the endpoint rate limit",
							ConstLabels: prometheus.Labels{
								"resource": resource,
								"endpoint": endpoint,
								"name":     name,
							},
						})
					}
					if _, ok := remainingGauges[resource][endpoint]; ok == false {
						remainingGauges[resource][endpoint] = promauto.NewGauge(prometheus.GaugeOpts{
							Name: metricPrefix + "_remaining",
							Help: "The current number of requests made in the period for the endpoint rate limit",
							ConstLabels: prometheus.Labels{
								"resource": resource,
								"endpoint": endpoint,
								"name":     name,
							},
						})
					}

					limitGauges[resource][endpoint].Set(float64(data.Limit))
					remainingGauges[resource][endpoint].Set(float64(data.Remaining))
				}
			}

			l.Println("Updated metrics")
			l.Printf("Waiting %d seconds until the next fetch time", intervalSecondsInt)

			// Wait for the next time period
			time.Sleep(time.Duration(intervalSecondsInt) * time.Second)
		}
	}()

	// Set up the HTTP handler and block
	addr := ":" + httpPort
	l.Printf("Setting up http service %s/%s", addr, httpPath)
	http.Handle("/"+httpPath, promhttp.Handler())
	http.ListenAndServe(addr, nil)
}
