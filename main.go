package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
)

const keyServerAddr = "serverAddr"

func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s: / request received\n", ctx.Value(keyServerAddr))
	io.WriteString(w, "Welcome to the root website")
}

func getHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s: /hello request received\n", ctx.Value(keyServerAddr))
	io.WriteString(w, "Hello HTTP!")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/hello", getHello)

	ctx, cancelCtx := context.WithCancel(context.Background())
	serverOne := &http.Server{
		Addr:    ":3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	serverTwo := &http.Server{
		Addr:    ":4444",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	go func() {
		err := serverOne.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server one closed")
		} else if err != nil {
			fmt.Printf("server one error: %s", err)
		}
		cancelCtx()
	}()

	go func() {
		err := serverTwo.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server two closed")
		} else if err != nil {
			fmt.Printf("server two error: %s", err)
		}
		cancelCtx()
	}()

	<-ctx.Done()
}
