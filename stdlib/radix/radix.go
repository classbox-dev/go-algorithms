package radix

// Sort sorts the given slice of uint64 using the radix sort algorithm.
//
// The reference implementation uses least significant digit (LSD) version of the algorithm,
// takes Θ(8n) time and uses Θ(n) of additional memory.
// The data is sorted byte by byte (starting from the least significant one) via stable counting sort.
func Sort(data []uint64) {
	old := data
	result := make([]uint64, len(data))

	for offset := 0; offset < 64; offset += 8 {
		ps := [256]int{}
		for _, v := range data {
			k := (v >> offset) & 0xff
			ps[k]++
		}

		for i, sum := 0, 0; i < 256; i++ {
			tmp := sum + ps[i]
			ps[i] = sum
			sum = tmp
		}
		for _, v := range data {
			k := (v >> offset) & 0xff
			result[ps[k]] = v
			ps[k]++
		}
		data, result = result, data
	}
	copy(old, data)
}
