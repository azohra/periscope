package tools

// Pair pair
type Pair struct {
	svcLabel, deployLabel interface{}
}

// Int32Ptr int 32 pointer
func Int32Ptr(i int32) *int32 {
	return &i
}
