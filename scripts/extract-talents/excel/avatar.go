package excel

var (
	AvatarExcelConfigData           []*Avatar
	AvatarSkillDepotExcelConfigData []*AvatarSkillDepot
	AvatarSkillExcelConfigData      []*AvatarSkill
	ProudSkillExcelConfigData       []*ProudSkill
)

func init() {
	load("ExcelBinOutput/AvatarExcelConfigData.json", &AvatarExcelConfigData)
	load("ExcelBinOutput/AvatarSkillDepotExcelConfigData.json", &AvatarSkillDepotExcelConfigData)
	load("ExcelBinOutput/AvatarSkillExcelConfigData.json", &AvatarSkillExcelConfigData)
	load("ExcelBinOutput/ProudSkillExcelConfigData.json", &ProudSkillExcelConfigData)
}

type Avatar struct {
	SkillDepotId uint32
	Id           uint32
}

func (a *Avatar) SkillDepot() *AvatarSkillDepot {
	return FindSkillDepot(a.SkillDepotId)
}

type AvatarSkillDepot struct {
	Id          uint32
	EnergySkill uint32
	Skills      []uint32
}

type AvatarSkill struct {
	Id                uint32
	NameTextMapHash   TextMapHash
	ProudSkillGroupId uint32
}

func (a *AvatarSkill) Name() string {
	return a.NameTextMapHash.String()
}

func (a *AvatarSkill) ProudSkill(level uint32) *ProudSkill {
	for _, v := range ProudSkillExcelConfigData {
		if v.ProudSkillGroupId == a.ProudSkillGroupId && v.Level == level {
			return v
		}
	}
	return nil
}

type ProudSkill struct {
	ProudSkillGroupId uint32
	Level             uint32
	ParamDescList     []TextMapHash
}

func FindAvatar(id uint32) *Avatar {
	for _, v := range AvatarExcelConfigData {
		if v.Id == id {
			return v
		}
	}
	return nil
}

func FindSkillDepot(id uint32) *AvatarSkillDepot {
	for _, v := range AvatarSkillDepotExcelConfigData {
		if v.Id == id {
			return v
		}
	}
	return nil
}

func FindSkill(id uint32) *AvatarSkill {
	for _, v := range AvatarSkillExcelConfigData {
		if v.Id == id {
			return v
		}
	}
	return nil
}
