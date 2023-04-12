package dm

type AvatarExcel struct {
	// ScriptDataPathHashSuffix     int32           `json:"scriptDataPathHashSuffix"`
	// ScriptDataPathHashPre        int32           `json:"scriptDataPathHashPre"`
	// SideIconName                 string        `json:"sideIconName"`
	// ChargeEfficiency             float64       `json:"chargeEfficiency"`
	// CombatConfigHashSuffix       int32           `json:"combatConfigHashSuffix"`
	// CombatConfigHashPre          int32           `json:"combatConfigHashPre"`
	// ManekinPathHashSuffix        int32           `json:"manekinPathHashSuffix"`
	// ManekinPathHashPre           int32           `json:"manekinPathHashPre"`
	// GachaCardNameHashPre         int32           `json:"gachaCardNameHashPre"`
	// CutsceneShow                 string        `json:"cutsceneShow"`
	// StaminaRecoverSpeed          float64       `json:"staminaRecoverSpeed"`
	// CandSkillDepotIds            []interface{} `json:"candSkillDepotIds"`
	// ManekinJSONConfigHashSuffix  int32           `json:"manekinJsonConfigHashSuffix"`
	// ManekinJSONConfigHashPre     int32           `json:"manekinJsonConfigHashPre"`
	// ManekinMotionConfig          int32           `json:"manekinMotionConfig"`
	// DescTextMapHash              int32           `json:"descTextMapHash"`
	// AvatarIdentityType string `json:"avatarIdentityType"`
	// AvatarPromoteRewardLevelList []int32         `json:"avatarPromoteRewardLevelList"`
	// AvatarPromoteRewardIDList    []int32         `json:"avatarPromoteRewardIdList"`
	// FeatureTagGroupID            int32           `json:"featureTagGroupID"`
	// InfoDescTextMapHash          int32           `json:"infoDescTextMapHash"`
	// PrefabPathRagdollHashSuffix    int32    `json:"prefabPathRagdollHashSuffix"`
	// PrefabPathRagdollHashPre       int32    `json:"prefabPathRagdollHashPre"`
	// AnimatorConfigPathHashSuffix   int32  `json:"animatorConfigPathHashSuffix"`
	// PrefabPathHashSuffix           int32  `json:"prefabPathHashSuffix"`
	// PrefabPathHashPre              int32    `json:"prefabPathHashPre"`
	// PrefabPathRemoteHashSuffix     int32  `json:"prefabPathRemoteHashSuffix"`
	// PrefabPathRemoteHashPre        int32    `json:"prefabPathRemoteHashPre"`
	// ControllerPathHashSuffix       int32  `json:"controllerPathHashSuffix"`
	// ControllerPathHashPre          int32    `json:"controllerPathHashPre"`
	// ControllerPathRemoteHashSuffix int32  `json:"controllerPathRemoteHashSuffix"`
	// ControllerPathRemoteHashPre    int32    `json:"controllerPathRemoteHashPre"`
	// LODPatternName                 string `json:"LODPatternName"`
	// InitialWeapon   int32     `json:"initialWeapon"`
	// Critical        float64 `json:"critical"`
	// CriticalHurt    float64 `json:"criticalHurt"`
	BodyType     string `json:"bodyType"`
	SkillDepotID int32  `json:"skillDepotId"`
	IconName     string `json:"iconName"`
	QualityType  string `json:"qualityType"`
	WeaponType   string `json:"weaponType"`
	// ImageName       string          `json:"imageName"`
	AvatarPromoteID int32           `json:"avatarPromoteId"`
	HpBase          float64         `json:"hpBase"`
	AttackBase      float64         `json:"attackBase"`
	DefenseBase     float64         `json:"defenseBase"`
	PropGrowCurves  []PropGrowCurve `json:"propGrowCurves"`
	ID              int32           `json:"id"`
	// NameTextMapHash int32           `json:"nameTextMapHash"`
	UseType string `json:"useType,omitempty"`
}

type PropGrowCurve struct {
	Type      string `json:"type"`
	GrowCurve string `json:"growCurve"`
}

type TextMap map[int32]string

