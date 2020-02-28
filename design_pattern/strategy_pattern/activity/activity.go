package mario

type Activity interface {
	CalculateSummary() (float64, error)
}

type ActivityDetails struct {
	Summary float64
	ActivityParam1 float64
	ActivityParam2 float64
}
