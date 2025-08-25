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
	"echoesofanoffering",
	"emblemofseveredfate",
	"finaleofthedeepgalleries",
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
	"longnightsoath",
	"luckydog",
	"maidenbeloved",
	"marechausseehunter",
	"martialartist",
	"nighttimewhispersintheechoingwoods",
	"noblesseoblige",
	"nymphsdream",
	"obsidiancodex",
	"oceanhuedclam",
	"paleflame",
	"prayersfordestiny",
	"prayersforillumination",
	"prayersforwisdom",
	"prayerstospringtime",
	"resolutionofsojourner",
	"retracingbolide",
	"scholar",
	"scrolloftheheroofcindercity",
	"shimenawasreminiscence",
	"songofdayspast",
	"tenacityofthemillelith",
	"theexile",
	"thunderingfury",
	"thundersoother",
	"tinymiracle",
	"travelingdoctor",
	"unfinishedreverie",
	"vermillionhereafter",
	"viridescentvenerer",
	"vourukashasglow",
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
	EchoesOfAnOffering
	EmblemOfSeveredFate
	FinaleOfTheDeepGalleries
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
	LongNightsOath
	LuckyDog
	MaidenBeloved
	MarechausseeHunter
	MartialArtist
	NighttimeWhispersInTheEchoingWoods
	NoblesseOblige
	NymphsDream
	ObsidianCodex
	OceanHuedClam
	PaleFlame
	PrayersForDestiny
	PrayersForIllumination
	PrayersForWisdom
	PrayersToSpringtime
	ResolutionOfSojourner
	RetracingBolide
	Scholar
	ScrollOfTheHeroOfCinderCity
	ShimenawasReminiscence
	SongOfDaysPast
	TenacityOfTheMillelith
	TheExile
	ThunderingFury
	Thundersoother
	TinyMiracle
	TravelingDoctor
	UnfinishedReverie
	VermillionHereafter
	ViridescentVenerer
	VourukashasGlow
	WanderersTroupe
)