type AvatarSkillDepot struct {
	ID          int32   `json:"id"`
	EnergySkill int32   `json:"energySkill"`
	Skills      []int32 `json:"skills"`
	// SubSkills               []int32    `json:"subSkills"`
	// ExtraAbilities          []string `json:"extraAbilities"`
	// Talents                 []int32    `json:"talents"`
	// TalentStarName          string   `json:"talentStarName"`
	// InherentProudSkillOpens []struct {
	// } `json:"inherentProudSkillOpens"`
	// SkillDepotAbilityGroup string `json:"skillDepotAbilityGroup"`
	// LeaderTalent           int32    `json:"leaderTalent,omitempty"`
}
type AvatarSkillExcel struct {
	ID           int32  `json:"id"`
	CostElemType string `json:"costElemType,omitempty"`
	// NameTextMapHash    int32     `json:"nameTextMapHash"`
	// AbilityName        string    `json:"abilityName"`
	// DescTextMapHash    int32       `json:"descTextMapHash"`
	// SkillIcon          string    `json:"skillIcon"`
	// CdTime             float64   `json:"cdTime,omitempty"`
	CostElemVal float64 `json:"costElemVal,omitempty"`
	// MaxChargeNum       int32       `json:"maxChargeNum"`
	// TriggerID          int32       `json:"triggerID,omitempty"`
	// LockShape          string    `json:"lockShape"`
	// LockWeightParams   []float64 `json:"lockWeightParams"`
	// IsAttackCameraLock bool      `json:"isAttackCameraLock"`
	// BuffIcon           string    `json:"buffIcon"`
	// GlobalValueKey     string    `json:"globalValueKey"`
	// CostStamina        float64   `json:"costStamina,omitempty"`
}

type AvatarFetterInfo struct {
	AvatarAssocType string `json:"avatarAssocType"`
	AvatarId        int32  `json:"avatarId"`
	FetterID        int32  `json:"fetterId"`
}

type AvatarPromote struct {
	AvatarPromoteID int32     `json:"avatarPromoteId"`
	AddProps        []AddProp `json:"addProps"`
	// PromoteLevel    int32     `json:"promoteLevel,omitempty"`
	// PromoteAudio    string `json:"promoteAudio"`
	// CostItems       []struct {
	// } `json:"costItems"`
	UnlockMaxLevel int32 `json:"unlockMaxLevel"`
	// ScoinCost           int32 `json:"scoinCost,omitempty"`
	// RequiredPlayerLevel int32 `json:"requiredPlayerLevel,omitempty"`
}
type AddProp struct {
	PropType string  `json:"propType"`
	Value    float64 `json:"value"`
}

type WeaponExcel struct {
	WeaponType string `json:"weaponType"`
	RankLevel  int32  `json:"rankLevel"`
	// WeaponBaseExp int32    `json:"weaponBaseExp"`
	SkillAffix []int32 `json:"skillAffix"`
	WeaponProp []struct {
		PropType  string  `json:"propType,omitempty"`
		InitValue float64 `json:"initValue,omitempty"`
		Type      string  `json:"type"`
	} `json:"weaponProp"`
	// AwakenTexture              string        `json:"awakenTexture"`
	// AwakenLightMapTexture      string        `json:"awakenLightMapTexture"`
	// AwakenIcon                 string        `json:"awakenIcon"`
	WeaponPromoteID int32 `json:"weaponPromoteId"`
	// StoryID                    int32           `json:"storyId"`
	// AwakenCosts                []interface{} `json:"awakenCosts"`
	// GachaCardNameHashSuffix    int32         `json:"gachaCardNameHashSuffix"`
	// DestroyRule                string        `json:"destroyRule"`
	// DestroyReturnMaterial      []int32         `json:"destroyReturnMaterial"`
	// DestroyReturnMaterialCount []int32         `json:"destroyReturnMaterialCount"`
	ID              int32 `json:"id"`
	NameTextMapHash int32 `json:"nameTextMapHash"`
	// DescTextMapHash            int32         `json:"descTextMapHash"`
	Icon     string `json:"icon"`
	ItemType string `json:"itemType"`
	// Weight                     int32           `json:"weight"`
	Rank int32 `json:"rank"`
	// GadgetID                   int32           `json:"gadgetId"`
}

type WeaponCurve struct {
}

type WeaponPromoteConfig struct {
	WeaponPromoteID int32 `json:"weaponPromoteId"`
	// CostItems       []struct {
	// } `json:"costItems"`
	AddProps       []AddProp
	UnlockMaxLevel int32 `json:"unlockMaxLevel"`
	PromoteLevel   int32 `json:"promoteLevel,omitempty"`
	// RequiredPlayerLevel int32 `json:"requiredPlayerLevel,omitempty"`
	// CoinCost            int32 `json:"coinCost,omitempty"`
}
