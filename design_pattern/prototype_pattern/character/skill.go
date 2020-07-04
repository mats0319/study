package mario

import (
    "fmt"
    "strconv"
)

type skill struct {
	SkillName        string
	SkillType        string
	SkillTime        int
	SkillCoefficient float64
	SkillNotes       string
}

const (
    defaultSkillName = "default skill"
	DamageSkill  = "伤害"
)

var defaultSkill = &skill{
    SkillName:        defaultSkillName,
    SkillType:        "",
    SkillTime:        0,
    SkillCoefficient: 0,
    SkillNotes:       "",
}

func CreateSkill(name, typ string, tim int, coefficient float64, notes string) *skill {
    if notes == "" {
        notes = "无"
    }

    return &skill{
        SkillName:        name,
        SkillType:        typ,
        SkillTime:        tim,
        SkillCoefficient: coefficient,
        SkillNotes:       notes,
    }
}

func (s *skill) String() string {
    if s.SkillName == defaultSkillName {
        return "无"
    }

    var skillTimeStr string
    if s.SkillTime == 0 {
        skillTimeStr = "瞬发"
    } else {
        skillTimeStr = strconv.Itoa(s.SkillTime) + "秒"
    }

    return fmt.Sprintf("技能名称：%s, 类型：%s, 持续时间：%s, 系数：%0.2f倍基础伤害, 备注：%s\n",
        s.SkillName, s.SkillType, skillTimeStr, s.SkillCoefficient, s.SkillNotes)
}

func (s *skill) Clone() *skill {
    return &skill{
        SkillName:        s.SkillName,
        SkillType:        s.SkillType,
        SkillTime:        s.SkillTime,
        SkillCoefficient: s.SkillCoefficient,
        SkillNotes:       s.SkillNotes,
    }
}
