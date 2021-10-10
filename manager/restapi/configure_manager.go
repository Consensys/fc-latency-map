// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/robfig/cron"

	"github.com/ConsenSys/fc-latency-map/manager/config"
	"github.com/ConsenSys/fc-latency-map/manager/internal/handlers"
	"github.com/ConsenSys/fc-latency-map/manager/jobs"
	"github.com/ConsenSys/fc-latency-map/manager/restapi/operations"
	"github.com/ConsenSys/fc-latency-map/manager/restapi/operations/check"
)

//go:generate swagger generate server --target ../../manager --name Manager --spec ../swagger.yml --principal interface{}

var conf = config.NewConfig()
var scheduler = cron.New()

func configureFlags(_ *operations.ManagerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.ManagerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.CheckGetHealthCheckHandler = check.GetHealthCheckHandlerFunc(func(params check.GetHealthCheckParams) middleware.Responder {
		payload := handlers.GetHealthCheckHandler()
		return check.NewGetHealthCheckOK().WithPayload(&payload)
	})
	api.CheckGetMetricsHandler = check.GetMetricsHandlerFunc(func(params check.GetMetricsParams) middleware.Responder {
		payload := handlers.GetMetricsHandler()
		return check.NewGetMetricsOK().WithPayload(&payload)
	})

	schedule := conf.GetString("CRON_SCHEDULE")
	log.Printf("Scheduling GetMesures task: %s\n", schedule)
	scheduler.AddFunc(schedule, func() {
		log.Printf("GetMesures task started at %s\n", time.Now())
		jobs.RunTaskGetMeasures()
	})
	scheduler.Start()

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {
		scheduler.Stop()
	}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(_ *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(_ *http.Server, _, _ string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
