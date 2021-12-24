package mario

import "errors"

type ActivityNormal struct {
	ActivityDetails
}

func (a *ActivityNormal) CalculateSummary() (float64, error) {
	return a.Summary, nil
}

type ActivityPercent struct {
	ActivityDetails
}

func (a *ActivityPercent) CalculateSummary() (float64, error) {
	return a.Summary * a.ActivityParam1, nil
}

type ActivityReturn struct {
	ActivityDetails
}

func (a *ActivityReturn) CalculateSummary() (float64, error) {
	if a.ActivityParam1 == 0 {
		return 0, errors.New("请输入满减的目标金额")
	}

	return a.Summary - float64(int(a.Summary/a.ActivityParam1))*a.ActivityParam2, nil
}
