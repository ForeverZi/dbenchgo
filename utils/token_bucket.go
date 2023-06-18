package utils

type TokenBucket struct {
	tokenChan chan struct{}
}

func NewTokenBucket(tokenCount int, fillToken bool) *TokenBucket {
	bucket := &TokenBucket{
		tokenChan: make(chan struct{}, tokenCount),
	}
	if fillToken {
		for i := 0; i < tokenCount; i++ {
			bucket.Put()
		}
	}
	return bucket
}

func (bucket *TokenBucket) Take() {
	<-bucket.tokenChan
}

func (bucket *TokenBucket) Put() {
	bucket.tokenChan <- struct{}{}
}
