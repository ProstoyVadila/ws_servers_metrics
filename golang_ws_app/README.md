## TODO

1. Implement valid models (custom marshal/unmarshal for snake -> camel cases)
2. Add metrics
3. Rewrite docker file for two binaries.
4. Add golang app to k6/scripts

### Metrics implementation example

```go
func routerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		REQUEST_INPROGRESS.Inc()

		startTime := time.Now()
		route := mux.CurrentRoute(r)
		path, err := route.GetPathTemplate()
		if err != nil {
			log.Fatalln("cannot get path template")
		}

		next.ServeHTTP(w, r)

		takenTime := time.Since(startTime)
		REQUEST_RESPOND_TIME.WithLabelValues(path).Observe(takenTime.Seconds())
		REQUEST_RESPOND_TIME_HIST.WithLabelValues(path).Observe(takenTime.Seconds())
		REQUEST_INPROGRESS.Dec()
		REQUEST_COUNT.Inc()
	})
}
```
