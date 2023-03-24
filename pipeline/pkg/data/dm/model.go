package dm

type AvatarExcel struct {
	// ScriptDataPathHashSuffix     int64           `json:"scriptDataPathHashSuffix"`
	// ScriptDataPathHashPre        int64           `json:"scriptDataPathHashPre"`
	// SideIconName                 string        `json:"sideIconName"`
	// ChargeEfficiency             float64       `json:"chargeEfficiency"`
	// CombatConfigHashSuffix       int64           `json:"combatConfigHashSuffix"`
	// CombatConfigHashPre          int64           `json:"combatConfigHashPre"`
	// ManekinPathHashSuffix        int64           `json:"manekinPathHashSuffix"`
	// ManekinPathHashPre           int64           `json:"manekinPathHashPre"`
	// GachaCardNameHashPre         int64           `json:"gachaCardNameHashPre"`
	// CutsceneShow                 string        `json:"cutsceneShow"`
	// StaminaRecoverSpeed          float64       `json:"staminaRecoverSpeed"`
	// CandSkillDepotIds            []interface{} `json:"candSkillDepotIds"`
	// ManekinJSONConfigHashSuffix  int64           `json:"manekinJsonConfigHashSuffix"`
	// ManekinJSONConfigHashPre     int64           `json:"manekinJsonConfigHashPre"`
	// ManekinMotionConfig          int64           `json:"manekinMotionConfig"`
	// DescTextMapHash              int64           `json:"descTextMapHash"`
	// AvatarIdentityType string `json:"avatarIdentityType"`
	// AvatarPromoteRewardLevelList []int64         `json:"avatarPromoteRewardLevelList"`
	// AvatarPromoteRewardIDList    []int64         `json:"avatarPromoteRewardIdList"`
	// FeatureTagGroupID            int64           `json:"featureTagGroupID"`
	// InfoDescTextMapHash          int64           `json:"infoDescTextMapHash"`
	// PrefabPathRagdollHashSuffix    int64    `json:"prefabPathRagdollHashSuffix"`
	// PrefabPathRagdollHashPre       int64    `json:"prefabPathRagdollHashPre"`
	// AnimatorConfigPathHashSuffix   int64  `json:"animatorConfigPathHashSuffix"`
	// PrefabPathHashSuffix           int64  `json:"prefabPathHashSuffix"`
	// PrefabPathHashPre              int64    `json:"prefabPathHashPre"`
	// PrefabPathRemoteHashSuffix     int64  `json:"prefabPathRemoteHashSuffix"`
	// PrefabPathRemoteHashPre        int64    `json:"prefabPathRemoteHashPre"`
	// ControllerPathHashSuffix       int64  `json:"controllerPathHashSuffix"`
	// ControllerPathHashPre          int64    `json:"controllerPathHashPre"`
	// ControllerPathRemoteHashSuffix int64  `json:"controllerPathRemoteHashSuffix"`
	// ControllerPathRemoteHashPre    int64    `json:"controllerPathRemoteHashPre"`
	// LODPatternName                 string `json:"LODPatternName"`
	// InitialWeapon   int64     `json:"initialWeapon"`
	// Critical        float64 `json:"critical"`
	// CriticalHurt    float64 `json:"criticalHurt"`
	BodyType     string `json:"bodyType"`
	SkillDepotID int64  `json:"skillDepotId"`
	IconName     string `json:"iconName"`
	QualityType  string `json:"qualityType"`
	WeaponType   string `json:"weaponType"`
	// ImageName       string          `json:"imageName"`
	AvatarPromoteID int64           `json:"avatarPromoteId"`
	HpBase          float64         `json:"hpBase"`
	AttackBase      float64         `json:"attackBase"`
	DefenseBase     float64         `json:"defenseBase"`
	PropGrowCurves  []PropGrowCurve `json:"propGrowCurves"`
	ID              int64           `json:"id"`
	NameTextMapHash int64           `json:"nameTextMapHash"`
	UseType         string          `json:"useType,omitempty"`
}

type PropGrowCurve struct {
	Type      string `json:"type"`
	GrowCurve string `json:"growCurve"`
}

type TextMap map[int64]string

