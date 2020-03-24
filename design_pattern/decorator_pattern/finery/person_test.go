package mario

func ExamplePerson_Show() {
	var (
		person Finery = &Person{Name: "Mario"}
		dec1 = &LongSleeves{}
		dec2 = &Trousers{}
		dec3 = &Shoes{}
	)

	dec1.Dec.Decorate(person)
	dec2.Dec.Decorate(dec1)
	dec3.Dec.Decorate(dec2)
	dec3.Show()

	return

	// Output:
	// 我叫: Mario, 我穿着: 长袖外套  裤子  鞋
}
