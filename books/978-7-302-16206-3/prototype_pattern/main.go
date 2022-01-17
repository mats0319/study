package main

import (
    mario "github.com/mats9693/study/books/978-7-302-16206-3/prototype_pattern/character"
    "log"
)

func main() {
    simpleAttack := mario.CreateSkill("Attack", mario.DamageSkill, 0, 1.0, "")
    continuousAttack := mario.CreateSkill("Attack-combo", mario.DamageSkill, 3, 0.8, "")

    charOne := mario.CreateCharacterWithUid("Mario")

    {
        charTwo := charOne.Clone()
        charTwo.ModifyName("Mario second")
        charTwo.UpdateSkill(simpleAttack)

        log.Printf("compare shallow clone res, after modify name and skills:\n%s\n%s\n",
            charOne.String(), charTwo.String())
        log.Println("modify skill info is dependent.")
    }

    log.Println("-------")

    {
        charThree := charOne.CloneDeep()
        charThree.ModifyName("Mario three")
        charThree.UpdateSkill(continuousAttack)

        log.Printf("compare deep clone res, after modify name and skills:\n%s\n%s\n",
            charOne.String(), charThree.String())
        log.Println("modify skill info is independent.")
    }

    return
}
