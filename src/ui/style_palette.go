package ui

import "image/color"

// PICO8 色板
var (
	COLOR_BLACK       = color.NRGBA{R: 0, G: 0, B: 0, A: 255}       // 0: 黑色
	COLOR_DARK_BLUE   = color.NRGBA{R: 29, G: 43, B: 83, A: 255}    // 1: 深蓝
	COLOR_DARK_PURPLE = color.NRGBA{R: 126, G: 37, B: 83, A: 255}   // 2: 深紫
	COLOR_DARK_GREEN  = color.NRGBA{R: 0, G: 135, B: 81, A: 255}    // 3: 深绿
	COLOR_BROWN       = color.NRGBA{R: 171, G: 82, B: 54, A: 255}   // 4: 棕色
	COLOR_DARK_GRAY   = color.NRGBA{R: 95, G: 87, B: 79, A: 255}    // 5: 深灰
	COLOR_LIGHT_GRAY  = color.NRGBA{R: 194, G: 195, B: 199, A: 255} // 6: 浅灰
	COLOR_WHITE       = color.NRGBA{R: 255, G: 241, B: 232, A: 255} // 7: 白色
	COLOR_RED         = color.NRGBA{R: 255, G: 0, B: 77, A: 255}    // 8: 红色
	COLOR_ORANGE      = color.NRGBA{R: 255, G: 163, B: 0, A: 255}   // 9: 橙色
	COLOR_YELLOW      = color.NRGBA{R: 255, G: 236, B: 39, A: 255}  // 10: 黄色
	COLOR_GREEN       = color.NRGBA{R: 0, G: 228, B: 54, A: 255}    // 11: 绿色
	COLOR_BLUE        = color.NRGBA{R: 41, G: 173, B: 255, A: 255}  // 12: 蓝色
	COLOR_INDIGO      = color.NRGBA{R: 131, G: 118, B: 156, A: 255} // 13: 靛蓝/灰紫
	COLOR_PINK        = color.NRGBA{R: 255, G: 119, B: 168, A: 255} // 14: 粉红
	COLOR_PEACH       = color.NRGBA{R: 255, G: 204, B: 170, A: 255} // 15: 桃红/肤色
)

var Deco_palette_light DecoColorPalette = DecoColorPalette{
	Deco_fg:                 color.NRGBA{R: 26, G: 26, B: 26, A: 255},    // 主文本/图标 (深灰 #1A1A1A)
	Deco_bg:                 color.NRGBA{R: 243, G: 243, B: 243, A: 255}, // 背景色 (#F3F3F3)
	Decoicon_fg_idle:        color.NRGBA{R: 26, G: 26, B: 26, A: 255},    // 图标常态
	Decoicon_fg_hover:       color.NRGBA{R: 0, G: 0, B: 0, A: 255},       // 图标悬停 (纯黑加深)
	Decoicon_fg_press:       color.NRGBA{R: 0, G: 0, B: 0, A: 180},       // 图标按下
	Decoicon_bg_idle:        color.NRGBA{R: 0, G: 0, B: 0, A: 0},         // 按钮背景常态 (透明)
	Decoicon_bg_hover:       color.NRGBA{R: 0, G: 0, B: 0, A: 38},        // 按钮背景悬停 (约 15% 黑色半透明)
	Decoicon_bg_press:       color.NRGBA{R: 0, G: 0, B: 0, A: 65},        // 按钮背景按下 (约 25% 黑色半透明)
	Decoicon_fg_idle_close:  color.NRGBA{R: 26, G: 26, B: 26, A: 255},    // 关闭图标常态
	Decoicon_fg_hover_close: color.NRGBA{R: 255, G: 255, B: 255, A: 255}, // 关闭图标悬停 (反白)
	Decoicon_fg_press_close: color.NRGBA{R: 255, G: 255, B: 255, A: 200}, // 关闭图标按下
	Decoicon_bg_idle_close:  color.NRGBA{R: 0, G: 0, B: 0, A: 0},         // 关闭背景常态 (透明)
	Decoicon_bg_hover_close: color.NRGBA{R: 232, G: 17, B: 35, A: 255},   // 关闭背景悬停 (鲜艳红 #E81123)
	Decoicon_bg_press_close: color.NRGBA{R: 196, G: 43, B: 28, A: 255},   // 关闭背景按下 (深红 #C42B1C)
}

