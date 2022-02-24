package chamux

// Searches for +END
func splitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// That means we've scanned to the end
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	// search for +END
	var gotPlus, gotE, gotN bool
	var prev byte
	for i, b := range data {
		switch {
		case b == '+':
			prev = b
			gotPlus = true
		case gotPlus && b == 'E' && prev == '+':
			prev = b
			gotE = true
		case gotPlus && gotE && b == 'N' && prev == 'E':
			prev = b
			gotN = true
		case gotPlus && gotE && gotN && b == 'D' && prev == 'N':
			return i + 1, data[:len(data)-4], nil
		}
		gotPlus = false
	}
	// The reader contents processed here are all read out,
	// but the contents are not empty, so the remaining data needs to be returned.
	if atEOF {
		return len(data), data, nil
	}
	// Represents that you can't split up now, and requests more data from Reader
	return 0, nil, nil
}
