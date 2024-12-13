package dice

type RandomNumberGenerator struct {
	sm *SplitMix64
}

func NewRandomNumberGenerator(seed uint64) *RandomNumberGenerator {
	return &RandomNumberGenerator{sm: NewSplitMix64(seed)}
}

func (rng *RandomNumberGenerator) Next() uint64 {
	return rng.sm.Next()
}
