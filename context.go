package main

import (
	"context"

	"github.com/sirupsen/logrus"
)

// ContextKey is a key used to look up a context value
// https://golang.org/pkg/context/#WithValue
type ContextKey string

func setContextValue(ctx context.Context, k string, v interface{}) context.Context {
	logrus.Debugf("set context key %v to value %v", ContextKey(k), v)
	newCtx := context.WithValue(ctx, ContextKey(k), v)
	return newCtx
}

func getContextValue(ctx context.Context, k string) interface{} {
	if v := ctx.Value(ContextKey(k)); v != nil {
		logrus.Debugf("found context value:", v)
		return v
	}
	return nil
}
