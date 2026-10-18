package ui

type Language int

const (
	CN Language = iota
	EN
	JP
)

type Word [3]string

var TransTable = map[string]Word{
	// 侧边栏
	"home":       {CN: "主页", EN: "Home", JP: "ホーム"},
	"color_mode": {CN: "颜色模式", EN: "Color Mode", JP: "カラーモード"},
	"bookshelf":  {CN: "书架", EN: "Bookshelf", JP: "本棚"},
	"setting":    {CN: "设置", EN: "Settings", JP: "設定"},

	// 预览页
	"reset":  {CN: "重置", EN: "Reset", JP: "リセット"},
	"next":   {CN: "下一页", EN: "Next", JP: "次へ"},
	"last":   {CN: "上一页", EN: "Previous", JP: "前へ"},
	"rotate": {CN: "旋转", EN: "Rotate", JP: "回転"},
	"shader": {CN: "着色器", EN: "Shader", JP: "シェーダー"},
	"layout": {CN: "布局", EN: "Layout", JP: "レイアウト"},
	"record": {CN: "记录", EN: "Record", JP: "履歴"},

	// 设置页
	"info":         {CN: "信息", EN: "Info", JP: "情報"},
	"custom":       {CN: "个性化", EN: "Custom", JP: "カスタム"},
	"experimental": {CN: "实验性", EN: "Experimental", JP: "実験的"},
	"test":         {CN: "测试", EN: "Test", JP: "テスト"},

	// 设置_个性化
	"language":  {CN: "语言", EN: "Language", JP: "言語"},
	"palette":   {CN: "色板", EN: "Palette", JP: "パレット"},
	"dark_mode": {"暗色模式", "Dark Mode", ""},

	// 配色
	"koishi": {CN: "古明地恋", EN: "Koishi", JP: "古明地こいし"},

	// 设置_信息
	"version":   {"版本", "Version", "バージョン"},
	"developer": {"开发者", "Developer", "開発者"},
	"yufu":      {"雨符", "YUFU", "YUFU"},
	"website":   {"网站", "Website", "ウェブサイト"},

	// 其他
	"software_desc": {
		"轻量,美观,快速的PDF阅读器",
		"A lightweight, attractive, and fast PDF reader",
		"軽量で美しく、高速なPDFリーダー",
	},
}

func GetWord(id string, lang Language) string {
	if word, ok := TransTable[id]; !ok {
		return id
	} else {
		return word[lang]
	}
}
