package main

import (
	"encoding/json"
	"log"
)

type Avatars []struct {
	// ScriptDataPathHashSuffix     int           `json:"scriptDataPathHashSuffix"`
	// ScriptDataPathHashPre        int           `json:"scriptDataPathHashPre"`
	// SideIconName                 string        `json:"sideIconName"`
	// ChargeEfficiency             float64       `json:"chargeEfficiency"`
	// CombatConfigHashSuffix       int           `json:"combatConfigHashSuffix"`
	// CombatConfigHashPre          int           `json:"combatConfigHashPre"`
	// ManekinPathHashSuffix        int           `json:"manekinPathHashSuffix"`
	// ManekinPathHashPre           int           `json:"manekinPathHashPre"`
	// Mjgngjhbagi                  int64         `json:"MJGNGJHBAGI"`
	// GachaCardNameHashPre         int           `json:"gachaCardNameHashPre"`
	// Cbfoekcenea                  int64         `json:"CBFOEKCENEA"`
	// Pagadeakhac                  int           `json:"PAGADEAKHAC"`
	// CutsceneShow                 string        `json:"cutsceneShow"`
	// StaminaRecoverSpeed          float64       `json:"staminaRecoverSpeed"`
	// CandSkillDepotIds            []interface{} `json:"candSkillDepotIds"`
	// ManekinJSONConfigHashSuffix  int           `json:"manekinJsonConfigHashSuffix"`
	// ManekinJSONConfigHashPre     int           `json:"manekinJsonConfigHashPre"`
	// ManekinMotionConfig          int           `json:"manekinMotionConfig"`
	// DescTextMapHash              int           `json:"descTextMapHash"`
	// AvatarIdentityType string `json:"avatarIdentityType"`
	// AvatarPromoteRewardLevelList []int         `json:"avatarPromoteRewardLevelList"`
	// AvatarPromoteRewardIDList    []int         `json:"avatarPromoteRewardIdList"`
	// FeatureTagGroupID            int           `json:"featureTagGroupID"`
	// InfoDescTextMapHash          int           `json:"infoDescTextMapHash"`
	// PrefabPathRagdollHashSuffix    int    `json:"prefabPathRagdollHashSuffix"`
	// PrefabPathRagdollHashPre       int    `json:"prefabPathRagdollHashPre"`
	// AnimatorConfigPathHashSuffix   int64  `json:"animatorConfigPathHashSuffix"`
	// Njiekklklld                    int    `json:"NJIEKKLKLLD"`
	// PrefabPathHashSuffix           int64  `json:"prefabPathHashSuffix"`
	// PrefabPathHashPre              int    `json:"prefabPathHashPre"`
	// PrefabPathRemoteHashSuffix     int64  `json:"prefabPathRemoteHashSuffix"`
	// PrefabPathRemoteHashPre        int    `json:"prefabPathRemoteHashPre"`
	// ControllerPathHashSuffix       int64  `json:"controllerPathHashSuffix"`
	// ControllerPathHashPre          int    `json:"controllerPathHashPre"`
	// ControllerPathRemoteHashSuffix int64  `json:"controllerPathRemoteHashSuffix"`
	// ControllerPathRemoteHashPre    int    `json:"controllerPathRemoteHashPre"`
	// LODPatternName                 string `json:"LODPatternName"`
	// InitialWeapon   int     `json:"initialWeapon"`
	// Critical        float64 `json:"critical"`
	// CriticalHurt    float64 `json:"criticalHurt"`
	SkillDepotID    int            `json:"skillDepotId"`
	BodyType        string         `json:"bodyType"`
	IconName        string         `json:"iconName"`
	QualityType     string         `json:"qualityType"`
	WeaponType      string         `json:"weaponType"`
	ImageName       string         `json:"imageName"`
	AvatarPromoteID int            `json:"avatarPromoteId"`
	HpBase          float64        `json:"hpBase"`
	AttackBase      float64        `json:"attackBase"`
	DefenseBase     float64        `json:"defenseBase"`
	PropGrowCurves  PropGrowCurves `json:"propGrowCurves"`
	ID              int            `json:"id"`
	NameTextMapHash int            `json:"nameTextMapHash"`
	UseType         string         `json:"useType,omitempty"`
}

type PropGrowCurves []struct {
	Type string `json:"type"`

	GrowCurve string `json:"growCurve"`
}

func getAvatarArray() (Avatars, []int) {
	avatarDataJson, err := fetchJsonFromUrl("https://raw.githubusercontent.com/Dimbreath/GenshinData/master/ExcelBinOutput/AvatarExcelConfigData.json")
	if err != nil {
		log.Fatal(err)
	}
	var avatars Avatars
	if err := json.Unmarshal([]byte(avatarDataJson), &avatars); err != nil {
		log.Fatal(err)
	}
	// remove testing/invalid chars
	var filteredAvatars Avatars
	var textMapIds []int
	for _, avatar := range avatars {
		if avatar.UseType == "AVATAR_FORMAL" {
			filteredAvatars = append(filteredAvatars, avatar)
			textMapIds = append(textMapIds, avatar.NameTextMapHash)
		}
	}
	return filteredAvatars, textMapIds

}
