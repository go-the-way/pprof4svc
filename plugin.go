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
	"math/rand"
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
)

const (
	pprofIndexRoute   = "/debug/pprof/"
	pprofNameRoute    = "/debug/pprof/:name"
	pprofCmdlineRoute = "/debug/pprof/cmdline"
	pprofProfileRoute = "/debug/pprof/profile"
	pprofSymbolRoute  = "/debug/pprof/symbol"
	pprofTraceRoute   = "/debug/pprof/trace"
)

type plugin struct {
	entrypoint, token string
	prefix            string
	pprofIndex        string
	pprofName         string
	pprofCmdline      string
	pprofProfile      string
	pprofSymbol       string
	pprofTrace        string
}

func DefaultPlugin(token string) *plugin { return Plugin(pprofIndexRoute, token) }

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
	}
}

func (p *plugin) Plug(engine *gin.Engine) {
	type (
		fn    = http.HandlerFunc
		fnCtx = func(ctx *gin.Context)
	)
	wrapped := func(f fn) fnCtx { return func(ctx *gin.Context) { f(ctx.Writer, ctx.Request) } }
	engine.GET(p.entrypoint, p.handler)
	engine.GET(p.pprofIndex, wrapped(pprof.Index))
	engine.GET(p.pprofCmdline, wrapped(pprof.Cmdline))
	engine.GET(p.pprofProfile, wrapped(pprof.Profile))
	engine.GET(p.pprofSymbol, wrapped(pprof.Symbol))
	engine.GET(p.pprofTrace, wrapped(pprof.Trace))
	engine.GET(p.pprofName, pprof0)
}

func (p *plugin) handler(ctx *gin.Context) {
	token0, _ := ctx.GetQuery("token")
	if token0 != p.token {
		ctx.String(http.StatusUnauthorized, "Unauthorized")
		return
	}
	ctx.Redirect(http.StatusMovedPermanently, p.pprofIndex)
}

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
