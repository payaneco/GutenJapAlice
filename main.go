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
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"gopkg.in/kyokomi/emoji.v1"
)

//ファイル存在チェック
func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func Download(filename string, url string) bool {
	out, err := os.Create(filename)
	defer out.Close()
	doc, err := http.Get(url)
	defer doc.Body.Close()
	if err != nil {
		fmt.Print("url scarapping failed")
		return false
	}
	_, err = io.Copy(out, doc.Body) //使わない変数はblank identifierにしまっちゃおうね
	if err != nil {
		fmt.Print("file output failed")
		return false
	}
	return true
}

type Record struct {
	chapter int
	period  int
	line    int
	text    string
}

//todo 質問
type Bookmark struct {
	Chapter int `json:"chapter"`
	Period  int `json:"period"`
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
	scanner := bufio.NewScanner(transform.NewReader(fp, japanese.ShiftJIS.NewDecoder()))
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
	//ファイルのダウンロード
	files := map[string]string{
		"alice_ja.txt": "http://www.genpaku.org/alice01/alice01j.txt",
		"alice_en.txt": "http://www.gutenberg.org/files/11/11-0.txt",
		"alice_it.txt": "http://www.gutenberg.org/cache/epub/28371/pg28371.txt",
	}
	for filename, url := range files {
		//ファイルがあればスキップ
		if Exists(filename) {
			continue
		}
		if !Download(filename, url) {
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
		Tweet(chapter, period)
		nc, np := GetNextBookmark(chapter, period)
		s := fmt.Sprintf(`{"chapter":%v,"period":%v}`, nc, np)
		ioutil.WriteFile(*filename, []byte(s), 0666)
	}
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

func Tweet(chapter, period int) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	langs := []int{ita, eng, jap}
	for _, lang := range langs {
		s := GetPeriod(db, lang, chapter, period)
		var prefix string
		switch lang {
		case ita: prefix = ":it:"
		case eng: prefix = ":uk:"
		case jap: prefix = ":jp:"
		default: prefix = ""
		}
		emoji.Println(prefix + s)
	}
	//各国語で表示
	//:国旗 99-99(99/11) - 改行込みで14文字
	//38文字くらいで分割して表示?
	//1章1節(99/99) -11文字+改行3+伊英日6 バッファは20
	//伊：CAPITOLO I. GIÙ NELLA CONIGLIERA. -33文字
	//英：CHAPTER I. Down the Rabbit-Hole
	//日：1. うさぎの穴をまっさかさま       -15文字
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
