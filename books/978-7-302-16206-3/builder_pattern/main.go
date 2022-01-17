package main

import (
    "fmt"
    "github.com/mats9693/study/books/978-7-302-16206-3/builder_pattern/meal"
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
