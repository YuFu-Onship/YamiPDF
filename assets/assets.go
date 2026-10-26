package assets

import _ "embed"

//go:embed font/icon_outline.ttf
var FontIconOutline []byte

//go:embed font/icon_fill.ttf
var FontIconFill []byte

//go:embed font/emoji.ttf
var FontEmoji []byte

//go:embed font/sans.ttf
var FontSans []byte

// 图标索引
var (
	Icon_home          = "\uE2C2"
	Icon_moon          = "\uE330"
	Icon_books         = "\uE758"
	Icon_setting       = "\uE272"
	Icon_down          = "\uE136"
	Icon_up            = "\uE13C"
	Icon_palette       = "\uE6C8"
	Icon_language      = "\uE4A2"
	Icon_info          = "\uE2CE"
	Icon_check         = "\uE182"
	Icon_caret_down    = "\uE136"
	Icon_caret_up      = "\uE13C"
	Icon_vertical_line = "\uED70"
	Icon_experimental  = "\uE79E"
	Icon_loading       = "\uE66C"
	Icon_check_square  = "\uE186"
	Icon_square        = "\uE45E" // 可用作 max
	Icon_min           = "\uE32A"
	Icon_cards         = "\uE0F8" // 可用作 un_max
	Icon_x             = "\uE4F6" // 可用作 close
)
