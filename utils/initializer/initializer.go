package initializer

type PriorityLevel int

const (
	VeryHighPriority PriorityLevel = iota
	HighPriority
	MidPriority
	LowPriority
)

type Initializer interface {
	Initialize() func()
}

var initializers [4][]Initializer

func Register(in Initializer, pl ...PriorityLevel) {
	if len(pl) > 0 {
		initializers[pl[0]] = append(initializers[pl[0]], in)
	}
	initializers[MidPriority] = append(initializers[MidPriority], in)
}

func Initialize() func() {
	var finalizers []func()
	for i := range initializers {
		for _, j := range initializers[i] {
			finalizers = append(finalizers, j.Initialize())
		}
	}

	return func() {
		for i := range finalizers {
			if finalizers[i] != nil {
				finalizers[i]()
			}
		}
	}
}
