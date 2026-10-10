package ui

type Language int

const (
	CN Language = iota
	EN
	JP
)

type Word [3]string

var TransTable = map[string]Word{
	"home":       {CN: "主页", EN: "Home", JP: "ホーム"},
	"color_mode": {CN: "颜色模式", EN: "Color Mode", JP: "カラーモード"},
	"bookshelf":  {CN: "书架", EN: "Bookshelf", JP: "本棚"},
	"setting":    {CN: "设置", EN: "Settings", JP: "設定"},
	"reset":      {CN: "重置", EN: "Reset", JP: "リセット"},
	"next":       {CN: "下一页", EN: "Next", JP: "次へ"},
	"last":       {CN: "上一页", EN: "Previous", JP: "前へ"},
	"rotate":     {CN: "旋转", EN: "Rotate", JP: "回転"},
	"shader":     {CN: "着色器", EN: "Shader", JP: "シェーダー"},
	"layout":     {CN: "布局", EN: "Layout", JP: "レイアウト"},
	"record":     {CN: "记录", EN: "Record", JP: "记录 / 履歴"},
}

func GetWord(id string, lang Language) string {
	if word, ok := TransTable[id]; !ok {
		return id
	} else {
		return word[lang]
	}
}
