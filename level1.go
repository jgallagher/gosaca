package gosaca

// recursive version of ComputeSuffixArray for levels 1+
func computeSuffixArray1(S, SA []int, k int) {
	n := len(S)

	// bit-invert any S-type characters in S
	// TODO - can our caller do this for us as he builds S?
	for i := n - 2; i >= 0; i-- {
		if (S[i+1] < 0 && S[i] <= ^S[i+1]) || // S[i+1] is S-type and we're <= it
			(S[i+1] >= 0 && S[i] < S[i+1]) { // S[i+1] is L-type and we're strictly < it
			S[i] = ^S[i]
		}
	}

	// *********************************************
	// Stage 1: Induced-sort the LMS-substrings of S
	// *********************************************

	// step 1 - initialize SA as empty
	setAllToEmpty(SA)

	// step 2 - put all LMS substrings into buckets based on their first character
	for i := n - 2; i >= 0; i-- {
		if S[i] >= 0 {
			// S[i] is L-type
			continue
		}
		// S[i] is S-type; walk back until S[i-1] is L-type or -1
		for i >= 1 && S[i-1] < 0 {
			// S[i-1] is also S-type
			i--
		}

		if i > 0 {
			// found LMS starting at S[i] - need to insert it into end of its bucket
			// by Property 4.1, S[i] holds pointer to the end of the bucket,
			// but we need to walk backwards if we've already put other LMS suffixes into the same bucket
			end := ^S[i]
			for SA[end] != empty {
				end--
			}
			SA[end] = i
		}
	}

	// step 3 - induced sort the L-type suffixes of S into their buckets
	induceSortL1(S, SA)

	// step 4 - induced sort the S-type suffixes of S into their buckets
	induceSortS1(S, SA)

	// compact all the now-sorted LMS substrings into the first n1 positions of SA
	n1 := 0
	for i := 0; i < n; i++ {
		if SA[i] > 0 && // S[0] is not LMS by definition
			S[SA[i]] < 0 && // S[SA[i]] is S-type
			S[SA[i]-1] >= 0 { // S[SA[i-1]] is L-type
			// S[i] is LMS
			SA[n1] = SA[i]
			n1++
		}
	}

	// *********************************************
	// Stage 2: Rename the LMS substrings
	// *********************************************

	// provably, n1 is at most floor(n/2), so the following overlapping works
	SA1 := SA[:n1] // SA1 overlaps the front of SA
	work := SA[n1:] // workspace overlaps the rest of SA
	S1 := SA[n-n1:] // S1 overlaps the end of SA (including part of "work", but rename deals with that correctly)
	k1 := rename1(S, SA1, work, S1)

	// *********************************************
	// Stage 3: Sort recursively
	// *********************************************
	if k1 == n1 {
		for i := 0; i < n1; i++ {
			SA1[S1[i]] = i
		}
	} else {
		computeSuffixArray1(S1, SA1, k1)
	}

	// NOT DESCRIBED IN PAPER BUT STILL NECESSARY (see SA-IS)
	// We need to undo the renaming of the LMS suffixes.
	// We no longer need S1, so reuse it to hold all the LMS indices.
	j := n1 - 1
	for i := n - 2; i >= 0; i-- {
		if S[i] >= 0 {
			// L-type; ignore
			continue
		}
		// S[i] is S-type; walk backwards to find LMS
		for i >= 1 && S[i-1] < 0 {
			// S[i-1] is also S-type; keep moving back
			i--
		}
		// S[0] is not LMS by definition, but otherwise S[i] is LMS
		if i > 0 {
			S1[j] = i
			j--
		}
	}
	// Now convert SA1 from renamed values to true values.
	for i := 0; i < n1; i++ {
		SA1[i] = S1[SA1[i]]
	}

	// *********************************************
	// Stage 4: Induced-sort SA(S) from SA1(S1)
	// *********************************************

	// step 1 - initialize SA[n1:] as empty
	setAllToEmpty(SA[n1:])

	// step 2 - put all the sorted LMS suffixes of S into their buckets in SA
	for i := n1 - 1; i >= 0; i-- {
		j := SA[i]
		SA[i] = empty
		c := ^S[j]
		if j == 0 {
			panic("unexpected j == 0")
		}
		// look backwards until we find an empty slot in the bucket
		for SA[c] != empty {
			c--
		}
		SA[c] = j
	}

	// step 3 - induced sort the L-type suffixes of S into their buckets
	induceSortL1(S, SA)

	// step 4 - induced sort the S-type suffixes of S into their buckets
	induceSortS1(S, SA)
}

