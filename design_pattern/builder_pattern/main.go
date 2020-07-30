package main

import (
    "fmt"
    mario "github.com/mats9693/study/design_pattern/builder_pattern/meal"
)

func main() {
    director := mario.Director{Builders: []mario.BuilderI{&mario.Builder1{}, &mario.Builder2{}}}

    director.SetBuilderIndex(0)
    setMealIns := director.MakeSetMeal()
    fmt.Println(setMealIns.String())

    director.SetBuilderIndex(1)
    setMealIns = director.MakeSetMeal()
    fmt.Println(setMealIns.String())
}
