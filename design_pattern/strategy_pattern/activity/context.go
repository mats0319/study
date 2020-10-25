package mario

type StrategyContext struct {
	Activity
}

func (s *StrategyContext) CalculateSummary(activity Activity) (float64, error) {
	s.Activity = activity

	return s.Activity.CalculateSummary()
}
