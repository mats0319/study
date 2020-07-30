package mario

import (
	"errors"
	"fmt"
)

type ActivityTemplate struct {
	Name   string
	Awards string
}

var _ Activity = (*ActivityTemplate)(nil)

func (a *ActivityTemplate) Start(playerName string) string {
	return fmt.Sprintf("Player: %s select game: %s, enjoy a good time! :)\n", playerName, a.Name)
}

func (a *ActivityTemplate) Failed(playerName string, percent float64) string {
	return fmt.Sprintf("Player %s failed, with %0.2f%% of %s game.\n", playerName, percent, a.Name)
}

func (a *ActivityTemplate) Success(playerName string) string {
	return fmt.Sprintf("Player %s passed the game: %s!\n", playerName, a.Name)
}

func (a *ActivityTemplate) Reward() string {
	return fmt.Sprintf("You can get %s as award from %s game.\n", a.Awards, a.Name)
}

func (a *ActivityTemplate) activityStage(playerName string, passGame bool, percent ...float64) (string, error) {
    res := ""
    res += "\n" + a.Start(playerName)
    if passGame {
        res += a.Success(playerName)
    } else {
        if len(percent) != 1 {
            return "", errors.New("invalid input param")
        }

        res += a.Failed(playerName, percent[0])
    }
    res += a.Reward()

    return res, nil
}

type ActivityInstance struct {
	PlayerName string
	*ActivityTemplate
}

func (a *ActivityInstance) Play(passGame bool, percent ...float64) (string, error) {
    return  a.activityStage(a.PlayerName, passGame, percent...)
}


