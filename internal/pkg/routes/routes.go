package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"

	"github.com/fromanirh/kubevirt-template-indexer/pkg/templateindex"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var log logr.Logger
var index *templateindex.TemplateIndexer

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		log.Info(fmt.Sprintf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		))
	})
}

var routes = Routes{
	Route{
		"oses",
		"GET",
		"/oses",
		oses,
	},
	Route{
		"workloads",
		"GET",
		"/workloads",
		workloads,
	},
	Route{
		"sizes",
		"GET",
		"/sizes",
		sizes,
	},
	Route{
		"templates",
		"GET",
		"/templates",
		templates,
	},
}

func oses(w http.ResponseWriter, r *http.Request) {
	summarize("os", w, r)
}

func workloads(w http.ResponseWriter, r *http.Request) {
	summarize("workload", w, r)
}

func sizes(w http.ResponseWriter, r *http.Request) {
	summarize("size", w, r)
}

func filterOptions(r *http.Request) templateindex.FilterOptions {
	query := r.URL.Query()
	opts := templateindex.FilterOptions{}
	for _, param := range []string{"os", "workload", "size"} {
		if value := query.Get(param); value != "" {
			opts[param] = value
		}
	}
	return opts
}

func templates(w http.ResponseWriter, r *http.Request) {

	descriptions, err := index.DescribeBy(filterOptions(r))
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(descriptions)
	if err != nil {
		panic(err)
	}
}

func summarize(label string, w http.ResponseWriter, r *http.Request) {
	summaries, err := index.SummarizeBy(label)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(summaries)
	if err != nil {
		panic(err)
	}
}

func Serve(host string, port int, index_ *templateindex.TemplateIndexer, log_ logr.Logger) error {
	index = index_
	log = log_

	return http.ListenAndServe(
		fmt.Sprintf("%s:%d", host, port),
		NewRouter(),
	)
}
