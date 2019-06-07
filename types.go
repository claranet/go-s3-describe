package main

type s3Bucket struct {
	region          string
	name            string
	bucketSizeBytes map[string]float64
	numberOfObjects float64
	isPublic        bool
}
