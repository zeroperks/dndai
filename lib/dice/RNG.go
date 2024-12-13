package dice

type RNG interface {
	Next() uint64
}
