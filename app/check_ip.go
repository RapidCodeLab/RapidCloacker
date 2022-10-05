package main

import (
	"github.com/valyala/fasthttp"
)

func (app *application) checkIP(ctx *fasthttp.RequestCtx) {

	if !ctx.IsGet() {
		ctx.Response.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		return
	}

	ip := ctx.UserValue("ip").(string)

	if !app.validateIP(ip) {
		ctx.Response.SetStatusCode(fasthttp.StatusNoContent)
		return
	}

	ctx.Response.SetStatusCode(fasthttp.StatusOK)

}
