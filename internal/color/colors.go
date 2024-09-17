package color

import "github.com/Jamlie/colors"

const (
	greenSuccess    = 46
	orangeWarn      = 208
	redError        = 196
	grayPlaceholder = 248
)

var (
	Success     = colors.NewCustomId(greenSuccess)
	Warn        = colors.NewCustomId(orangeWarn, colors.WithBold)
	Error       = colors.NewCustomId(redError, colors.WithBold)
	Placeholder = colors.NewCustomId(grayPlaceholder)
)
