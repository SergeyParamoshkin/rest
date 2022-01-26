//
// REST
// ====
// This example demonstrates a HTTP REST web service with some fixture data.
// Follow along the example and patterns.
//
// Also check routes.json for the generated docs from passing the -routes flag,
// to run yourself do: `go run . -routes`
//
// Boot the server:
// ----------------
// $ go run main.go
//
// Client requests:
// ----------------
// $ curl http://localhost:3333/
// root.
//
// $ curl http://localhost:3333/articles
// [{"id":"1","title":"Hi"},{"id":"2","title":"sup"}]
//
// $ curl http://localhost:3333/articles/1
// {"id":"1","title":"Hi"}
//
// $ curl -X DELETE http://localhost:3333/articles/1
// {"id":"1","title":"Hi"}
//
// $ curl http://localhost:3333/articles/1
// "Not Found"
//
// $ curl -X POST -d '{"id":"will-be-omitted","title":"awesomeness"}' http://localhost:3333/articles
// {"id":"97","title":"awesomeness"}
//
// $ curl http://localhost:3333/articles/97
// {"id":"97","title":"awesomeness"}
//
// $ curl http://localhost:3333/articles
// [{"id":"2","title":"sup"},{"id":"97","title":"awesomeness"}]
//
package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/SergeyParamoshkin/rest/internal/article"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/docgen"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const ServiceName = "rest"

type CtxKey int8

const (
	// CtxKeyLogger is context key.
	CtxKeyLogger CtxKey = iota
)

type App struct {
	sugarLogger *zap.SugaredLogger
	config      Config
}

// nolint
func main() {

	// nolint
	var (
		routes   = flag.Bool("routes", getEnvBool(ServiceName+"_routes", false), "Generate router documentation")
		addr     = flag.String("addr", getEnv(ServiceName+"_ADDR", ":3333"), "application port")
		diagPort = flag.String("diagaddr", getEnv(ServiceName+"_DIAG_ADDR", ":9999"), "diag port")
		logLevel = flag.Int64("loglevel", getEnvInt64(ServiceName+"LOG_LEVEL", -1),
			`
		-1 - Debug
		 0 - Info 
		 1 - Warn
		 2 - Error
		 3 - DPanic
		 4 - Panic
		 5 - Fatal
		`)
	)

	flag.Parse()

	zpc := zap.NewProductionConfig()

	zpc.Level.SetLevel(zapcore.Level(*logLevel))
	logger, err := zpc.Build()
	if err != nil {
		panic(err)
	}

	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	a := App{
		sugarLogger: sugar,
	}

	config := prometheus.Config{}
	c := controller.New(
		processor.New(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			export.CumulativeExportKindSelector(),
			processor.WithMemory(true),
		),
	)
	exporter, err := prometheus.New(config, c)
	if err != nil {
		a.sugarLogger.Panicf("failed to initialize prometheus exporter %v", err)
	}
	global.SetMeterProvider(exporter.MeterProvider())

	meter := global.Meter(ServiceName)

	ClientCompletedCount := metric.Must(meter).NewInt64Counter(
		"http/client/completed_count",
		metric.WithDescription("Count of completed requests, by HTTP method and response status"),
	).Bind([]attribute.KeyValue{attribute.String("status", "200")}...)
	defer ClientCompletedCount.Unbind()

	// observerLock := new(sync.RWMutex)
	// observerValueToReport := new(float64)
	// observerLabelsToReport := new([]attribute.KeyValue)
	// cb := func(_ context.Context, result metric.Float64ObserverResult) {
	// 	(*observerLock).RLock()
	// 	value := *observerValueToReport
	// 	labels := *observerLabelsToReport
	// 	(*observerLock).RUnlock()
	// 	result.Observe(value, labels...)
	// }
	// _ = metric.Must(meter).NewFloat64ValueObserver("ex.com.one", cb,
	// 	metric.WithDescription("A ValueObserver set to 1.0"),
	// )

	// valuerecorder := metric.Must(meter).NewFloat64ValueRecorder("ex.com.two")
	// counter := metric.Must(meter).NewFloat64Counter("ex.com.three")

	// commonLabels := []attribute.KeyValue{lemonsKey.Int(10), attribute.String("A", "1"), attribute.String("B", "2"), attribute.String("C", "3")}
	// notSoCommonLabels := []attribute.KeyValue{lemonsKey.Int(13)}

	// ctx := context.Background()

	// (*observerLock).Lock()
	// *observerValueToReport = 1.0
	// *observerLabelsToReport = commonLabels
	// (*observerLock).Unlock()
	// meter.RecordBatch(
	// 	ctx,
	// 	commonLabels,
	// 	valuerecorder.Measurement(2.0),
	// 	counter.Measurement(12.0),
	// )

	// time.Sleep(5 * time.Second)

	// (*observerLock).Lock()
	// *observerValueToReport = 1.0
	// *observerLabelsToReport = notSoCommonLabels
	// (*observerLock).Unlock()
	// meter.RecordBatch(
	// 	ctx,
	// 	notSoCommonLabels,
	// 	valuerecorder.Measurement(2.0),
	// 	counter.Measurement(22.0),
	// )

	// time.Sleep(5 * time.Second)

	// (*observerLock).Lock()
	// *observerValueToReport = 13.0
	// *observerLabelsToReport = commonLabels
	// (*observerLock).Unlock()
	// meter.RecordBatch(
	// 	ctx,
	// 	commonLabels,
	// 	valuerecorder.Measurement(12.0),
	// 	counter.Measurement(13.0),
	// )
	r := chi.NewRouter()

	diagRouter := chi.NewRouter()
	diagRouter.Get("/metrics", exporter.ServeHTTP)

	r.Use(middleware.RequestID)
	r.Use(a.Logger)
	// r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("root."))
		if err != nil {
			sugar.Errorw(err.Error())
		}
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		logger := r.Context().Value(CtxKeyLogger).(*zap.SugaredLogger)
		ClientCompletedCount.Add(r.Context(), 1)
		_, err := w.Write([]byte("pong"))
		logger.Infow("pong")
		if err != nil {
			logger.Errorw(err.Error())
		}
	})

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		sugar.Panicw("panic")
	})

	// RESTy routes for "articles" resource
	r.Route("/articles", func(r chi.Router) {
		r.With(paginate).Get("/", article.ListArticles)
		r.Post("/", article.CreateArticle)       // POST /articles
		r.Get("/search", article.SearchArticles) // GET /articles/search

		r.Route("/{articleID}", func(r chi.Router) {
			r.Use(article.ArticleCtx)            // Load the *Article on the request context
			r.Get("/", article.GetArticle)       // GET /articles/123
			r.Put("/", article.UpdateArticle)    // PUT /articles/123
			r.Delete("/", article.DeleteArticle) // DELETE /articles/123
		})

		// GET /articles/whats-up
		r.With(article.ArticleCtx).Get("/{articleSlug:[a-z-]+}", article.GetArticle)
	})

	// Mount the admin sub-router, which btw is the same as:
	// r.Route("/admin", func(r chi.Router) { admin routes here })
	r.Mount("/admin", adminRouter())

	// Passing -routes to the program will generate docs for the above
	// router definition. See the `routes.json` file in this folder for
	// the output.
	if *routes {
		// fmt.Println(docgen.JSONRoutesDoc(r))
		// nolint
		fmt.Println(docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{
			ProjectPath: "github.com/go-chi/chi/v5",
			Intro:       "Welcome to the chi/_examples/rest generated docs.",
		}))

		return
	}

	FileServer(r, "/swagger-ui", Swagger())

	go func() {
		err = http.ListenAndServe(*addr, r)
		if err != nil {
			a.sugarLogger.Errorw(err.Error())
		}
	}()

	err = http.ListenAndServe(*diagPort, diagRouter)
	if err != nil {
		a.sugarLogger.Errorw(err.Error())
	}

}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}

	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func (a *App) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.sugarLogger.Debug(next, r)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxKeyLogger, a.sugarLogger)))
	})
}