var Deco_palette_dark DecoColorPalette = DecoColorPalette{
	Deco_fg:                 color.NRGBA{R: 255, G: 255, B: 255, A: 255}, // 主文本/图标 (纯白)
	Deco_bg:                 color.NRGBA{R: 32, G: 32, B: 32, A: 255},    // 背景色 (#202020)
	Decoicon_fg_idle:        color.NRGBA{R: 255, G: 255, B: 255, A: 255}, // 图标常态
	Decoicon_fg_hover:       color.NRGBA{R: 255, G: 255, B: 255, A: 255}, // 图标悬停
	Decoicon_fg_press:       color.NRGBA{R: 255, G: 255, B: 255, A: 180}, // 图标按下
	Decoicon_bg_idle:        color.NRGBA{R: 0, G: 0, B: 0, A: 0},         // 按钮背景常态 (透明)
	Decoicon_bg_hover:       color.NRGBA{R: 255, G: 255, B: 255, A: 51},  // 按钮背景悬停 (20% 白色半透明)
	Decoicon_bg_press:       color.NRGBA{R: 255, G: 255, B: 255, A: 85},  // 按钮背景按下 (约 33% 白色半透明)
	Decoicon_fg_idle_close:  color.NRGBA{R: 255, G: 255, B: 255, A: 255}, // 关闭图标常态
	Decoicon_fg_hover_close: color.NRGBA{R: 255, G: 255, B: 255, A: 255}, // 关闭图标悬停
	Decoicon_fg_press_close: color.NRGBA{R: 255, G: 255, B: 255, A: 200}, // 关闭图标按下
	Decoicon_bg_idle_close:  color.NRGBA{R: 0, G: 0, B: 0, A: 0},         // 关闭背景常态 (透明)
	Decoicon_bg_hover_close: color.NRGBA{R: 232, G: 17, B: 35, A: 255},   // 关闭背景悬停 (鲜艳红 #E81123)
	Decoicon_bg_press_close: color.NRGBA{R: 196, G: 43, B: 28, A: 255},   // 关闭背景按下 (深红 #C42B1C)
}

var Palette_windows StylePalette = StylePalette{
	Light: ColorPalette{
		Fg_1: color.NRGBA{R: 0x1A, G: 0x1A, B: 0x1A, A: 255},
		Bg_1: color.NRGBA{R: 0xF3, G: 0xF3, B: 0xF3, A: 255},
		Fg_2: color.NRGBA{R: 0x61, G: 0x61, B: 0x61, A: 255},
		Bg_2: color.NRGBA{R: 0xE5, G: 0xE5, B: 0xE5, A: 255},
		Fg_3: color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 255},
		Bg_3: color.NRGBA{R: 0x00, G: 0x5A, B: 0x9E, A: 255},
	},
	Dark: ColorPalette{
		Fg_1: color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 255},
		Bg_1: color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 255},
		Fg_2: color.NRGBA{R: 0x9E, G: 0x9E, B: 0x9E, A: 255},
		Bg_2: color.NRGBA{R: 0x2D, G: 0x2D, B: 0x2D, A: 255},
		Fg_3: color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 255},
		Bg_3: color.NRGBA{R: 0x00, G: 0x78, B: 0xD4, A: 255},
	},
}

var Palette_teto StylePalette = StylePalette{
	Light: ColorPalette{
		Fg_1: color.NRGBA{R: 0x3C, G: 0x3C, B: 0x41, A: 255},
		Bg_1: color.NRGBA{R: 0xF3, G: 0xF3, B: 0xF3, A: 255},
		Fg_2: color.NRGBA{R: 0x32, G: 0x32, B: 0x37, A: 255},
		Bg_2: color.NRGBA{R: 0xE6, G: 0xE6, B: 0xEB, A: 255},
		Fg_3: color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 255},
		Bg_3: color.NRGBA{R: 0xFF, G: 0x2D, B: 0x58, A: 255},
	},
	Dark: ColorPalette{
		Fg_1: color.NRGBA{R: 0xFF, G: 0xD2, B: 0xDC, A: 255},
		Bg_1: color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 255},
		Fg_2: color.NRGBA{R: 0xF0, G: 0xC8, B: 0xD2, A: 255},
		Bg_2: color.NRGBA{R: 0x3C, G: 0x37, B: 0x37, A: 255},
		Fg_3: color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 255},
		Bg_3: color.NRGBA{R: 0xBE, G: 0x1E, B: 0x3C, A: 255},
	},
}

var Palette_koishi StylePalette = StylePalette{
	Light: ColorPalette{
		Fg_1: color.NRGBA{R: 0x2E, G: 0x3D, B: 0x30, A: 255},
		Bg_1: color.NRGBA{R: 0xF3, G: 0xF3, B: 0xF3, A: 255},
		Fg_2: color.NRGBA{R: 0x28, G: 0x46, B: 0x3c, A: 255},
		Bg_2: color.NRGBA{R: 0xfa, G: 0xf0, B: 0xbe, A: 255},
		Fg_3: color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 255},
		Bg_3: color.NRGBA{R: 0x46, G: 0xbe, B: 0x96, A: 255},
	},
	Dark: ColorPalette{
		Fg_1: color.NRGBA{R: 0xd2, G: 0xF0, B: 0xdc, A: 255},
		Bg_1: color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 255},
		Fg_2: color.NRGBA{R: 0xbc, G: 0xD7, B: 0xc5, A: 255},
		Bg_2: color.NRGBA{R: 0x32, G: 0x3c, B: 0x37, A: 255},
		Fg_3: color.NRGBA{R: 0xF4, G: 0xF4, B: 0xF4, A: 255},
		Bg_3: color.NRGBA{R: 0x32, G: 0xa0, B: 0x82, A: 255},
	},
}
