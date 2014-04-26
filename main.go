package main

import (
	"errors"
	"net/http"

	"github.com/pnelson/waitress"
)

type RootContext struct {
	*waitress.Context
}

type CreateShortcutForm struct {
	URL string `json:"url"`
}

type CreateShortcutResponse struct {
	Id       int64  `json:"id"`
	URL      string `json:"url"`
	ShortURL string `json:"short_url"`
}

var (
	ErrInvalidURL = errors.New("invalid url")
)

func (ctx *RootContext) Create() http.Handler {
	form := &CreateShortcutForm{}
	if err := ctx.DecodeJSON(form); err != nil {
		return ctx.Abort(500)
	}

	if len(form.URL) == 0 {
		return ctx.Abort(422)
	}

	shortcut := NewShortcut()
	shortcut.URL = form.URL
	shortcut.Save()

	builder := ctx.Build("GET", "Redirect")
	builder.Set("encoded", shortcut.Encode())

	url, ok := builder.Build()
	if !ok {
		return ctx.Abort(500)
	}

	response := &CreateShortcutResponse{
		Id:       shortcut.Id,
		URL:      shortcut.URL,
		ShortURL: url.String(),
	}

	b, err := ctx.EncodeJSON(response)
	if err != nil {
		return ctx.Abort(500)
	}

	ctx.Status(201)
	return ctx.WriteJSON(b)
}

func (ctx *RootContext) Redirect(encoded string) http.Handler {
	shortcut, err := FindShortcut(encoded)
	if err != nil {
		return ctx.Abort(404)
	}

	return ctx.RedirectToWithCode(shortcut.URL, 307)
}

func Application() *waitress.Application {
	app := waitress.New((*RootContext)(nil))
	app.Route(`/`, "Create", []string{"POST"})
	app.Route(`/<encoded>`, "Redirect", []string{"GET"})
	return app
}

func main() {
	app := Application()
	app.Run()
}