// TODO pre-post
func induceSortL1(S, SA []int) {
	n := len(S)

	// special case to deal with the (virtual) sentinel:
	// S[n-1] is L-type because of the sentinel, and if we were treating
	// the sentinel as a real character, it would be at the front of SA[]
	// (it's effectively stored in "SA[-1]").
	//
	// Because c is L-type, we know SA[c] is empty, so we're in case 1 of section 4.1
	c := S[n-1]
	if c+1 < n && SA[c+1] == empty {
		SA[c+1] = n - 1
		SA[c] = -1
	} else {
		SA[c] = n - 1
	}

	for i := 0; i < n; i++ {
		if SA[i] < 0 {
			// SA[i] is empty or being used as a counter; nothing to do
			continue
		}
		j := SA[i] - 1
		// if we just grabbed the character before an LMS suffix, we need to clear
		// out that LMS suffix (induceSortS1 assumes only L-type suffixes are in SA)
		if S[SA[i]] < 0 {
			SA[i] = empty
		}
		if j < 0 {
			// SA[i] was == 0; there is no preceding character to look at
			continue
		}
		c := S[j]
		if c < 0 {
			// S[j] is S-type; move on
			continue
		}

		switch {
		case SA[c] >= 0:
			// section 4.1 case 2
			// left shift the previous bucket until we find the head
			// NOT MENTIONED IN PAPER: if we overwrite SA[i], we need
			// to *stay here* for the next iteration
			val := SA[c]
			overwroteSAi := (c == i)
			// TODO clean this up
			stop := false
			for x := c - 1; x >= 0 && !stop; x-- {
				prev := val
				val = SA[x]
				if val < 0 && val != empty {
					stop = true
				}
				SA[x] = prev
				if x == i {
					overwroteSAi = true
				}
			}
			if overwroteSAi {
				// decrement i; it will be incremented by the for loop, forcing us to look at SA[i] again next time
				i--
			}
			fallthrough // we now know SA[c] is empty so fall through

		case SA[c] == empty:
			// section 4.1 case 1
			if c+1 < n && SA[c+1] == empty {
				SA[c+1] = j
				SA[c] = -1
			} else {
				SA[c] = j
			}
			break

		default:
			// section 4.1 case 3 (SA[c] is a counter)
			d := SA[c]
			pos := c - d + 1
			if pos < n && SA[pos] == empty {
				SA[pos] = j
				SA[c]--
			} else {
				// left-shift SA[c+1:pos-1], inserting j into SA[pos-1]
				prev := j
				overwroteSAi := (c == i)
				for x := pos - 1; x >= c; x-- {
					val := SA[x]
					SA[x] = prev
					prev = val
					if x == i {
						overwroteSAi = true
					}
				}
				if overwroteSAi {
					i--
				}
			}
			break
		}
	}

	// NOT MENTIONED IN PAPER: We need to go back over SA and fix
	// any leftover counter values via left shifting the buckets appropriately.
	for i := 0; i < n; i++ {
		if SA[i] == empty || SA[i] >= 0 {
			continue
		}
		d := SA[i]
		pos := i - d + 1
		prev := empty
		for x := pos - 1; x >= i; x-- {
			val := SA[x]
			SA[x] = prev
			prev = val
		}
	}
}

// TODO pre-post
func induceSortS1(S, SA []int) {
	n := len(S)

	for i := n - 1; i >= 0; i-- {
		if SA[i] <= 0 {
			// SA[i] is empty or being used as a counter; nothing to do
			continue
		}
		j := SA[i] - 1
		c := ^S[j]
		if c < 0 {
			// S[j]==c is L-type; move on
			continue
		}

		switch {
		case SA[c] >= 0:
			// section 4.2 case 2
			val := SA[c]
			overwroteSAi := (c == i)
			stop := false
			for x := c + 1; x < n && !stop; x++ {
				prev := val
				val = SA[x]
				if val < 0 && val != empty {
					stop = true
				}
				SA[x] = prev
				if x == i {
					overwroteSAi = true
				}
			}
			if overwroteSAi {
				i++
			}
			fallthrough

		case SA[c] == empty:
			// section 4.2 case 1 
			if c-1 >= 0 && SA[c-1] == empty {
				SA[c-1] = j
				SA[c] = -1
			} else {
				SA[c] = j
			}
			break

		default:
			// section 4.2 case 3
			d := SA[c]
			pos := c + d - 1
			if pos >= 0 && SA[pos] == empty {
				SA[pos] = j
				SA[c]--
			} else {
				// right-shift SA[pos+1:c-1], inserting j into SA[pos+1]
				prev := j
				overwroteSAi := (c == i) // TODO - can default to false?
				for x := pos + 1; x <= c; x++ {
					val := SA[x]
					SA[x] = prev
					prev = val
					if x == i {
						overwroteSAi = true
					}
				}
				if overwroteSAi {
					i++
				}
			}
			break
		}
	}
	// TODO PICKUP AT LINE 320
}
