package ui

type Language int

const (
	CN Language = iota
	EN
	JP
)

type Word [3]string

var TransTable = map[string]Word{
	"home":       {CN: "主页", EN: "Home", JP: "ホーム"}, // Go 允许通过索引赋值
	"color_mode": {"颜色模式", "Color Mode", "カラーモード"},
	"bookshelf":  {"书架", "Bookshelf", "本棚"},
	"setting":    {"设置", "Settings", "設定"},
}

func GetWord(id string, lang Language) string {
	if word, ok := TransTable[id]; !ok {
		return id
	} else {
		return word[lang]
	}
}
