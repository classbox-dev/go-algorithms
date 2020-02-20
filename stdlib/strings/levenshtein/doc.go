// Package levenshtein implements functions to compute Levenshtein distance between strings.
/*
Definitions

Given two strings, we define transformation between them as a sequence of the following elementary operations:

	I - insert a character
	R - replace a character
	D - delete a character
	M - match, leave a character unchanged

For example, the string "vintner" can be transformed into "writers" as follows:

	v intner   // source string
	wri t ers  // destination string
	RIMDMDMMI  // editing operations (edit transcript)

An edit transcript is a string over the alphabet {I, R, D, M} that describes such transformation.

The Levenshtein distance is a minimum number of operations I, R, and D required to transform one string to another.

In this package we are only interested in edit transcripts that correspond to the Levenshtein distance i.e. contain the minimum
number of operations I, R, and D.

Resources

https://en.wikipedia.org/wiki/Levenshtein_distance

Dan Gusfield, "Algorithms on Strings, Trees, and Sequences", Section 11.2
*/
package levenshtein
