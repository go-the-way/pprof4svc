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

// Package pprof4svc provides functionality for exposing Go runtime statistics and tracing via HTTP endpoints.
// This file implements the trace control endpoint for the pprof service.
package pprof4svc

import (
	"net/http"
	"runtime/trace"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// mu is a global mutex to ensure thread-safe access to trace operations.
var mu = &sync.Mutex{}

// trace handles HTTP requests to the trace control endpoint.
// It starts a runtime trace for a specified duration and writes the trace data to the HTTP response.
func trace0(ctx *gin.Context) {
	// Attempt to acquire the mutex and check if tracing is already active
	if !mu.TryLock() || trace.IsEnabled() {
		// Return an error if tracing is already active or the mutex cannot be acquired
		serveError(ctx.Writer, http.StatusBadRequest, "Tracing is already active")
		return
	}
	// Ensure the mutex is released after the function completes
	defer mu.Unlock()

	// Get the duration query parameter, defaulting to 10 seconds if not specified
	dur0str := ctx.Query("dur")
	if dur0str == "" {
		dur0str = "10s"
	}
	// Parse the duration string into a time.Duration
	dur0, _ := time.ParseDuration(dur0str)
	if dur0 == 0 {
		dur0 = time.Second * 10
	}
	// Start tracing, writing trace data to the HTTP response writer
	trace.Start(ctx.Writer)
	// Wait for the specified duration
	time.Sleep(dur0)
	// Stop tracing
	trace.Stop()
}
