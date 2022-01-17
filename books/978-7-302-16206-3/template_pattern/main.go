package main

import (
    mario "github.com/mats9693/study/books/978-7-302-16206-3/template_pattern/activity"
    "log"
)

func main() {
    var (
        activityIns = &mario.ActivityTemplate{
            Name:   "Limit Challenge",
            Awards: "Money, Exp, Sword",
        }

        marioIns = &mario.ActivityInstance{
            PlayerName:       "Mario",
            ActivityTemplate: activityIns,
        }
        phoenixIns = &mario.ActivityInstance{
            PlayerName:       "Phoenix",
            ActivityTemplate: activityIns,
        }
    )

    res, err := marioIns.Play(true)
    if err != nil {
        log.Println("limit challenge failed, error:", err)
    } else {
        log.Println(res)
    }

    res, err = phoenixIns.Play(false, 80)
    if err != nil {
        log.Println("normal instance failed, error:", err)
    } else {
        log.Println(res)
    }
}
