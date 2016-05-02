package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type RouteFunc func(db *sql.DB, params map[string]string, w http.ResponseWriter, r *http.Request)

type Route struct {
	Path    string
	Handler RouteFunc
}

type CompiledRoute struct {
	Path    string
	Handler RouteFunc
	Regex   *regexp.Regexp
	Params  []string
}

func CompileRoutes(routes []Route) []CompiledRoute {
	compiledRoutes := []CompiledRoute{}
	for _, route := range routes {
		originalRoute := route.Path
		var params []string
		for open := strings.Index(route.Path, "{"); open > 0; open = strings.Index(route.Path, "{") {
			close := strings.Index(route.Path, "}")
			if close < 0 {
				panic("malformed route")
			}
			param := route.Path[open+1 : close]

			params = append(params, param)
			route.Path = route.Path[:open] + `(.+)` + route.Path[close+1:]
		}
		regex := regexp.MustCompile(route.Path)
		compiledRoutes = append(compiledRoutes, CompiledRoute{Path: originalRoute, Handler: route.Handler, Regex: regex, Params: params})
	}
	return compiledRoutes
}

func Handler(db *sql.DB, routes []CompiledRoute, w http.ResponseWriter, r *http.Request) {
	requested := r.URL.Path[1:]

	for _, compiledRoute := range routes {
		match := compiledRoute.Regex.FindStringSubmatch(requested)
		if len(match) == 0 {
			continue
		}

		params := map[string]string{}
		for i, v := range compiledRoute.Params {
			params[v] = match[i+1]
		}
		compiledRoute.Handler(db, params, w, r)
		return
	}

	fmt.Fprintf(w, "You didn't ask for anything")
}
