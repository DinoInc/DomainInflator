package Engine

type DeviationContext int64

const (
	PropertyIdentifier DeviationContext = 0
	EnumIdentifier     DeviationContext = 1
)

type Deviation struct {
	Original    string
	Replacement string
	Context     DeviationContext
}
