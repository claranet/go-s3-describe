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

// Parse a S3 Bucket Policy into a structure
// func parseJsonBucketPolicy(jsonInput interface{}) bucketPolicy {
// 	// see https://blog.golang.org/json-and-go
// 	m := jsonInput.(map[string]interface{})
// 	var bpolicy bucketPolicy
// 	bpolicy.Version = m["Version"].(string)
// 	n := m["Statement"].([]interface{})
// 	for _, v := range n {
// 		var state statement
// 		for k, w0 := range v.(map[string]interface{}) {
// 			switch k {
// 			case "Sid":
// 				state.Sid = w0.(string)
// 			case "Effect":
// 				state.Effect = w0.(string)
// 			case "Action":
// 				switch w1 := w0.(type) {
// 				case string:
// 					state.Action = append(state.Action, w0.(string))
// 				case []interface{}:
// 					for _, w2 := range w1 {
// 						state.Action = append(state.Action, w2.(string))
// 					}
// 				}
// 			case "Principal":
// 				switch w1 := w0.(type) {
// 				case string:
// 					state.Principal = append(state.Principal, w0.(string))
// 				case map[string]interface{}:
// 					for _, w2 := range w1 {
// 						switch w3 := w2.(type) {
// 						case string:
// 							state.Principal = append(state.Principal, w2.(string))
// 						case []interface{}:
// 							for _, w4 := range w3 {
// 								state.Principal = append(state.Principal, w4.(string))
// 							}
// 						}
// 					}
// 				}
// 			case "Resource":
// 				switch w1 := w0.(type) {
// 				case string:
// 					state.Resource = append(state.Resource, w0.(string))
// 				case []interface{}:
// 					for _, w2 := range w1 {
// 						state.Resource = append(state.Resource, w2.(string))
// 					}
// 				}

// 			}

// 		}
// 		bpolicy.Statement = append(bpolicy.Statement, state)
// 	}
// 	return bpolicy
// }
