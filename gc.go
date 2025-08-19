// Copyright 2025 pprof4svc Author. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package pprof4svc provides functionality for exposing Go runtime statistics via HTTP endpoints.
// This file implements the GC statistics endpoint for the pprof service.
package pprof4svc

import (
	"fmt"
	"math"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
)

// gc0 handles HTTP requests to the GC statistics endpoint.
// It reads debug.GCStats and returns either a formatted text response or JSON based on the "json" query parameter.
func gc0(ctx *gin.Context) {
	var gs debug.GCStats
	// Read GC statistics from the runtime
	debug.ReadGCStats(&gs)
	// Check the "json" query parameter (case-insensitive)
	json0 := strings.ToLower(ctx.Query("json"))
	switch json0 {
	default:
		// Return formatted text output for GC stats by default
		ctx.String(http.StatusOK, gcStats(&gs))
	case "1", "t", "true":
		// Return JSON output for GC stats if json=1, t, or true
		ctx.JSON(http.StatusOK, gcStatsJSON(&gs))
	}
}

// gcStats formats debug.GCStats into a human-readable string.
// It includes the number of GC runs, total pause time, last GC time, and recent pause details, with time in milliseconds.
func gcStats(gs *debug.GCStats) string {
	// Initialize output with a header
	output := "=========================== Go Runtime GC Statistics ===========================\n"

	// Add number of garbage collections
	output += fmt.Sprintf("NumGC:       %d (Number of garbage collections)\n", gs.NumGC)
	// Convert total pause time from nanoseconds to milliseconds
	output += fmt.Sprintf("PauseTotal:  %.2f ms (Total GC pause time)\n", float64(gs.PauseTotal.Nanoseconds())/1e6)

	// Include last GC time if available, formatted as YYYY-MM-DD HH:MM:SS
	if !gs.LastGC.IsZero() {
		output += fmt.Sprintf("LastGC:      %s (Time of last garbage collection)\n", gs.LastGC.Format("2006-01-02 15:04:05"))
	} else {
		output += "LastGC:      Not available\n"
	}

	// List recent GC pause durations in milliseconds
	output += fmt.Sprintf("Recent Pauses (%d recorded):\n", len(gs.Pause))
	for i, pause := range gs.Pause {
		output += fmt.Sprintf("  Pause %d:   %.2f ms\n", i+1, float64(pause.Nanoseconds())/1e6)
	}

	// List recent GC pause end times, formatted as YYYY-MM-DD HH:MM:SS
	output += fmt.Sprintf("Recent Pause Ends (%d recorded):\n", len(gs.PauseEnd))
	for i, end := range gs.PauseEnd {
		output += fmt.Sprintf("  Pause End %d: %s\n", i+1, end.Format("2006-01-02 15:04:05"))
	}

	// Close output with a footer
	output += "=========================== Go Runtime GC Statistics ===========================\n"
	return output
}

// gcStatsJSON converts debug.GCStats into a JSON-compatible map.
// It formats pause times in milliseconds and timestamps as strings, suitable for JSON output.
func gcStatsJSON(gs *debug.GCStats) map[string]any {
	// Convert recent pause durations to strings in milliseconds
	recentPausesMs := make([]string, len(gs.Pause))
	for i, pause := range gs.Pause {
		recentPausesMs[i] = fmt.Sprintf("%.2f", float64(pause.Nanoseconds())/1e6)
	}

	// Format recent pause end times as YYYY-MM-DD HH:MM:SS
	recentPauseEnds := make([]string, len(gs.PauseEnd))
	for i, end := range gs.PauseEnd {
		recentPauseEnds[i] = end.Format("2006-01-02 15:04:05")
	}

	// Set last GC time, defaulting to "Not available" if zero
	lastGCTime := "Not available"
	if !gs.LastGC.IsZero() {
		lastGCTime = gs.LastGC.Format("2006-01-02 15:04:05")
	}

	// Return a map with GC statistics for JSON serialization
	return map[string]any{
		"NumGC":           gs.NumGC,
		"PauseTotalMs":    math.Round(float64(gs.PauseTotal.Nanoseconds())/1e6*100) / 100,
		"LastGC":          lastGCTime,
		"RecentPausesMs":  recentPausesMs,
		"RecentPauseEnds": recentPauseEnds,
	}
}
