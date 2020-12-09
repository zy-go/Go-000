package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/firfly/taliatoolkits/errgroup"
)

func main() {
    g := errgroup.WithContext(context.Background())

    g.Go(func(ctx context.Context) error {
        c := make(chan os.Signal)
        signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
        select {
        case <-c:
            g.StopAll()
        case <-ctx.Done():
        }
        return nil
    })

    g.Go(func(ctx context.Context) error { return newServer(ctx, ":8080", g.StopAll) })
    g.Go(func(ctx context.Context) error { return newServer(ctx, ":8081", g.StopAll) })
    g.Go(func(ctx context.Context) error { return newServer(ctx, ":8082", g.StopAll) })

    if err := g.Wait(); err != nil {
        log.Fatal(err)
    }
}

func newServer(ctx context.Context, addr string, cb func()) error {
    s := &http.Server{Addr: addr}
    s.RegisterOnShutdown(cb)
    go func() {
        <-ctx.Done()
        c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        s.Shutdown(c)
    }()

    if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed{
        return err
    }
    return nil
}
