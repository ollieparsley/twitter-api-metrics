# Twitter API Metrics
Monitor your Twitter API usage using app auth and outputting to prometheus. This is done by calling the Twitter API on a regular interval and updating the prometheus metrics. The Twitter API endpoint `application/rate_limit_status`. A metric key is created for following:

- Rate limit for each Twitter API endpoint for your app
- Remaining requests for each Twitter API endpoint for your app

[![Build Status](https://cloud.drone.io/api/badges/ollieparsley/twitter-api-metrics/status.svg)](https://cloud.drone.io/ollieparsley/twitter-api-metrics)

## Usage

### Docker

We have a Docker image for you to use on docker hub. By default you'll need to expose port `9100`. Metric would then be available at `/metrics`

```
ollieparsley/twitter-api-metrics
```

### Environment variables

To allow you to customise the service we have a number of environment variables you can use so that you don't need to edit the code yourself:

| Name                                   | Default  | Description                                                                                             |
|----------------------------------------|--------------------------|-----------------------------------------------------------------------------------------|
| `TWITTER_API_METRICS_HTTP_PORT`        | `9100`                   | The HTTP port that you will need to expose if running in docker                         |
| `TWITTER_API_METRICS_HTTP_PATH`        | `/metrics`               | The HTTP endpoint to to expose prometheus metrics on                                    |
| `TWITTER_API_METRICS_API_KEY`          | `Required`               | The Twitter app API Key (sometimes called consumer key)                                 |
| `TWITTER_API_METRICS_API_SECRET_KEY`   | `Required`               | The Twitter app API Secret Key (sometimes called consumer secret)                       |
| `TWITTER_API_METRICS_INTERVAL_SECONDS` | `10`                     | The interval used to refresh the twitter api rate limit metrics                         |
| `TWITTER_API_METRICS_METRICS_PREFIX`   | `twitter_api_ratelimits` | The prefix for the prometheus metrics names                                             |
| `TWITTER_API_METRICS_NAME`             | `default`                | A name to put into prometheus metric keys to help identify, if using multiple instances |
| `TWITTER_API_METRICS_LOG_LEVEL`        | `INFO`                   | Set the log level, defaults to INFO                                                     |

## Prometheus output

This is an example of the output you get from the service

```
# HELP twitter_api_ratelimits_counter Increment each time we call the API to get rate limits
# TYPE twitter_api_ratelimits_counter counter
twitter_api_ratelimits_counter{type=""} 2
# HELP twitter_api_ratelimits_limit The request limit for the endpoint rate limit
# TYPE twitter_api_ratelimits_limit gauge
twitter_api_ratelimits_limit{endpoint="/application/rate_limit_status",resource="application"} 180
twitter_api_ratelimits_limit{endpoint="/favorites/list",resource="favorites"} 75
twitter_api_ratelimits_limit{endpoint="/followers/ids",resource="followers"} 15
twitter_api_ratelimits_limit{endpoint="/followers/list",resource="followers"} 30
twitter_api_ratelimits_limit{endpoint="/friends/following/ids",resource="friends"} 15
twitter_api_ratelimits_limit{endpoint="/friends/following/list",resource="friends"} 15
twitter_api_ratelimits_limit{endpoint="/friends/ids",resource="friends"} 15
twitter_api_ratelimits_limit{endpoint="/friends/list",resource="friends"} 30
twitter_api_ratelimits_limit{endpoint="/friendships/list",resource="friendships"} 15
twitter_api_ratelimits_limit{endpoint="/friendships/show",resource="friendships"} 15
twitter_api_ratelimits_limit{endpoint="/help/configuration",resource="help"} 15
twitter_api_ratelimits_limit{endpoint="/help/languages",resource="help"} 15
twitter_api_ratelimits_limit{endpoint="/help/privacy",resource="help"} 15
twitter_api_ratelimits_limit{endpoint="/help/settings",resource="help"} 15
twitter_api_ratelimits_limit{endpoint="/help/tos",resource="help"} 15
twitter_api_ratelimits_limit{endpoint="/lists/list",resource="lists"} 15
twitter_api_ratelimits_limit{endpoint="/lists/members",resource="lists"} 75
twitter_api_ratelimits_limit{endpoint="/lists/members/show",resource="lists"} 15
twitter_api_ratelimits_limit{endpoint="/lists/memberships",resource="lists"} 75
twitter_api_ratelimits_limit{endpoint="/lists/ownerships",resource="lists"} 15
twitter_api_ratelimits_limit{endpoint="/lists/show",resource="lists"} 75
twitter_api_ratelimits_limit{endpoint="/lists/statuses",resource="lists"} 900
twitter_api_ratelimits_limit{endpoint="/lists/subscribers",resource="lists"} 15
twitter_api_ratelimits_limit{endpoint="/lists/subscribers/show",resource="lists"} 15
twitter_api_ratelimits_limit{endpoint="/lists/subscriptions",resource="lists"} 15
twitter_api_ratelimits_limit{endpoint="/search/tweets",resource="search"} 450
twitter_api_ratelimits_limit{endpoint="/statuses/lookup",resource="statuses"} 300
twitter_api_ratelimits_limit{endpoint="/statuses/oembed",resource="statuses"} 180
twitter_api_ratelimits_limit{endpoint="/statuses/retweeters/ids",resource="statuses"} 300
twitter_api_ratelimits_limit{endpoint="/statuses/retweets/:id",resource="statuses"} 300
twitter_api_ratelimits_limit{endpoint="/statuses/show/:id",resource="statuses"} 900
twitter_api_ratelimits_limit{endpoint="/statuses/user_timeline",resource="statuses"} 1500
twitter_api_ratelimits_limit{endpoint="/trends/available",resource="trends"} 75
twitter_api_ratelimits_limit{endpoint="/trends/closest",resource="trends"} 75
twitter_api_ratelimits_limit{endpoint="/trends/place",resource="trends"} 75
twitter_api_ratelimits_limit{endpoint="/users/lookup",resource="users"} 300
twitter_api_ratelimits_limit{endpoint="/users/profile_banner",resource="users"} 180
twitter_api_ratelimits_limit{endpoint="/users/reverse_lookup",resource="users"} 15
twitter_api_ratelimits_limit{endpoint="/users/show/:id",resource="users"} 900
twitter_api_ratelimits_limit{endpoint="/users/suggestions",resource="users"} 15
twitter_api_ratelimits_limit{endpoint="/users/suggestions/:slug",resource="users"} 15
twitter_api_ratelimits_limit{endpoint="/users/suggestions/:slug/members",resource="users"} 15
# HELP twitter_api_ratelimits_remaining The current number of requests made in the period for the endpoint rate limit
# TYPE twitter_api_ratelimits_remaining gauge
twitter_api_ratelimits_remaining{endpoint="/application/rate_limit_status",resource="application"} 114
twitter_api_ratelimits_remaining{endpoint="/favorites/list",resource="favorites"} 75
twitter_api_ratelimits_remaining{endpoint="/followers/ids",resource="followers"} 15
twitter_api_ratelimits_remaining{endpoint="/followers/list",resource="followers"} 30
twitter_api_ratelimits_remaining{endpoint="/friends/following/ids",resource="friends"} 15
twitter_api_ratelimits_remaining{endpoint="/friends/following/list",resource="friends"} 15
twitter_api_ratelimits_remaining{endpoint="/friends/ids",resource="friends"} 15
twitter_api_ratelimits_remaining{endpoint="/friends/list",resource="friends"} 30
twitter_api_ratelimits_remaining{endpoint="/friendships/list",resource="friendships"} 15
twitter_api_ratelimits_remaining{endpoint="/friendships/show",resource="friendships"} 15
twitter_api_ratelimits_remaining{endpoint="/help/configuration",resource="help"} 15
twitter_api_ratelimits_remaining{endpoint="/help/languages",resource="help"} 15
twitter_api_ratelimits_remaining{endpoint="/help/privacy",resource="help"} 15
twitter_api_ratelimits_remaining{endpoint="/help/settings",resource="help"} 15
twitter_api_ratelimits_remaining{endpoint="/help/tos",resource="help"} 15
twitter_api_ratelimits_remaining{endpoint="/lists/list",resource="lists"} 15
twitter_api_ratelimits_remaining{endpoint="/lists/members",resource="lists"} 75
twitter_api_ratelimits_remaining{endpoint="/lists/members/show",resource="lists"} 15
twitter_api_ratelimits_remaining{endpoint="/lists/memberships",resource="lists"} 75
twitter_api_ratelimits_remaining{endpoint="/lists/ownerships",resource="lists"} 15
twitter_api_ratelimits_remaining{endpoint="/lists/show",resource="lists"} 75
twitter_api_ratelimits_remaining{endpoint="/lists/statuses",resource="lists"} 900
twitter_api_ratelimits_remaining{endpoint="/lists/subscribers",resource="lists"} 15
twitter_api_ratelimits_remaining{endpoint="/lists/subscribers/show",resource="lists"} 15
twitter_api_ratelimits_remaining{endpoint="/lists/subscriptions",resource="lists"} 15
twitter_api_ratelimits_remaining{endpoint="/search/tweets",resource="search"} 450
twitter_api_ratelimits_remaining{endpoint="/statuses/lookup",resource="statuses"} 300
twitter_api_ratelimits_remaining{endpoint="/statuses/oembed",resource="statuses"} 180
twitter_api_ratelimits_remaining{endpoint="/statuses/retweeters/ids",resource="statuses"} 300
twitter_api_ratelimits_remaining{endpoint="/statuses/retweets/:id",resource="statuses"} 300
twitter_api_ratelimits_remaining{endpoint="/statuses/show/:id",resource="statuses"} 900
twitter_api_ratelimits_remaining{endpoint="/statuses/user_timeline",resource="statuses"} 1500
twitter_api_ratelimits_remaining{endpoint="/trends/available",resource="trends"} 75
twitter_api_ratelimits_remaining{endpoint="/trends/closest",resource="trends"} 75
twitter_api_ratelimits_remaining{endpoint="/trends/place",resource="trends"} 75
twitter_api_ratelimits_remaining{endpoint="/users/lookup",resource="users"} 300
twitter_api_ratelimits_remaining{endpoint="/users/profile_banner",resource="users"} 180
twitter_api_ratelimits_remaining{endpoint="/users/reverse_lookup",resource="users"} 15
twitter_api_ratelimits_remaining{endpoint="/users/show/:id",resource="users"} 900
twitter_api_ratelimits_remaining{endpoint="/users/suggestions",resource="users"} 15
twitter_api_ratelimits_remaining{endpoint="/users/suggestions/:slug",resource="users"} 15
twitter_api_ratelimits_remaining{endpoint="/users/suggestions/:slug/members",resource="users"} 15
```

## Known Issues

Current we only support Twitter app auth, not user auth. Please add an issue if you require user auth

## Contributions

To contribute please submit a PR and ensure you run `make qa` before submitting
