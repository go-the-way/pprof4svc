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

package pprof4svc

import (
	"fmt"
	"net/http"
	"runtime/pprof"
	"strconv"

	"github.com/gin-gonic/gin"
)

func pprof0(ctx *gin.Context) {
	name := ctx.Param("name")
	ctx.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	p := pprof.Lookup(name)
	if p == nil {
		serveError(ctx.Writer, http.StatusNotFound, "Unknown profile")
		return
	}
	debug, _ := strconv.Atoi(ctx.Request.FormValue("debug"))
	if debug != 0 {
		ctx.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	} else {
		ctx.Writer.Header().Set("Content-Type", "application/octet-stream")
		ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, name))
	}
	p.WriteTo(ctx.Writer, debug)
}

func serveError(w http.ResponseWriter, status int, txt string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Go-Pprof", "1")
	w.Header().Del("Content-Disposition")
	w.WriteHeader(status)
	fmt.Fprintln(w, txt)
}
