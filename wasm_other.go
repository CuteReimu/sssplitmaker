//go:build !windows

package main

import "errors"

func (a *App) FixLiveSplit() (string, error) {
	return "", errors.New("仅支持 Windows")
}
