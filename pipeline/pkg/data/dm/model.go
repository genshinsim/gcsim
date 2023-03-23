package dm

type AvatarExcel struct {
	// ScriptDataPathHashSuffix     int           `json:"scriptDataPathHashSuffix"`
	// ScriptDataPathHashPre        int           `json:"scriptDataPathHashPre"`
	// SideIconName                 string        `json:"sideIconName"`
	// ChargeEfficiency             float64       `json:"chargeEfficiency"`
	// CombatConfigHashSuffix       int           `json:"combatConfigHashSuffix"`
	// CombatConfigHashPre          int           `json:"combatConfigHashPre"`
	// ManekinPathHashSuffix        int           `json:"manekinPathHashSuffix"`
	// ManekinPathHashPre           int           `json:"manekinPathHashPre"`
	// GachaCardNameHashPre         int           `json:"gachaCardNameHashPre"`
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
	BodyType        string          `json:"bodyType"`
	SkillDepotID    int             `json:"skillDepotId"`
	IconName        string          `json:"iconName"`
	QualityType     string          `json:"qualityType"`
	WeaponType      string          `json:"weaponType"`
	ImageName       string          `json:"imageName"`
	AvatarPromoteID int             `json:"avatarPromoteId"`
	HpBase          float64         `json:"hpBase"`
	AttackBase      float64         `json:"attackBase"`
	DefenseBase     float64         `json:"defenseBase"`
	PropGrowCurves  []PropGrowCurve `json:"propGrowCurves"`
	ID              int             `json:"id"`
	NameTextMapHash int             `json:"nameTextMapHash"`
	UseType         string          `json:"useType,omitempty"`
}

type PropGrowCurve struct {
	Type      string `json:"type"`
	GrowCurve string `json:"growCurve"`
}

type TextMap map[int]string

type AvatarSkillDepot struct {
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
type AvatarSkillExcel struct {
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

type FetterInfo struct {
	AvatarAssocType string `json:"avatarAssocType"`
	AvatarId        int    `json:"avatarId"`
}

type AvatarPromote struct {
	AvatarPromoteID int       `json:"avatarPromoteId"`
	AddProps        []AddProp `json:"addProps"`
	PromoteLevel    int       `json:"promoteLevel,omitempty"`
	// PromoteAudio    string `json:"promoteAudio"`
	// CostItems       []struct {
	// } `json:"costItems"`
	UnlockMaxLevel int `json:"unlockMaxLevel"`
	// ScoinCost           int `json:"scoinCost,omitempty"`
	// RequiredPlayerLevel int `json:"requiredPlayerLevel,omitempty"`
}
type AddProp struct {
	PropType string  `json:"propType"`
	Value    float64 `json:"value"`
}

type WeaponExcel struct {
	WeaponType string `json:"weaponType"`
	RankLevel  int    `json:"rankLevel"`
	// WeaponBaseExp int    `json:"weaponBaseExp"`
	SkillAffix []int `json:"skillAffix"`
	WeaponProp []struct {
		PropType  string  `json:"propType,omitempty"`
		InitValue float64 `json:"initValue,omitempty"`
		Type      string  `json:"type"`
	} `json:"weaponProp"`
	// AwakenTexture              string        `json:"awakenTexture"`
	// AwakenLightMapTexture      string        `json:"awakenLightMapTexture"`
	// AwakenIcon                 string        `json:"awakenIcon"`
	WeaponPromoteID int `json:"weaponPromoteId"`
	// StoryID                    int           `json:"storyId"`
	// AwakenCosts                []interface{} `json:"awakenCosts"`
	// GachaCardNameHashSuffix    int64         `json:"gachaCardNameHashSuffix"`
	// DestroyRule                string        `json:"destroyRule"`
	// DestroyReturnMaterial      []int         `json:"destroyReturnMaterial"`
	// DestroyReturnMaterialCount []int         `json:"destroyReturnMaterialCount"`
	ID              int   `json:"id"`
	NameTextMapHash int64 `json:"nameTextMapHash"`
	// DescTextMapHash            int64         `json:"descTextMapHash"`
	Icon     string `json:"icon"`
	ItemType string `json:"itemType"`
	// Weight                     int           `json:"weight"`
	Rank int `json:"rank"`
	// GadgetID                   int           `json:"gadgetId"`
}

type WeaponCurve struct {
}

type WeaponPromoteConfig struct {
	WeaponPromoteID int `json:"weaponPromoteId"`
	// CostItems       []struct {
	// } `json:"costItems"`
	AddProps       []AddProp
	UnlockMaxLevel int `json:"unlockMaxLevel"`
	PromoteLevel   int `json:"promoteLevel,omitempty"`
	// RequiredPlayerLevel int `json:"requiredPlayerLevel,omitempty"`
	// CoinCost            int `json:"coinCost,omitempty"`
}
