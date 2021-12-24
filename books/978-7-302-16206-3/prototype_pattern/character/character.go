package mario

import (
	"fmt"
)

var uid = 0

func uidAddOne() int {
	uid++
	return uid
}

type character struct {
	Id    int
	Name  string
	Level int
	Skill *skill
}

func CreateCharacterWithUid(name string) *character {
	return &character{
		Id:    uidAddOne(),
		Name:  name,
		Level: 1,
		Skill: defaultSkill,
	}
}

func (c *character) String() string {
	return fmt.Sprintf("\n人物卡：\nID：%d, 昵称：%s, 等级：%d, 掌握技能：\n  %s",
		c.Id, c.Name, c.Level, c.Skill.String())
}

func (c *character) ModifyName(newName string) {
	c.Name = newName
}

func (c *character) UpdateSkill(skill *skill) {
	c.Skill.SkillName = skill.SkillName
	c.Skill.SkillType = skill.SkillType
	c.Skill.SkillTime = skill.SkillTime
	c.Skill.SkillCoefficient = skill.SkillCoefficient
	c.Skill.SkillNotes = skill.SkillNotes
}

func (c *character) Clone() *character {
	return &character{
		Id:    uidAddOne(),
		Name:  c.Name,
		Level: 1,
		Skill: c.Skill,
	}
}

func (c *character) CloneDeep() *character {
	return &character{
		Id:    uidAddOne(),
		Name:  c.Name,
		Level: 1,
		Skill: c.Skill.Clone(),
	}
}