//go:embed swagger-ui
var embededFiles embed.FS

func Swagger() http.FileSystem {
	log.Print("using embed mode")
	fsys, err := fs.Sub(embededFiles, "swagger-ui")

	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

// A completely separate router for administrator routes
func adminRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(AdminOnly)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("admin: index"))
		if err != nil {
			log.Println(err)
		}
	})
	r.Get("/accounts", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("admin: list accounts.."))
		if err != nil {
			log.Println(err)
		}
	})
	r.Get("/users/{userId}", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(fmt.Sprintf("admin: view user id %v", chi.URLParam(r, "userId"))))
		if err != nil {
			log.Println(err)
		}
	})

	return r
}

// AdminOnly middleware restricts access to just administrators.
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := r.Context().Value("acl.admin").(bool)
		if !ok || !isAdmin {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)

			return
		}
		next.ServeHTTP(w, r)
	})
}

// paginate is a stub, but very possible to implement middleware logic
// to handle the request params for handling a paginated request.
func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// just a stub.. some ideas are to look at URL query params for something like
		// the page number, or the limit, and send a query cursor down the chain
		next.ServeHTTP(w, r)
	})
}

// This is entirely optional, but I wanted to demonstrate how you could easily
// add your own logic to the render.Respond method.
// nolint
func init() {
	render.Respond = func(w http.ResponseWriter, r *http.Request, v interface{}) {
		if err, ok := v.(error); ok {

			// We set a default error status response code if one hasn't been set.
			if _, ok := r.Context().Value(render.StatusCtxKey).(int); !ok {
				w.WriteHeader(400)
			}

			// We log the error
			// nolint
			fmt.Printf("Logging err: %s\n", err.Error())

			// We change the response to not reveal the actual error message,
			// instead we can transform the message something more friendly or mapped
			// to some code / language, etc.
			render.DefaultResponder(w, r, render.M{"status": "error"})

			return
		}

		render.DefaultResponder(w, r, v)
	}
}
