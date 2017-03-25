package main

import (
	"testing"
)

func TestReplace(t *testing.T) {
	testMap := map[string]string{
		"\nCORO.\n\n(Coro al quale": Replace(ita, "(Coro al quale"),
		"began.\n\nAlice thought":   Replace(eng, "began. Alice thought"),
		"（訳者のおねがい":                  Replace(jap, "\n（訳者のおねがい"),
	}
	for expect, actual := range testMap {
		if expect != actual {
			t.Errorf("Expect: %v, Actual: %v", expect, actual)
		}
	}
}

func TestDownload(t *testing.T) {
	if !Download(ita) {
		t.Error("だめ")
	}
	if !Download(eng) {
		t.Error("だめ")
	}
	if !Download(jap) {
		t.Error("だめ")
	}
}

func TestSlice(t *testing.T) {
	ss := Slice("abc de", 5)
	if len(ss) < 2 {
		t.Errorf("数が足りない: %v", len(ss))
	}
	if ss[0] != "abc" || ss[1] != "de" {
		t.Errorf("値が違う: %v", ss)
	}
}

func TestSliceFixed(t *testing.T) {
	ss := SliceFixed("123あいうえ1", 3)
	if ss[0] != "123" || ss[1] != "あいう" || ss[2] != "え1" {
		t.Errorf("値が違う: %v", ss)
	}
}

func TestTweet(t *testing.T) {
	//Tweet(10, 40)
	//Tweet(12, 72)
	//t.Error("だめ")
}
