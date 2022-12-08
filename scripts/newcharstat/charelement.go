package main

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

type SkillDepot []struct {
	ID          int `json:"id"`
	EnergySkill int `json:"energySkill"`
	// Skills                  []int    `json:"skills"`
	// SubSkills               []int    `json:"subSkills"`
	// ExtraAbilities          []string `json:"extraAbilities"`
	// Talents                 []int    `json:"talents"`
	// TalentStarName          string   `json:"talentStarName"`
	// InherentProudSkillOpens []struct {
	// } `json:"inherentProudSkillOpens"`
	// SkillDepotAbilityGroup string `json:"skillDepotAbilityGroup"`
	// LeaderTalent           int    `json:"leaderTalent,omitempty"`
}
type AvatarSkillInfo []struct {
	ID           int    `json:"id"`
	CostElemType string `json:"costElemType,omitempty"`
	// NameTextMapHash    int64     `json:"nameTextMapHash"`
	// AbilityName        string    `json:"abilityName"`
	// DescTextMapHash    int       `json:"descTextMapHash"`
	// SkillIcon          string    `json:"skillIcon"`
	// CdTime             float64   `json:"cdTime,omitempty"`
	// CostElemVal        float64   `json:"costElemVal,omitempty"`
	// MaxChargeNum       int       `json:"maxChargeNum"`
	// TriggerID          int       `json:"triggerID,omitempty"`
	// LockShape          string    `json:"lockShape"`
	// LockWeightParams   []float64 `json:"lockWeightParams"`
	// IsAttackCameraLock bool      `json:"isAttackCameraLock"`
	// BuffIcon           string    `json:"buffIcon"`
	// GlobalValueKey     string    `json:"globalValueKey"`
	// CostStamina        float64   `json:"costStamina,omitempty"`
}

func getAvatarElementMap() map[int]string {
	skillDepot := getJsonFromFile[SkillDepot]("../ExcelBinOutput/AvatarSkillDepotExcelConfigData.json")

	avatarSkillInfo := getJsonFromFile[AvatarSkillInfo]("../ExcelBinOutput/AvatarSkillExcelConfigData.json")

	// reshape avatarSkillInfo to map of energyskillID to CostElemType
	energySkillMap := make(map[int]string)
	for _, v := range avatarSkillInfo {
		energySkillMap[v.ID] = v.CostElemType
	}

	// reshape skillDepot to map of skilldepotID to CostElemType
	elementMap := make(map[int]string)
	for _, skill := range skillDepot {
		if skill.EnergySkill == 0 {
			continue
		}
		elementMap[skill.ID] = energySkillMap[skill.EnergySkill]
	}

	return elementMap
}

func convertElement(element string) attributes.Element {
	switch element {
	case "Fire":
		return attributes.Pyro
	case "Water":
		return attributes.Hydro
	case "Wind":
		return attributes.Anemo
	case "Ice":
		return attributes.Cryo
	case "Electric":
		return attributes.Electro
	case "Rock":
		return attributes.Geo
	case "Grass":
		return attributes.Dendro
	default:
		return attributes.UnknownElement
	}
}
