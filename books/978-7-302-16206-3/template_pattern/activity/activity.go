package mario

type Activity interface {
	Start(playerName string) string
	Failed(playerName string, percent float64) string
	Success(playerName string) string
	Reward() string
}
