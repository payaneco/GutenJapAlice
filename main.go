package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"gopkg.in/kyokomi/emoji.v1"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Record struct {
	chapter int
	period  int
	line    int
	text    string
}

//フィールド名は大文字でないと外部パッケージから参照できない
type Bookmark struct {
	Chapter int `json:"chapter"`
	Period  int `json:"period"`
}

type Replacer struct {
	Lang        int    `json:"lang"`
	Target      string `json:"target"`
	Replacement string `json:replacement`
}

const dbFile = "./alice.db"
const ita = 1
const eng = 2
const jap = 3

var langMap = map[int]string{
	ita: "伊",
	eng: "英",
	jap: "日",
}

//ファイル存在チェック
func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func Download(lang int) bool {
	fileMap := map[int]string{
		ita: "alice_it.txt",
		eng: "alice_en.txt",
		jap: "alice_ja.txt",
	}
	urlMap := map[int]string{
		ita: "http://www.gutenberg.org/cache/epub/28371/pg28371.txt",
		eng: "http://www.gutenberg.org/files/11/11-0.txt",
		jap: "http://www.genpaku.org/alice01/alice01j.txt",
	}
	filename := fileMap[lang]
	//ファイルがあればスキップ
	if Exists(filename) {
		return true
	}
	url := urlMap[lang]
	doc, err := http.Get(url)
	defer doc.Body.Close()
	if err != nil {
		fmt.Print("url scarapping failed")
		return false
	}
	arr, err := ioutil.ReadAll(doc.Body)
	if err != nil {
		fmt.Print("file read failed")
		return false
	}
	text := string(arr)
	if lang == jap {
		//sjis -> utf-8
		text, _, err = transform.String(japanese.ShiftJIS.NewDecoder(), text)
		if err != nil {
			fmt.Print("file encode failed")
			return false
		}
	}
	text = strings.Replace(Replace(lang, text), "\r", "", -1)
	text = strings.Replace(text, "\n", "\r\n", -1)
	ioutil.WriteFile(filename, []byte(text), 0666)
	return true
}

func Replace(lang int, text string) string {
	bytes, err := ioutil.ReadFile("replace.json")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	// JSONデコード
	var replacers []Replacer
	if err := json.Unmarshal(bytes, &replacers); err != nil {
		log.Fatal(err)
	}
	s := text
	for _, r := range replacers {
		if lang != r.Lang {
			continue
		}
		old := r.Target
		new := strings.Replace(r.Replacement, `$1`, old, -1)
		s = strings.Replace(s, old, new, -1)
	}
	return s

}

