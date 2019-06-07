package main

// Functions to sort buckets by size
type bySize []s3Bucket

func (s bySize) Len() int      { return len(s) }
func (s bySize) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s bySize) Less(i, j int) bool {
	var a, b float64
	for _, c := range s[i].bucketSizeBytes {
		a += c
	}
	for _, c := range s[j].bucketSizeBytes {
		b += c
	}
	return b < a
}
