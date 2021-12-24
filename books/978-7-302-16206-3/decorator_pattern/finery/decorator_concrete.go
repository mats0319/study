package mario

import "fmt"

type LongSleeves struct {
	Dec Decorator
}

func (l *LongSleeves) Show() {
	l.Dec.Show()
	fmt.Printf(" 长袖外套 ")

	return
}

type ShortSleeves struct {
	Dec Decorator
}

func (s *ShortSleeves) Show() {
	s.Dec.Show()
	fmt.Printf(" 短袖外套 ")

	return
}

type Trousers struct {
	Dec Decorator
}

func (t *Trousers) Show() {
	t.Dec.Show()
	fmt.Printf(" 裤子 ")

	return
}

type Shoes struct {
	Dec Decorator
}

func (s *Shoes) Show() {
	s.Dec.Show()
	fmt.Printf(" 鞋 ")

	return
}