func CreateDB() {
	//64bit windowsで使うにはgccが必要です。
	//http://twinbird-htn.hatenablog.com/entry/2016/07/01/133824
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE "main" ("id" INTEGER PRIMARY KEY AUTOINCREMENT, "lang" INTEGER, "chapter" INTEGER, "period" INTEGER, "line" INTEGER, "text" VARCHAR)`)
	if err != nil {
		panic(err)
	}
}

func PushDB(lang int, list []Record) {
	//64bit windowsで使うにはgccが必要です。
	//http://twinbird-htn.hatenablog.com/entry/2016/07/01/133824
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	stmt, err := tx.Prepare(`INSERT INTO "main" ("lang", "chapter", "period", "line", "text") VALUES (?, ?, ?, ?, ?) `)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	for i := range list {
		if _, err = stmt.Exec(lang, list[i].chapter, list[i].period, list[i].line, list[i].text); err != nil {
			panic(err)
		}
	}
	tx.Commit()
}

func GetChapter(c string) int {
	switch c {
	case "I":
		return 1
	case "II":
		return 2
	case "III":
		return 3
	case "IV":
		return 4
	case "V":
		return 5
	case "VI":
		return 6
	case "VII":
		return 7
	case "VIII":
		return 8
	case "IX":
		return 9
	case "X":
		return 10
	case "XI":
		return 11
	case "XII":
		return 12
	default:
		return 999
	}
	return 999
}

func GetEngRecords() []Record {
	fp, err := os.Open("alice_en.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(fp)
	defer fp.Close()
	list := []Record{}
	//一行ずつ読み取って処理する
	c := 0
	p := 1
	l := 1
	isSkipping := false
	for scanner.Scan() {
		line := scanner.Text()
		r := regexp.MustCompile(`CHAPTER (.+?)\. (.+)`)
		if r.MatchString(line) {
			ss := r.FindStringSubmatch(line)
			c = GetChapter(ss[1])
			p = 1
			l = 1
		} else if c == 0 {
			//章立ての前はすべて飛ばす
			continue
		} else if regexp.MustCompile(`^[ \*]*$`).MatchString(line) {
			//複数行空白があっても段落は1つのみ加算
			if !isSkipping {
				p++
				l = 1
			}
			isSkipping = true
			continue
		} else if regexp.MustCompile(`^ *THE +END *$`).MatchString(line) {
			//THE END
			break
		}
		isSkipping = false
		list = append(list, Record{chapter: c, period: p, line: l, text: line})
		l++
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return list
}

func GetItaRecords() []Record {
	fp, err := os.Open("alice_it.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(fp)
	defer fp.Close()
	list := []Record{}
	//一行ずつ読み取って処理する
	c := 0
	p := 1
	l := 1
	isSkipping := false
	for scanner.Scan() {
		line := scanner.Text()
		r := regexp.MustCompile(`CAPITOLO (.+?)\.`)
		if r.MatchString(line) {
			ss := r.FindStringSubmatch(line)
			c = GetChapter(ss[1])
			p = 1
			l = 1
			//2行先のタイトルを取得
			scanner.Scan()
			scanner.Scan()
			line = line + " " + scanner.Text()
		} else if c == 0 {
			//章立ての前はすべて飛ばす
			continue
		} else if strings.TrimSpace(line) == "[Illustrazione]" {
			continue
		} else if regexp.MustCompile(`^[ \*]*$`).MatchString(line) {
			//複数行空白があっても段落は1つのみ加算
			if !isSkipping {
				p++
				l = 1
			}
			isSkipping = true
			continue
		} else if regexp.MustCompile(`^ *FINE. *$`).MatchString(line) {
			//FINE
			break
		}
		isSkipping = false
		list = append(list, Record{chapter: c, period: p, line: l, text: line})
		l++
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return list
}

func GetJapRecords() []Record {
	fp, err := os.Open("alice_ja.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(fp)
	defer fp.Close()
	list := []Record{}
	//一行ずつ読み取って処理する
	c := 0
	p := 1
	l := 1
	isSkipping := false
	for scanner.Scan() {
		line := scanner.Text()
		r := regexp.MustCompile(`^(\d+?)\. `)
		if r.MatchString(line) {
			ss := r.FindStringSubmatch(line)
			c, _ = strconv.Atoi(ss[1])
			p = 1
			l = 1
		} else if c == 0 {
			//章立ての前はすべて飛ばす
			continue
		} else if regexp.MustCompile(`^[ \*\-　]*$`).MatchString(line) {
			//複数行空白があっても段落は1つのみ加算
			if !isSkipping {
				p++
				l = 1
			}
			isSkipping = true
			continue
		} else if regexp.MustCompile(`訳したやつのいろんな言い訳`).MatchString(line) {
			//めでたしめでたし
			break
		}
		isSkipping = false
		list = append(list, Record{chapter: c, period: p, line: l, text: line})
		l++
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return list
}

func GetFiles() {
	langs := []int{ita, eng, jap}
	for _, lang := range langs {
		//ファイルのダウンロード
		if !Download(lang) {
			return
		}
	}
	os.Remove(dbFile)
	CreateDB()
	PushDB(eng, GetEngRecords())
	PushDB(ita, GetItaRecords())
	PushDB(jap, GetJapRecords())
}

func main() {
	filename := flag.String("b", "", "bookmark.jsonのパス")
	flag.Parse()
	if *filename == "" {
		GetFiles()
	} else {
		chapter, period := GetBookmark(*filename)
		sliceMap := GetSliceMap(chapter, period)
		TweetMap(sliceMap)
		SaveNextBookmark(chapter, period, filename)
	}
}
func SaveNextBookmark(chapter int, period int, filename *string) {
	nc, np := GetNextBookmark(chapter, period)
	s := fmt.Sprintf(`{"chapter":%v,"period":%v}`, nc, np)
	ioutil.WriteFile(*filename, []byte(s), 0666)
}

func QueryFirstInt(db *sql.DB, query string, args int) int {
	rows, err := db.Query(query, args)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	rows.Next()
	var val int
	err = rows.Scan(&val)
	if err != nil {
		panic(err)
	}
	return val
}

func GetNextBookmark(chapter, period int) (int, int) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	qPeriod := `select ifnull(max(period), 0) from main where chapter = ?`
	mPeriod := QueryFirstInt(db, qPeriod, chapter)
	if mPeriod > period {
		return chapter, period + 1
	}
	qChapter := `select max(chapter) from main`
	mChapter := QueryFirstInt(db, qChapter, 0)
	if mChapter > chapter {
		return chapter + 1, 1
	}
	return 1, 1

}

func GetBookmark(filename string) (int, int) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	// JSONデコード
	var b Bookmark
	if err := json.Unmarshal(bytes, &b); err != nil {
		log.Fatal(err)
	}
	return b.Chapter, b.Period
}

func TweetMap(sliceMap map[int][]string) {
	//各国語で表示
	//イタリア語と英語を交互に表示
	max := len(sliceMap[ita])
	if max < len(sliceMap[eng]) {
		max = len(sliceMap[eng])
	}
	for i := 0; i < max; i++ {
		if i < len(sliceMap[ita]) {
			emoji.Println(sliceMap[ita][i])
		}
		if i < len(sliceMap[eng]) {
			emoji.Println(sliceMap[eng][i])
		}
	}
	//日本語を最後に表示
	for _, js := range sliceMap[jap] {
		emoji.Println(js)
	}
}

func GetSliceMap(chapter int, period int) map[int][]string {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	langs := []int{ita, eng, jap}
	sliceMap := make(map[int][]string)
	for _, lang := range langs {
		s := GetPeriod(db, lang, chapter, period)
		var format string
		var ss []string
		switch lang {
		case ita:
			format = ":it:%v-%v(%v/%v)\n%v"
			ss = Slice(s, 120)
		case eng:
			format = ":uk:%v-%v(%v/%v)\n%v"
			ss = Slice(s, 120)
		case jap:
			format = ":jp:%v-%v(%v/%v)\n%v"
			ss = SliceFixed(s, 120)
		default:
			format = ""
		}
		var texts []string
		for i, body := range ss {
			text := fmt.Sprintf(format, chapter, period, i+1, len(ss), body)
			texts = append(texts, text)

		}
		sliceMap[lang] = texts
	}
	return sliceMap
}

func Slice(text string, max int) []string {
	var slice []string
	var sb bytes.Buffer
	for _, s := range strings.Fields(text) {
		if sb.Len()+len(s) >= max {
			slice = append(slice, sb.String())
			sb.Reset()
			sb.WriteString(s)
		} else {
			if sb.Len() != 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(s)
		}
	}
	slice = append(slice, sb.String())
	return slice
}

func SliceFixed(text string, splitlen int) []string {
	var slice []string
	runes := []rune(text)
	for i := 0; i < len(runes); i += splitlen {
		if i+splitlen < len(runes) {
			slice = append(slice, string(runes[i:(i + splitlen)]))
		} else {
			slice = append(slice, string(runes[i:]))
		}
	}
	return slice
}

func GetPeriod(db *sql.DB, lang, chapter, period int) string {
	q := `select text from (select * from main where lang = ? and chapter = ? and period = ? order by line)`
	rows, err := db.Query(q, lang, chapter, period)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var sb bytes.Buffer
	var text string
	for rows.Next() {
		err = rows.Scan(&text)
		if err != nil {
			panic(err)
		}
		if lang == jap {
			sb.WriteString(strings.TrimSpace(strings.Trim(text, "　")))
		} else {
			sb.WriteString(" " + strings.TrimSpace(text))
		}
	}
	return strings.TrimSpace(sb.String())
}
