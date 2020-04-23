// Package fulltext implements a silly fulltext search.
/*
Given a list of "documents", the package implements an extremely simple fulltext search.
It allows to find documents containing all words from the query, not necessarily in the given order.

The lack of requirement on matching the word order is important.
It enables a trivial implementation _without_ specialised algorithms or data structures (such as suffix trees/arrays).
*/
package fulltext
