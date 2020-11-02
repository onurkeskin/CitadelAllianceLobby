package domain

import (
	"net/http"

	"golang.org/x/net/context"
)

type ContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

func (h ContextHandlerFunc) ServeHTTP(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	h(ctx, rw, r)
}

type MiddlewareFunc func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (m MiddlewareFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	m(rw, r, next)
}

type ContextMiddlewareFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (m ContextMiddlewareFunc) ServeHTTP(ctx context.Context, rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	m(ctx, rw, r, next)
}

type IMiddleware interface {
	Handler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type ContextMiddleware interface {
	Handler(ctx context.Context, rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}
