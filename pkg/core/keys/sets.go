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

func (s *Set) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	str = strings.ToLower(str)
	for i, v := range setNames {
		if v == str {
			*s = Set(i)
			return nil
		}
	}
	return errors.New("unrecognized set key")
}

func (s Set) String() string {
	return setNames[s]
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
	"fragmentofharmonicwhimsy",
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
	"nighttimewhispersintheechoingwoods",
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
	"songofdayspast",
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
	FragmentOfHarmonicWhimsy
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
	NighttimeWhispersInTheEchoingWoods
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
	SongOfDaysPast
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
