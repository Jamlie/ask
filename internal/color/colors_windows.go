//go:build windows

package color

import "github.com/Jamlie/colors"

var (
	Success     = colors.New(colors.GreenFg)
	Warn        = colors.New(colors.YellowFg, colors.WithBold)
	Error       = colors.New(colors.RedFg, colors.WithBold)
	Placeholder = colors.New(colors.WhiteFg)
)
