package dm

//nolint:tagliatelle // need to match datamine
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
	NameTextMapHash int64           `json:"nameTextMapHash"`
	UseType         string          `json:"useType,omitempty"`
}

//nolint:tagliatelle // need to match datamine
type PropGrowCurve struct {
	Type      string `json:"type"`
	GrowCurve string `json:"growCurve"`
}

type TextMap map[int64]string

//nolint:tagliatelle // need to match datamine
type AvatarSkillDepot struct {
	ID          int32   `json:"id"`
	EnergySkill int32   `json:"energySkill"`
	Skills      []int32 `json:"skills"`
	// SubSkills               []int32    `json:"subSkills"`
	// ExtraAbilities          []string `json:"extraAbilities"`
	// Talents                 []int32    `json:"talents"`
	// TalentStarName          string   `json:"talentStarName"`
	InherentProudSkillOpens []struct {
		ProudSkillGroupId      int32 `json:"proudSkillGroupId"`
		NeedAvatarPromoteLevel int32 `json:"needAvatarPromoteLevel"`
	} `json:"inherentProudSkillOpens"`
	// SkillDepotAbilityGroup string `json:"skillDepotAbilityGroup"`
	// LeaderTalent           int32    `json:"leaderTalent,omitempty"`
}

//nolint:tagliatelle // need to match datamine
type AvatarSkillExcel struct {
	ID              int32  `json:"id"`
	CostElemType    string `json:"costElemType,omitempty"`
	NameTextMapHash int64  `json:"nameTextMapHash"`
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
	ProudSkillGroupID int32 `json:"proudSkillGroupId"`
}

//nolint:tagliatelle // need to match datamine
type AvatarFetterInfo struct {
	AvatarAssocType string `json:"avatarAssocType"`
	AvatarId        int32  `json:"avatarId"`
	FetterID        int32  `json:"fetterId"`
}

//nolint:tagliatelle // need to match datamine
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

//nolint:tagliatelle // need to match datamine
type AddProp struct {
	PropType string  `json:"propType"`
	Value    float64 `json:"value"`
}

//nolint:tagliatelle // need to match datamine
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
	NameTextMapHash int64 `json:"nameTextMapHash"`
	// DescTextMapHash            int32         `json:"descTextMapHash"`
	Icon     string `json:"icon"`
	ItemType string `json:"itemType"`
	// Weight                     int32           `json:"weight"`
	// Rank int32 `json:"rank"`
	// GadgetID                   int32           `json:"gadgetId"`
}

type WeaponCurve struct {
}

//nolint:tagliatelle // need to match datamine
type WeaponPromote struct {
	WeaponPromoteID int32 `json:"weaponPromoteId"`
	// CostItems       []struct {
	// } `json:"costItems"`
	AddProps       []AddProp
	UnlockMaxLevel int32 `json:"unlockMaxLevel"`
	PromoteLevel   int32 `json:"promoteLevel,omitempty"`
	// RequiredPlayerLevel int32 `json:"requiredPlayerLevel,omitempty"`
	// CoinCost            int32 `json:"coinCost,omitempty"`
}

//nolint:tagliatelle // need to match datamine
type ReliquarySetExcel struct {
	SetID        int64 `json:"setId"`
	EquipAffixID int64 `json:"EquipAffixId"`
}

//nolint:tagliatelle // need to match datamine
type ReliquaryExcel struct {
	SetID         int64  `json:"setId"`
	EquipType     string `json:"equipType"`
	Icon          string `json:"icon"`
	AppendPropNum int32  `json:"appendPropNum"`
}

//nolint:tagliatelle // need to match datamine
type EquipAffixExcel struct {
	AffixID         int64 `json:"affixId"`
	ID              int64 `json:"id"`
	NameTextMapHash int64 `json:"nameTextMapHash"`
	Level           int32 `json:"level"`
	// DescTextMapHash            int32         `json:"descTextMapHash"`
	// AddProps []AddProp `json:"addProps"`
}

//nolint:tagliatelle // need to match datamine
type ProudSkillExcel struct {
	ProudSkillID          int32  `json:"proudSkillId"`
	ProudSkillGroupID     int32  `json:"proudSkillGroupId"`
	Level                 int32  `json:"level"`
	ProudSkillType        int32  `json:"proudSkillType"`
	NameTextMapHash       int64  `json:"nameTextMapHash"`
	DescTextMapHash       int64  `json:"descTextMapHash"`
	UnlockDescTextMapHash int64  `json:"unlockDescTextMapHash"`
	Icon                  string `json:"icon"`
	// CoinCost              int    `json:"coinCost"`
	// CostItems             []struct {
	// 	ID    int `json:"id,omitempty"`
	// 	Count int `json:"count,omitempty"`
	// } `json:"costItems"`
	// FilterConds      []string      `json:"filterConds"`
	// BreakLevel       int           `json:"breakLevel"`
	// ParamDescList    []interface{} `json:"paramDescList"`
	// LifeEffectParams []string      `json:"lifeEffectParams"`
	// OpenConfig       string        `json:"openConfig"`
	// AddProps []struct {
	// } `json:"addProps"`
	ParamList []float64 `json:"paramList"`
}

//nolint:tagliatelle // need to match datamine
type HpDrop struct {
	DropId    int32   `json:"dropId"`
	HpPercent float64 `json:"hpPercent"`
}

//nolint:tagliatelle // need to match datamine
type MonsterExcel struct {
	MonsterName     string          `json:"monsterName"`
	Typ             string          `json:"type"`
	HpDrops         []HpDrop        `json:"hpDrops"`
	DescribeId      int32           `json:"describeId"`
	KillDropId      int32           `json:"killDropId"`
	HpBase          float64         `json:"hpBase"`
	PropGrowCurves  []PropGrowCurve `json:"propGrowCurves"`
	FireSubHurt     float64         `json:"fireSubHurt"`
	GrassSubHurt    float64         `json:"grassSubHurt"`
	WaterSubHurt    float64         `json:"waterSubHurt"`
	ElecSubHurt     float64         `json:"elecSubHurt"`
	WindSubHurt     float64         `json:"windSubHurt"`
	IceSubHurt      float64         `json:"iceSubHurt"`
	RockSubHurt     float64         `json:"rockSubHurt"`
	PhysicalSubHurt float64         `json:"physicalSubHurt"`
	Id              int32           `json:"id"`
	// ...
}

//nolint:tagliatelle // need to match datamine
type MonsterDescribeExcel struct {
	Id              int32  `json:"id"`
	NameTextMapHash int64  `json:"nameTextMapHash"`
	Icon            string `json:"icon"`
	// ...
}

//nolint:tagliatelle // need to match datamine
type MonsterCurveExcel struct {
	Level      int32 `json:"level"`
	CurveInfos []struct {
		Type  string  `json:"type"`
		Value float64 `json:"value"`
		// ...
	} `json:"curveInfos"`
}