type AvatarSkillDepot struct {
	ID          int64   `json:"id"`
	EnergySkill int64   `json:"energySkill"`
	Skills      []int64 `json:"skills"`
	// SubSkills               []int64    `json:"subSkills"`
	// ExtraAbilities          []string `json:"extraAbilities"`
	// Talents                 []int64    `json:"talents"`
	// TalentStarName          string   `json:"talentStarName"`
	// InherentProudSkillOpens []struct {
	// } `json:"inherentProudSkillOpens"`
	// SkillDepotAbilityGroup string `json:"skillDepotAbilityGroup"`
	// LeaderTalent           int64    `json:"leaderTalent,omitempty"`
}
type AvatarSkillExcel struct {
	ID           int64  `json:"id"`
	CostElemType string `json:"costElemType,omitempty"`
	// NameTextMapHash    int64     `json:"nameTextMapHash"`
	// AbilityName        string    `json:"abilityName"`
	// DescTextMapHash    int64       `json:"descTextMapHash"`
	// SkillIcon          string    `json:"skillIcon"`
	// CdTime             float64   `json:"cdTime,omitempty"`
	CostElemVal float64 `json:"costElemVal,omitempty"`
	// MaxChargeNum       int64       `json:"maxChargeNum"`
	// TriggerID          int64       `json:"triggerID,omitempty"`
	// LockShape          string    `json:"lockShape"`
	// LockWeightParams   []float64 `json:"lockWeightParams"`
	// IsAttackCameraLock bool      `json:"isAttackCameraLock"`
	// BuffIcon           string    `json:"buffIcon"`
	// GlobalValueKey     string    `json:"globalValueKey"`
	// CostStamina        float64   `json:"costStamina,omitempty"`
}

type AvatarFetterInfo struct {
	AvatarAssocType string `json:"avatarAssocType"`
	AvatarId        int64  `json:"avatarId"`
	FetterID        int64  `json:"fetterId"`
}

type AvatarPromote struct {
	AvatarPromoteID int64     `json:"avatarPromoteId"`
	AddProps        []AddProp `json:"addProps"`
	// PromoteLevel    int64     `json:"promoteLevel,omitempty"`
	// PromoteAudio    string `json:"promoteAudio"`
	// CostItems       []struct {
	// } `json:"costItems"`
	UnlockMaxLevel int64 `json:"unlockMaxLevel"`
	// ScoinCost           int64 `json:"scoinCost,omitempty"`
	// RequiredPlayerLevel int64 `json:"requiredPlayerLevel,omitempty"`
}
type AddProp struct {
	PropType string  `json:"propType"`
	Value    float64 `json:"value"`
}

type WeaponExcel struct {
	WeaponType string `json:"weaponType"`
	RankLevel  int64  `json:"rankLevel"`
	// WeaponBaseExp int64    `json:"weaponBaseExp"`
	SkillAffix []int64 `json:"skillAffix"`
	WeaponProp []struct {
		PropType  string  `json:"propType,omitempty"`
		InitValue float64 `json:"initValue,omitempty"`
		Type      string  `json:"type"`
	} `json:"weaponProp"`
	// AwakenTexture              string        `json:"awakenTexture"`
	// AwakenLightMapTexture      string        `json:"awakenLightMapTexture"`
	// AwakenIcon                 string        `json:"awakenIcon"`
	WeaponPromoteID int64 `json:"weaponPromoteId"`
	// StoryID                    int64           `json:"storyId"`
	// AwakenCosts                []interface{} `json:"awakenCosts"`
	// GachaCardNameHashSuffix    int64         `json:"gachaCardNameHashSuffix"`
	// DestroyRule                string        `json:"destroyRule"`
	// DestroyReturnMaterial      []int64         `json:"destroyReturnMaterial"`
	// DestroyReturnMaterialCount []int64         `json:"destroyReturnMaterialCount"`
	ID              int64 `json:"id"`
	NameTextMapHash int64 `json:"nameTextMapHash"`
	// DescTextMapHash            int64         `json:"descTextMapHash"`
	Icon     string `json:"icon"`
	ItemType string `json:"itemType"`
	// Weight                     int64           `json:"weight"`
	Rank int64 `json:"rank"`
	// GadgetID                   int64           `json:"gadgetId"`
}

type WeaponCurve struct {
}

type WeaponPromoteConfig struct {
	WeaponPromoteID int64 `json:"weaponPromoteId"`
	// CostItems       []struct {
	// } `json:"costItems"`
	AddProps       []AddProp
	UnlockMaxLevel int64 `json:"unlockMaxLevel"`
	PromoteLevel   int64 `json:"promoteLevel,omitempty"`
	// RequiredPlayerLevel int64 `json:"requiredPlayerLevel,omitempty"`
	// CoinCost            int64 `json:"coinCost,omitempty"`
}
