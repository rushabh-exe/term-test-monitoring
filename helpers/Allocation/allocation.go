package allocate

var startST = 1

func Allot(c, cap int64, i int) int64 {
	if i%2 == 1 {
		return c
	}
	if cap == c {
		startST = 1
		return cap
	}
	if cap > c {
		return c
	}
	if cap < c {
		l_cap := c - cap
		if startST == 1 {
			startST = int(c) + 1
		} else {
			startST += int(cap)
		}
		return c - l_cap
	}

	return -1
}
