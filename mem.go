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

// Package pprof4svc provides functionality for exposing Go runtime statistics via HTTP endpoints.
// This file implements the memory statistics endpoint for the pprof service.
package pprof4svc

import (
	"fmt"
	"math"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// mem0 handles HTTP requests to the memory statistics endpoint.
// It reads runtime.MemStats and returns either a formatted text response or JSON based on the "json" query parameter.
func mem0(ctx *gin.Context) {
	var ms runtime.MemStats
	// Read memory statistics from the runtime
	runtime.ReadMemStats(&ms)
	// Check the "json" query parameter (case-insensitive)
	json0 := strings.ToLower(ctx.Query("json"))
	switch json0 {
	default:
		// Return formatted text output for memory stats by default
		ctx.String(http.StatusOK, memStats(&ms))
	case "1", "t", "true":
		// Return JSON output for memory stats if json=1, t, or true
		ctx.JSON(http.StatusOK, memStatsJSON(&ms))
	}
}

// memStats formats runtime.MemStats into a human-readable string.
// It includes memory allocation, GC, and system memory stats, with memory sizes in human-readable units and time in milliseconds.
func memStats(ms *runtime.MemStats) string {
	// Initialize output with a header
	output := "=========================== Go Runtime Memory Statistics ===========================\n"

	// Add memory allocation statistics, formatted using convertBytes
	output += fmt.Sprintf("HeapAlloc:   %s (Current heap memory in use)\n", convertBytes(ms.HeapAlloc))
	output += fmt.Sprintf("TotalAlloc:  %s (Cumulative total memory allocated on heap)\n", convertBytes(ms.TotalAlloc))
	output += fmt.Sprintf("Sys:         %s (Total memory obtained from OS)\n", convertBytes(ms.Sys))
	output += fmt.Sprintf("HeapSys:     %s (Memory reserved for heap from OS)\n", convertBytes(ms.HeapSys))
	output += fmt.Sprintf("HeapIdle:    %s (Heap memory reserved but not in use)\n", convertBytes(ms.HeapIdle))
	output += fmt.Sprintf("HeapInuse:   %s (Heap memory currently in use)\n", convertBytes(ms.HeapInuse))
	output += fmt.Sprintf("HeapReleased:%s (Heap memory returned to OS)\n", convertBytes(ms.HeapReleased))

	// Add heap object and allocation counts
	output += fmt.Sprintf("HeapObjects: %d (Number of allocated heap objects)\n", ms.HeapObjects)
	output += fmt.Sprintf("Mallocs:     %d (Total number of mallocs)\n", ms.Mallocs)
	output += fmt.Sprintf("Frees:       %d (Total number of frees)\n", ms.Frees)

	// Add GC-related statistics
	output += fmt.Sprintf("NumGC:       %d (Number of garbage collections)\n", ms.NumGC)
	// Convert total GC pause time from nanoseconds to milliseconds
	output += fmt.Sprintf("PauseTotal:  %.2f ms (Total GC pause time)\n", float64(ms.PauseTotalNs)/1e6)
	// Format GC CPU fraction as a percentage
	output += fmt.Sprintf("GCCPUFraction: %.2f%% (Fraction of CPU used by GC)\n", ms.GCCPUFraction*100)

	// Include last GC time if available, formatted as YYYY-MM-DD HH:MM:SS
	if ms.LastGC > 0 {
		lastGCTime := time.Unix(0, int64(ms.LastGC)).Format("2006-01-02 15:04:05")
		output += fmt.Sprintf("LastGC:      %s (Time of last garbage collection)\n", lastGCTime)
	} else {
		output += "LastGC:      Not available\n"
	}

	// Add additional memory statistics for stack, mcache, mspan, and other system memory
	output += fmt.Sprintf("StackInuse:  %s (Memory used by stack)\n", convertBytes(ms.StackInuse))
	output += fmt.Sprintf("StackSys:    %s (Memory reserved for stack from OS)\n", convertBytes(ms.StackSys))
	output += fmt.Sprintf("MCacheInuse: %s (Memory used by mcache)\n", convertBytes(ms.MCacheInuse))
	output += fmt.Sprintf("MCacheSys:   %s (Memory reserved for mcache from OS)\n", convertBytes(ms.MCacheSys))
	output += fmt.Sprintf("MSpanInuse:  %s (Memory used by mspan)\n", convertBytes(ms.MSpanInuse))
	output += fmt.Sprintf("MSpanSys:    %s (Memory reserved for mspan from OS)\n", convertBytes(ms.MSpanSys))
	output += fmt.Sprintf("OtherSys:    %s (Other system memory)\n", convertBytes(ms.OtherSys))

	// Close output with a footer
	output += "=========================== Go Runtime Memory Statistics ==========================="
	return output
}

// memStatsJSON converts runtime.MemStats into a JSON-compatible map.
// It formats memory sizes in human-readable units, time in milliseconds, and timestamps as strings.
func memStatsJSON(ms *runtime.MemStats) map[string]any {
	// Set last GC time, defaulting to "Not available" if zero
	lastGCTime := "Not available"
	if ms.LastGC > 0 {
		lastGCTime = time.Unix(0, int64(ms.LastGC)).Format("2006-01-02 15:04:05")
	}
	// Return a map with memory statistics for JSON serialization
	return map[string]any{
		"HeapAlloc":     convertBytes(ms.HeapAlloc),
		"TotalAlloc":    convertBytes(ms.TotalAlloc),
		"Sys":           convertBytes(ms.Sys),
		"HeapSys":       convertBytes(ms.HeapSys),
		"HeapIdle":      convertBytes(ms.HeapIdle),
		"HeapInuse":     convertBytes(ms.HeapInuse),
		"HeapReleased":  convertBytes(ms.HeapReleased),
		"HeapObjects":   ms.HeapObjects,
		"Mallocs":       ms.Mallocs,
		"Frees":         ms.Frees,
		"NumGC":         ms.NumGC,
		"PauseTotalMs":  math.Round(float64(ms.PauseTotalNs)/1e6*100) / 100,
		"GCCPUFraction": math.Round(ms.GCCPUFraction*100*100) / 100,
		"LastGC":        lastGCTime,
		"StackInuse":    convertBytes(ms.StackInuse),
		"StackSys":      convertBytes(ms.StackSys),
		"MCacheInuse":   convertBytes(ms.MCacheInuse),
		"MCacheSys":     convertBytes(ms.MCacheSys),
		"MSpanInuse":    convertBytes(ms.MSpanInuse),
		"MSpanSys":      convertBytes(ms.MSpanSys),
		"OtherSys":      convertBytes(ms.OtherSys),
	}
}

// convertBytes converts a byte count to a human-readable string (B, KB, MB, GB).
// It rounds the value to two decimal places and selects the appropriate unit.
func convertBytes(bytes uint64) string {
	units := []string{"B", "KB", "MB", "GB"}
	value := float64(bytes)
	unitIndex := 0

	// Scale the value to the appropriate unit
	for value >= 1024 && unitIndex < len(units)-1 {
		value /= 1024
		unitIndex++
	}

	// Format the value with two decimal places and the selected unit
	return fmt.Sprintf("%.2f %s", math.Round(value*100)/100, units[unitIndex])
}
