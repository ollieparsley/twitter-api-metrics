package main

import (
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
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var logger *logrus.Logger = logrus.New()

// getEnv Get an env vairable and set a default
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv("TWITTER_API_METRICS_" + key); ok {
		return value
	}
	return fallback
}

func main() {
	// Setup the logger
	logLevel, err := logrus.ParseLevel(getEnv("LOG_LEVEL", "INFO"))
	if err != nil {
		logLevel = logrus.InfoLevel
	}

	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logLevel)

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
		logger.Errorf("Error converting interval: %s", err.Error())
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

			logger.Debug("Fetching rate limit information")

			// Fetch rate limits
			rateLimits, resp, err := client.RateLimits.Status(&twitter.RateLimitParams{Resources: []string{}})
			if resp != nil && resp.StatusCode == 429 {
				logger.Errorf("Hit rate limit when getting rate limits (ironic huh?). Sleeping until reset time: %s", err.Error())
				resetTimeString := resp.Header.Get("X-Rate-Limit-Reset")
				resetTime := time.Now().Add(15 * time.Minute)
				if resetTimeString != "" {
					i, err := strconv.ParseInt(resetTimeString, 10, 64)
					if err != nil {
						logger.Errorf("X-Rate-Limit-Reset header was not a valid integer: %s error: %s", resetTimeString, err.Error())
					}
					resetTime = time.Unix(i, 0)
				}

				logger.Debugf("Sleeping until rate limit reset at %s", resetTime.Format(time.RFC1123Z))
				time.Sleep(resetTime.Sub(time.Now()))
				continue
			}
			if err != nil {
				logger.Errorf("Error getting rate limits. Sleeping until next check: %s", err.Error())
				time.Sleep(time.Duration(intervalSecondsInt) * time.Second)
				continue
			}

			logger.Debug("Finished fetching rate limit information")

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

			logger.Debug("Updated metrics")
			logger.Debugf("Waiting %d seconds until the next fetch time", intervalSecondsInt)

			// Wait for the next time period
			time.Sleep(time.Duration(intervalSecondsInt) * time.Second)
		}
	}()

	// Set up the HTTP handler and block
	logger.Infof("Listening for HTTP requests at: 0.0.0.0:%v", httpPort)
	http.Handle("/"+httpPath, promhttp.Handler())
	http.ListenAndServe(":"+httpPort, nil)
}
