package main

import (
    "fmt"
    mario "github.com/mats9693/study/books/978-7-302-16206-3/strategy_pattern/activity"
    "log"
)

func main() {
    sc := &mario.StrategyContext{}
    ad := mario.ActivityDetails{}
    activityName := ""

    log.Printf("请输入金额与折扣活动（形如‘100 p 0.8 0’，4个参数，以空格分开，支持的折扣活动：原价-normal，打折-percent，满减-return）：\n")
    if n, err := fmt.Scanln(&ad.Summary, &activityName, &ad.ActivityParam1, &ad.ActivityParam2); (n != 3 && n != 4) || err != nil {
        log.Printf("表达式解析错误！\t成功解析数：%d，错误：%v\n使用默认参数：100 n 0 0\n", n, err)
        ad.Summary = 100
        activityName = "n"
        ad.ActivityParam1 = 0
        ad.ActivityParam2 = 0
    }

    var (
        result float64
        err    error
    )
    switch activityName {
    case "n":
        result, err = sc.CalculateSummary(&mario.ActivityNormal{ActivityDetails: ad})
        activityName = "原价"
    case "p":
        result, err = sc.CalculateSummary(&mario.ActivityPercent{ActivityDetails: ad})
        activityName = "打折"
    case "r":
        result, err = sc.CalculateSummary(&mario.ActivityReturn{ActivityDetails: ad})
        activityName = "满减"
    default:
        log.Fatalln("未知的折扣活动：", activityName)
    }

    if err != nil {
        log.Fatalln("计算错误：", err)
    }

    log.Printf("\n原价：%f\t活动：%s\n现价：%f\n", ad.Summary, activityName, result)

    return
}
