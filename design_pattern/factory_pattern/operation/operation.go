package mario

type Operation interface {
	CalculateResult() (float64, error)
}

type OperationNumber struct {
	NumberA float64
	NumberB float64
}
