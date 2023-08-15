package keys

import (
	"encoding/json"
	"errors"
	"strings"
)

type Set int

func (s *Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(setNames[*s])
}

func (c *Set) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	for i, v := range setNames {
		if v == s {
			*c = Set(i)
			return nil
		}
	}
	return errors.New("unrecognized set key")
}

func (c Set) String() string {
	return setNames[c]
}

var setNames = []string{
	"",
	"adventurer",
	"archaicpetra",
	"berserker",
	"blizzardstrayer",
	"bloodstainedchivalry",
	"braveheart",
	"crimsonwitchofflames",
	"deepwoodmemories",
	"defenderswill",
	"desertpavilionchronicle",
	"vourukashasglow",
	"echoesofanoffering",
	"emblemofseveredfate",
	"flowerofparadiselost",
	"gambler",
	"gladiatorsfinale",
	"gildeddreams",
	"goldentroupe",
	"heartofdepth",
	"huskofopulentdreams",
	"instructor",
	"lavawalker",
	"luckydog",
	"maidenbeloved",
	"marechausseehunter",
	"martialartist",
	"noblesseoblige",
	"nymphsdream",
	"oceanhuedclam",
	"paleflame",
	"prayersfordestiny",
	"prayersforillumination",
	"prayersforwisdom",
	"prayerstospringtime",
	"resolutionofsojourner",
	"retracingbolide",
	"scholar",
	"shimenawasreminiscence",
	"tenacityofthemillelith",
	"theexile",
	"thunderingfury",
	"thundersoother",
	"tinymiracle",
	"travelingdoctor",
	"vermillionhereafter",
	"viridescentvenerer",
	"wandererstroupe",
}

const (
	NoSet Set = iota
	Adventurer
	ArchaicPetra
	Berserker
	BlizzardStrayer
	BloodstainedChivalry
	BraveHeart
	CrimsonWitchOfFlames
	DeepwoodMemories
	DefendersWill
	DesertPavilionChronicle
	VourukashasGlow
	EchoesOfAnOffering
	EmblemOfSeveredFate
	FlowerOfParadiseLost
	Gambler
	GladiatorsFinale
	GildedDreams
	GoldenTroupe
	HeartOfDepth
	HuskOfOpulentDreams
	Instructor
	Lavawalker
	LuckyDog
	MaidenBeloved
	MarechausseeHunter
	MartialArtist
	NoblesseOblige
	NymphsDream
	OceanHuedClam
	PaleFlame
	PrayersForDestiny
	PrayersForIllumination
	PrayersForWisdom
	PrayersToSpringtime
	ResolutionOfSojourner
	RetracingBolide
	Scholar
	ShimenawasReminiscence
	TenacityOfTheMillelith
	TheExile
	ThunderingFury
	Thundersoother
	TinyMiracle
	TravelingDoctor
	VermillionHereafter
	ViridescentVenerer
	WanderersTroupe
)
