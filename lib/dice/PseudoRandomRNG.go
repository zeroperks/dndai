package dice

type PseudoRandomRNG struct {
	sm *SplitMix64
}

func NewPseudoRandomRNG(seed uint64) *PseudoRandomRNG {
	return &PseudoRandomRNG{sm: NewSplitMix64(seed)}
}

func (prng *PseudoRandomRNG) Next() uint64 {
	return prng.sm.Next()
}
