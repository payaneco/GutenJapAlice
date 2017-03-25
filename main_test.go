package main

import (
	"testing"
)

func TestReplace(t *testing.T) {
	testMap := map[string]string{
		"\nCORO.\n\n(Coro al quale":                      Replace(ita, "(Coro al quale"),
		"and he’s treading on my tail.\nSee how eagerly": Replace(eng, "and he’s treading on my tail.\n\nSee how eagerly"),
		"began.\n\nAlice thought":                        Replace(eng, "began. Alice thought"),
		"（訳者のおねがい":                                       Replace(jap, "\n（訳者のおねがい"),
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

func TestTweet(t *testing.T) {
	Tweet(10, 25)
	//Tweet(12, 72)
	//t.Error("だめ")
}
