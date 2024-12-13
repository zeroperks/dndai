package dice

type KarmicDiceRNG struct {
	rng      RNG
	lastRoll uint64
	karma    int
}

func NewKarmicDiceRNG(rng RNG) *KarmicDiceRNG {
	return &KarmicDiceRNG{rng: rng}
}

func (k *KarmicDiceRNG) Next() uint64 {
	roll := k.rng.Next()
	if k.karma > 0 {
		roll = (roll + k.lastRoll) / 2
		k.karma--
	} else {
		k.karma++
	}
	k.lastRoll = roll
	return roll
}
