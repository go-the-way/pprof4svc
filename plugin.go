// Copyright 2025 pprof4svc Author. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package pprof4svc provides a Gin-based plugin for integrating pprof, memory, GC, and trace endpoints
// into an HTTP service with token-based authentication and randomized route prefixes for security.
package pprof4svc

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
)

// Constants defining the routes for pprof, memory, GC, and trace endpoints.
const (
	pprofIndexRoute   = "/debug/pprof/"        // Base route for pprof index
	pprofNameRoute    = "/debug/pprof/:name"   // Route for specific pprof profiles (e.g., heap, goroutine)
	pprofCmdlineRoute = "/debug/pprof/cmdline" // Route for command line arguments
	pprofProfileRoute = "/debug/pprof/profile" // Route for CPU profile
	pprofSymbolRoute  = "/debug/pprof/symbol"  // Route for symbol lookup
	pprofTraceRoute   = "/debug/pprof/trace"   // Route for execution trace
	memRoute          = "/debug/mem"           // Route for memory statistics
	gcRoute           = "/debug/gc"            // Route for GC statistics
	traceRoute        = "/debug/trace"         // Route for trace control
)

// plugin represents the configuration for the pprof service plugin.
// It holds the entrypoint, authentication token, and prefixed routes for various endpoints.
type plugin struct {
	entrypoint   string // Main entrypoint for accessing the pprof service
	token        string // Token for authenticating access to the pprof service
	prefix       string // Random prefix for securing routes
	pprofIndex   string // Prefixed route for pprof index
	pprofName    string // Prefixed route for specific pprof profiles
	pprofCmdline string // Prefixed route for command line arguments
	pprofProfile string // Prefixed route for CPU profile
	pprofSymbol  string // Prefixed route for symbol lookup
	pprofTrace   string // Prefixed route for execution trace
	memRoute     string // Prefixed route for memory statistics
	gcRoute      string // Prefixed route for GC statistics
	traceRoute   string // Prefixed route for trace control
}

// DefaultPlugin creates a plugin with the default pprof entrypoint and the provided token.
// It uses the default pprof index route as the entrypoint.
func DefaultPlugin(token string) *plugin {
	return Plugin(pprofIndexRoute, token)
}

// Plugin creates a new plugin instance with the specified entrypoint and token.
// It generates a random prefix for securing routes and initializes all routes with the prefix.
func Plugin(entrypoint, token string) *plugin {
	prefix := randPrefix()
	return &plugin{
		entrypoint:   entrypoint,
		token:        token,
		prefix:       prefix,
		pprofIndex:   prefix + pprofIndexRoute,
		pprofName:    prefix + pprofNameRoute,
		pprofCmdline: prefix + pprofCmdlineRoute,
		pprofProfile: prefix + pprofProfileRoute,
		pprofSymbol:  prefix + pprofSymbolRoute,
		pprofTrace:   prefix + pprofTraceRoute,
		memRoute:     prefix + memRoute,
		gcRoute:      prefix + gcRoute,
		traceRoute:   prefix + traceRoute,
	}
}

// Plug registers the plugin's routes with the provided Gin engine.
// It sets up handlers for pprof, memory, GC, and trace endpoints, with token-based authentication.
func (p *plugin) Plug(engine *gin.Engine) {
	// Define type aliases for handler functions to simplify wrapping
	type (
		fn    = http.HandlerFunc       // Standard HTTP handler function
		fnCtx = func(ctx *gin.Context) // Gin context handler function
	)
	// wrapped converts a standard HTTP handler to a Gin handler
	wrapped := func(f fn) fnCtx {
		return func(ctx *gin.Context) {
			f(ctx.Writer, ctx.Request)
		}
	}
	// Register routes with the Gin engine
	engine.GET(p.entrypoint, p.handler)
	engine.GET(p.pprofIndex, wrapped(pprof.Index))
	engine.GET(p.pprofCmdline, wrapped(pprof.Cmdline))
	engine.GET(p.pprofProfile, wrapped(pprof.Profile))
	engine.GET(p.pprofSymbol, wrapped(pprof.Symbol))
	engine.GET(p.pprofTrace, wrapped(pprof.Trace))
	engine.GET(p.pprofName, pprof0)
	engine.GET(p.memRoute, mem0)
	engine.GET(p.gcRoute, gc0)
	engine.GET(p.traceRoute, trace0)
}

// handler authenticates requests to the entrypoint using a token query parameter.
// If the token is valid, it redirects to the pprof index route; otherwise, it returns an unauthorized error.
func (p *plugin) handler(ctx *gin.Context) {
	token0, _ := ctx.GetQuery("token")
	if token0 != p.token {
		ctx.String(http.StatusUnauthorized, "Unauthorized")
		return
	}
	ctx.Redirect(http.StatusMovedPermanently, p.pprofIndex)
}

// randPrefix generates a random string prefix for securing routes.
// The prefix is 40 characters long, using alphanumeric characters and underscores.
func randPrefix() string {
	rand0 := func(len0 int) (str string) {
		chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_"
		for i := 0; i < len0; i++ {
			str += fmt.Sprintf("%c", chars[rand.Intn(len(chars))])
		}
		return
	}
	return "/" + rand0(40)
}
