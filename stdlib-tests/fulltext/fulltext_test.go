package fulltext

import (
	"hsecode.com/stdlib-tests/internal/utils"
	"hsecode.com/stdlib/strings/fulltext"
	"reflect"
	"strings"
	"testing"
)

var dictWords = [...]string{
	"e", "podostemonaceous", "podsolization", "squamaceous", "tinctumutation", "lifelikeness", "valorisation", "vallisneriaceous",
	"zumbooruk", "chartographically", "benzalphenylhydrazone", "seculars", "selzogene", "uncommutative", "helicopted", "abaca",
	"ament", "amentaceous", "amental", "amentia", "Amentiferae", "amentiferous", "amentiform", "amentulum", "amentum", "amerce",
	"amerceable", "amercement", "amercer", "amerciament", "America", "American", "Americana", "Americanese", "Americanism",
	"Americanist", "Americanistic", "Americanitis", "Americanization", "Americanize", "Americanizer", "Americanly", "Americanoid",
	"Americaward", "Americawards",
}

var extraWords = [...]string{
	"makefile", "microsoft", "mutex", "namespace", "num", "ocaml", "png", "postscript", "pragma", "preprocessor", "radix", "ru", "sizeof",
}

func TestUnit__EmptyIndex(t *testing.T) {
	errorMsg := "Something is wrong with searching in an empty index"

	idx := fulltext.New([]string{})
	if !reflect.DeepEqual(idx.Search(""), []int{}) {
		t.Fatal(errorMsg)
	}
	if !reflect.DeepEqual(idx.Search("a"), []int{}) {
		t.Fatal(errorMsg)
	}
}

func TestUnit__EmptyQuery(t *testing.T) {
	idx := fulltext.New([]string{"a", "b", "c"})
	if !reflect.DeepEqual(idx.Search(""), []int{}) {
		t.Fatal("Something is wrong with searching in an empty query")
	}
}

func TestUnit__DuplicateDocs(t *testing.T) {
	idx := fulltext.New([]string{
		"a b c",
		"a b c",
		"a b c",
	})
	if !reflect.DeepEqual(idx.Search("a"), []int{0, 1, 2}) {
		t.Fatal("Something is wrong with searching duplicate documents")
	}
}

func TestUnit__NonExisting(t *testing.T) {
	errorMsg := "Something is wrong with searching non-existing documents"

	idx := fulltext.New([]string{
		"a b c",
		"d e",
		"z",
		"m",
		"k l",
		"u i p",
	})
	if !reflect.DeepEqual(idx.Search("t"), []int{}) {
		t.Fatal(errorMsg)
	}
	if !reflect.DeepEqual(idx.Search("t a"), []int{}) {
		t.Fatal(errorMsg)
	}
	if !reflect.DeepEqual(idx.Search("a t"), []int{}) {
		t.Fatal(errorMsg)
	}
}

func TestUnit__FoundUnique(t *testing.T) {
	errorMsg := "Something is wrong with searching a unique document"

	idx := fulltext.New([]string{
		"a b c",
		"d e",
		"z",
		"m",
		"k l",
		"u i p",
	})
	if !reflect.DeepEqual(idx.Search("m"), []int{3}) {
		t.Fatal(errorMsg)
	}
	if !reflect.DeepEqual(idx.Search("u i"), []int{5}) {
		t.Fatal(errorMsg)
	}
	if !reflect.DeepEqual(idx.Search("a b c"), []int{0}) {
		t.Fatal(errorMsg)
	}
}

func TestUnit__FoundMultiple(t *testing.T) {
	errorMsg := "Something is wrong with searching multiple document"

	idx := fulltext.New([]string{
		"a b c",
		"c b e",
		"b c e",
		"e c b",
		"x",
	})
	if !reflect.DeepEqual(idx.Search("e b"), []int{1, 2, 3}) {
		t.Fatal(errorMsg)
	}
	if !reflect.DeepEqual(idx.Search("a b"), []int{0}) {
		t.Fatal(errorMsg)
	}
	if !reflect.DeepEqual(idx.Search("b a"), []int{0}) {
		t.Fatal(errorMsg)
	}
	if !reflect.DeepEqual(idx.Search("x a"), []int{}) {
		t.Fatal(errorMsg)
	}
}

func TestUnit__AnyOrder(t *testing.T) {
	errorMsg := "Something is wrong with searching documents with various word order"

	idx := fulltext.New([]string{
		"a b c d",
		"d a b c",
		"c d a b",
		"b c d a",
		"b d a c",
	})
	if !reflect.DeepEqual(idx.Search("d c a b"), []int{0, 1, 2, 3, 4}) {
		t.Fatal(errorMsg)
	}
}

func TestUnit__ByWords(t *testing.T) {
	errorMsg := "Looks like the search does not respect word boundaries"

	idx := fulltext.New([]string{
		"a b c d",
		"dab c",
		"c dab",
		"bc da",
		"bdac",
	})
	if !reflect.DeepEqual(idx.Search("bdac"), []int{4}) {
		t.Fatal(errorMsg)
	}
	if !reflect.DeepEqual(idx.Search("da"), []int{3}) {
		t.Fatal(errorMsg)
	}
	if !reflect.DeepEqual(idx.Search("a"), []int{0}) {
		t.Fatal(errorMsg)
	}
}

func TestUnit__LargeDict(t *testing.T) {
	docs := make([]string, 0, 600000)
	for i, word := range dictWords {
		n := (i + 1) * (i + 1) * 19
		for j := 0; j < n; j++ {
			docs = append(docs, word)
		}
	}
	idx := fulltext.New(docs)
	for i, word := range dictWords {
		n := (i + 1) * (i + 1) * 19
		results := idx.Search(word)
		if len(results) != n {
			t.Fatal("Unexpected number of documents found")
		}
		if docs[results[0]] != word {
			t.Fatal("Unexpected document found")
		}
	}
}

func TestPerf__LargeDict(t *testing.T) {
	docs := make([]string, 0, 600000)
	for i, word := range dictWords {
		n := (i + 1) * (i + 1) * 19
		for j := 0; j < n; j++ {
			docs = append(docs, word)
		}
	}
	idx := fulltext.New(docs)
	for r := 0; r <= 20; r++ {
		for _, word := range dictWords {
			idx.Search(word)
		}
	}
}

func makeDocWords(words []string, n int) []string {
	total := len(words)
	ws := make([]string, 0, n)
	for j := 0; j < n; j++ {
		iw := utils.Rand.Intn(total)
		ws = append(ws, words[iw])
	}
	return ws
}

func shuffle(a []string) {
	utils.Rand.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
}

func TestPerf__SearchMultiple(t *testing.T) {
	N := 30000
	NW := 15

	docs := make([]string, 0, N+100)

	for i := 0; i < N; i++ {
		nw := utils.Rand.Intn(NW) + 1
		ws := makeDocWords(dictWords[:], nw)
		doc := strings.Join(ws, " ")
		docs = append(docs, doc)
	}

	for i, extraWord := range extraWords {
		ws := makeDocWords(dictWords[:], utils.Rand.Intn(NW)+1)
		ws = append(ws, extraWord)

		reps := i + 1

		for j := 0; j < reps; j++ {
			docs = append(docs, strings.Join(ws, " "))
		}

		idx := fulltext.New(docs)
		for r := 0; r < 50; r++ {
			shuffle(ws)
			q := strings.Join(ws, " ")

			if len(idx.Search(q)) != reps {
				t.Fatal("Unexpected number of documents found")
			}
		}
		docs = docs[0:N]
	}
}
