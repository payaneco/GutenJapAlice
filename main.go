package main

import (
	"fmt"
	"os"
	"net/http"
	"io"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"bufio"
	"regexp"
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
	period int
	line int
	text string
}

func PushDB(table string, list []Record) {
	//64bit windowsで使うにはgccが必要です。
	//http://twinbird-htn.hatenablog.com/entry/2016/07/01/133824
	var dbFile string = "./test.db"

	os.Remove( dbFile )

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil { panic(err) }
	defer db.Close()

	_, err = db.Exec( `CREATE TABLE "` + table + `" ("id" INTEGER PRIMARY KEY AUTOINCREMENT, "chapter" INTEGER, "period" INTEGER, "line" INTEGER, "text" VARCHAR(2000))` )
	if err != nil { panic(err) }

	tx, err := db.Begin()
	if err != nil { panic(err) }
	stmt, err := tx.Prepare(`INSERT INTO "` + table + `" ("chapter", "period", "line", "text") VALUES (?, ?, ?, ?) `)
	//stmt, err := db.Prepare( `INSERT INTO "ja" ("chapter", "period", "line", "text") VALUES (?, ?, ?, ?) ` )
	if err != nil { panic(err) }
	defer stmt.Close()

	for i := range list {
		if _, err = stmt.Exec(list[i].chapter, list[i].period, list[i].line, list[i].text); err != nil { panic(err) }
	}
	tx.Commit()
}

func GetChapter(c string) int {
	chMap := map[string]int{
		"I": 1,
		"II": 2,
		"III": 3,
		"IV": 4,
		"V": 5,
		"VI": 6,
		"VII": 7,
		"VIII": 8,
		"IX": 9,
		"X": 10,
		"XI": 11,
		"XII": 12,
	}
	if val, ok := chMap[c]; ok {
		return val
	}
	fmt.Println(c)
	return 999
}

func GetEnRecords() []Record{
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
		r := regexp.MustCompile(`CHAPTER (.+?). (.+)`)
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
		list = append(list, Record{chapter:c, period:p, line:l, text:line})
		l++
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return list
}

func main() {
	//ファイルのダウンロード
	files := map[string]string{
		"alice_ja.txt" : "http://www.genpaku.org/alice01/alice01j.txt",
		"alice_en.txt" : "http://www.gutenberg.org/files/11/11-0.txt",
		"alice_it.txt" : "http://www.gutenberg.org/cache/epub/28371/pg28371.txt",
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
	PushDB("en", GetEnRecords())
}
