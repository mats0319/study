package main

import mario "github.com/mats9693/study/books/978-7-302-16206-3/decorator_pattern/finery"

func main() {
    var (
        person mario.Finery = &mario.Person{Name: "Mario"}
        dec1                = &mario.LongSleeves{}
        dec2                = &mario.Trousers{}
        dec3                = &mario.Shoes{}
    )

    // 这里就不做交互了，因为不是重点。做了交互有一种本末倒置的感觉
    dec1.Dec.Decorate(person)
    dec2.Dec.Decorate(dec1)
    dec3.Dec.Decorate(dec2)
    dec3.Show()

    return
}
