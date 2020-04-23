package fulltext

import (
	"sort"
	"strings"
)

// Index implements a simple fulltext search
type Index struct {
	idx map[string][]int
}

// New creates a fulltext search index for the given slice of documents.
// The documents are non-empty strings of lowercase words divided by spaces.
// Each document is identified by its index in the docs slice.
func New(docs []string) *Index {
	idx := new(Index)
	idx.idx = make(map[string][]int)

	for docID, doc := range docs {
		terms := strings.Split(doc, " ")
		for _, term := range terms {
			ds, ok := idx.idx[term]
			if !ok {
				ds = make([]int, 0, 1)
			}
			idx.idx[term] = append(ds, docID)
		}
	}
	for term, dss := range idx.idx {
		idx.idx[term] = sortDedup(dss)
	}
	return idx
}

// Search returns a sorted slice of unique ids of documents
// that contain all words from the query, not necessarily in the given order.
// The query is always a string of lowercase words divided by spaces.
func (idx *Index) Search(query string) []int {
	if len(query) == 0 {
		return []int{}
	}
	found := make([]int, 0)

	terms := strings.Split(query, " ")
	for i, term := range terms {
		ds, ok := idx.idx[term]
		if !ok {
			ds = []int{}
		}
		if i == 0 {
			found = append(found, ds...)
		} else {
			found = intersect(found, ds)
		}
		if len(found) == 0 {
			break
		}
	}
	return found
}

func sortDedup(s []int) []int {
	if len(s) == 0 {
		return s
	}
	sort.Ints(s)
	j := 1
	for i := 1; i < len(s); i++ {
		if s[i] != s[j-1] {
			s[j] = s[i]
			j++
		}
	}
	return s[0:j]
}

func intersect(s1, s2 []int) []int {
	if len(s1) == 0 || len(s2) == 0 {
		return []int{}
	}
	if len(s1) > len(s2) {
		s1, s2 = s2, s1
	}
	result := make([]int, 0)
	i := 0

	for j := 0; j < len(s1); j++ {
		for i < len(s2) && s2[i] < s1[j] {
			i++
		}
		if i >= len(s2) {
			break
		}
		if s1[j] == s2[i] {
			result = append(result, s1[j])
		}
	}
	return result
}
