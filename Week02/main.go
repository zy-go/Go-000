package main

import (
    "log"

    "github.com/pkg/errors"
)

var ErrNoRows = errors.New("no found")

func dao() error {
    err := ErrNoRows
    if err != nil {
        return err
    }

    // do something
    return nil
}

func service() error {
    return dao()
}

func biz() {
    if err := service(); err != nil {
        log.Printf("%+v", err)
        return
    }

    // do something
    return
}

func main() {
    biz()
}
