package chamux

// Searches for +END
func splitPlusEnd(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// That means we've scanned to the end
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	var gotE, gotN bool
	var prev byte
	for i, b := range data {
		if prev == '+' && b == 'E' {
			gotE = true
		} else if gotE && prev == 'E' && b == 'N' {
			gotN = true
		} else if gotN && prev == 'N' && b == 'D' {
			return i + 1, data[:i-3], nil
		} else {
			gotE = false
			gotN = false
		}
		prev = b
	}
	// The reader contents processed here are all read out,
	// but the contents are not empty, so the remaining data needs to be returned.
	if atEOF {
		return len(data), data, nil
	}
	// Represents that you can't split up now, and requests more data from Reader
	return 0, nil, nil
}
