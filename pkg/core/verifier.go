package core

import (
    r "golan/pkg/route"
)

func Verify(targets []r.Target) error {
    for _, target := range targets {
        _, err := target.Verify(r.RetrieveCurrentRoutesV4())
        if err != nil {
            return err
        }
    }
    return nil
}
