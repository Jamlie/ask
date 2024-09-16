package logger

import "github.com/Jamlie/colors"

var (
	Success = colors.NewCustomId(46)
	Warn    = colors.NewCustomId(208, colors.WithBold)
	Error   = colors.NewCustomId(196, colors.WithBold)
)
