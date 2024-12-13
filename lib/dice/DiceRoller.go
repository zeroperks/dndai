package dice

import "fmt"

type DiceRoller struct {
	rng RNG
}

func NewDiceRoller(rng RNG) *DiceRoller {
	return &DiceRoller{rng: rng}
}

func (dr *DiceRoller) RollDie(sides int) int {
	return int(dr.rng.Next()%uint64(sides)) + 1
}

func (dr *DiceRoller) RollMultiple(sides, count int) []int {
	results := make([]int, count)
	for i := 0; i < count; i++ {
		results[i] = dr.RollDie(sides)
	}
	return results
}

func (dr *DiceRoller) RollAndPrint(sides, count int) {
	fmt.Printf("Rolling a d%d dice %d times:\n", sides, count)
	for i, result := range dr.RollMultiple(sides, count) {
		fmt.Printf("Roll %d: %d\n", i+1, result)
	}
}
