package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
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

const dbFile = "./test.db"

func PushDB(table string, list []Record) {
	//64bit windowsで使うにはgccが必要です。
	//http://twinbird-htn.hatenablog.com/entry/2016/07/01/133824
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE "` + table + `" ("id" INTEGER PRIMARY KEY AUTOINCREMENT, "lang" CHAR(2), "chapter" INTEGER, "period" INTEGER, "line" INTEGER, "text" VARCHAR(2000))`)
	if err != nil {
		panic(err)
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	stmt, err := tx.Prepare(`INSERT INTO "` + table + `" ("lang", "chapter", "period", "line", "text") VALUES (?, ?, ?, ?, ?) `)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	for i := range list {
		if _, err = stmt.Exec(table, list[i].chapter, list[i].period, list[i].line, list[i].text); err != nil {
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

func GetEnRecords() []Record {
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

func GetItRecords() []Record {
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
	PushDB("en", GetEnRecords())
	PushDB("it", GetItRecords())
}

func main() {
	index := flag.Int("i", 0, "インデックスを設定")
	flag.Parse()
	if *index == 0 {
		GetFiles()
	} else {
		Tweet(*index)
	}
}
func Tweet(index int) {
	fmt.Println(index)
	dbFile := "./test.db"

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("select text from (select * from en where chapter = 1 union all select * from it where chapter = 1 order by chapter, period, line, lang)")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	//todo この腐ったコードを何とかしろ！
	for i := 0; i < index; i++ {
		rows.Next()
	}
	var text string
	err = rows.Scan(&text)
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
}
