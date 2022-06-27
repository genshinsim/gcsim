export interface IGOOD {
  format: "GOOD"; // A way for people to recognize this format.
  version: number; // GOOD API version.
  source: string; // The app that generates this data.
  characters?: GOODCharacter[];
  artifacts?: GOODArtifact[];
  weapons?: GOODWeapon[];
  materials?: Object;
}
export interface GOODArtifact {
  setKey: GOODArtifactSetKey; //e.g. "GladiatorsFinale"
  slotKey: GOODSlotKey; //e.g. "plume"
  level: number; //0-20 inclusive
  rarity: number; //1-5 inclusive
  mainStatKey: GOODStatKey;
  location: GOODCharacterKey | ""; //where "" means not equipped.
  lock: boolean; //Whether the artifact is locked in game.
  substats: ISubstat[];
}

export interface ISubstat {
  key: GOODStatKey;
  value: number;
}

export type GOODSlotKey = "flower" | "plume" | "sands" | "goblet" | "circlet";

export interface GOODWeapon {
  key: GOODWeaponKey; //"CrescentPike"
  level: number; //1-90 inclusive
  ascension: number; //0-6 inclusive. need to disambiguate 80/90 or 80/80
  refinement: number; //1-5 inclusive
  location: GOODCharacterKey | ""; //where "" means not equipped.
  lock: boolean; //Whether the weapon is locked in game.
}

export interface GOODCharacter {
  key: GOODCharacterKey; //e.g. "Rosaria"
  level: number; //1-90 inclusive
  constellation: number; //0-6 inclusive
  ascension: number; //0-6 inclusive. need to disambiguate 80/90 or 80/80
  talent: {
    //does not include boost from constellations. 1-15 inclusive
    auto: number;
    skill: number;
    burst: number;
  };
}
export interface Weapon {
  key: string;
  name: string;
  icon: string;
  level: number;
  ascension: number;
  refinement: number;
}

export type GOODStatKey =
  | "hp" //HP
  | "hp_" //HP%
  | "atk" //ATK
  | "atk_" //ATK%
  | "def" //DEF
  | "def_" //DEF%
  | "eleMas" //Elemental Mastery
  | "enerRech_" //Energy Recharge%
  | "heal_" //Healing Bonus%
  | "critRate_" //CRIT Rate%
  | "critDMG_" //CRIT DMG%
  | "physical_dmg_" //Physical DMG Bonus%
  | "anemo_dmg_" //Anemo DMG Bonus%
  | "geo_dmg_" //Geo DMG Bonus%
  | "electro_dmg_" //Electro DMG Bonus%
  | "hydro_dmg_" //Hydro DMG Bonus%
  | "pyro_dmg_" //Pyro DMG Bonus%
  | "cryo_dmg_" //Cryo DMG Bonus%
  | ""; //Some scanners use this

export type GOODArtifactSetKey =
  | "Adventurer" //Adventurer
  | "ArchaicPetra" //Archaic Petra
  | "Berserker" //Berserker
  | "BlizzardStrayer" //Blizzard Strayer
  | "BloodstainedChivalry" //Bloodstained Chivalry
  | "BraveHeart" //Brave Heart
  | "CrimsonWitchOfFlames" //Crimson Witch of Flames
  | "DefendersWill" //Defender's Will
  | "EmblemOfSeveredFate" //Emblem of Severed Fate
  | "Gambler" //Gambler
  | "GladiatorsFinale" //Gladiator's Finale
  | "HeartOfDepth" //Heart of Depth
  | "HuskOfOpulentDreams" //Husk of Opulent Dreams
  | "Instructor" //Instructor
  | "Lavawalker" //Lavawalker
  | "LuckyDog" //Lucky Dog
  | "MaidenBeloved" //Maiden Beloved
  | "MartialArtist" //Martial Artist
  | "NoblesseOblige" //Noblesse Oblige
  | "OceanHuedClam" //Ocean-Hued Clam
  | "PaleFlame" //Pale Flame
  | "PrayersForDestiny" //Prayers for Destiny
  | "PrayersForIllumination" //Prayers for Illumination
  | "PrayersForWisdom" //Prayers for Wisdom
  | "PrayersToSpringtime" //Prayers to Springtime
  | "ResolutionOfSojourner" //Resolution of Sojourner
  | "RetracingBolide" //Retracing Bolide
  | "Scholar" //Scholar
  | "ShimenawasReminiscence" //Shimenawa's Reminiscence
  | "TenacityOfTheMillelith" //Tenacity of the Millelith
  | "TheExile" //The Exile
  | "ThunderingFury" //Thundering Fury
  | "Thundersoother" //Thundersoother
  | "TinyMiracle" //Tiny Miracle
  | "TravelingDoctor" //Traveling Doctor
  | "ViridescentVenerer" //Viridescent Venerer
  | "WanderersTroupe"; //Wanderer's Troupe

export type GOODCharacterKey =
  | "Albedo" //Albedo
  | "Aloy" //Aloy
  | "Amber" //Amber
  | "AratakiItto" //Arataki Itto
  | "Barbara" //Barbara
  | "Beidou" //Beidou
  | "Bennett" //Bennett
  | "Chongyun" //Chongyun
  | "Diluc" //Diluc
  | "Diona" //Diona
  | "Eula" //Eula
  | "Fischl" //Fischl
  | "Ganyu" //Ganyu
  | "Gorou" //Gorou
  | "HuTao" //Hu Tao
  | "Jean" //Jean
  | "KaedeharaKazuha" //Kaedehara Kazuha
  | "Kaeya" //Kaeya
  | "KamisatoAyaka" //Kamisato Ayaka
  | "KamisatoAyato" //Kamisato Ayato
  | "Keqing" //Keqing
  | "Klee" //Klee
  | "KujouSara" //Kujou Sara
  | "Lisa" //Lisa
  | "Mona" //Mona
  | "Ningguang" //Ningguang
  | "Noelle" //Noelle
  | "Qiqi" //Qiqi
  | "RaidenShogun" //Raiden Shogun
  | "Razor" //Razor
  | "Rosaria" //Rosaria
  | "SangonomiyaKokomi" //Sangonomiya Kokomi
  | "Sayu" //Sayu
  | "Shenhe" //Shenhe
  | "Sucrose" //Sucrose
  | "Tartaglia" //Tartaglia
  | "Thoma" //Thoma
  | "Traveler" //Traveler
  | "Venti" //Venti
  | "Xiangling" //Xiangling
  | "Xiao" //Xiao
  | "Xingqiu" //Xingqiu
  | "Xinyan" //Xinyan
  | "YaeMiko" //Yae Miko
  | "Yanfei" //Yanfei
  | "Yoimiya" //Yoimiya
  | "Yelan" //Yelan
  | "Zhongli" //Zhongli
  | "KukiShinobu" //Kuki
  | "YunJin"; //Yun Jin

export type GOODWeaponKey =
  | "Akuoumaru" //Akuoumaru
  | "AlleyHunter" //Alley Hunter
  | "AmenomaKageuchi" //Amenoma Kageuchi
  | "AmosBow" //Amos' Bow
  | "ApprenticesNotes" //Apprentice's Notes
  | "AquilaFavonia" //Aquila Favonia
  | "BeginnersProtector" //Beginner's Protector
  | "BlackTassel" //Black Tassel
  | "BlackcliffAgate" //Blackcliff Agate
  | "BlackcliffLongsword" //Blackcliff Longsword
  | "BlackcliffPole" //Blackcliff Pole
  | "BlackcliffSlasher" //Blackcliff Slasher
  | "BlackcliffWarbow" //Blackcliff Warbow
  | "BloodtaintedGreatsword" //Bloodtainted Greatsword
  | "CinnabarSpindle" //Cinnabar Spindle
  | "CompoundBow" //Compound Bow
  | "CoolSteel" //Cool Steel
  | "CrescentPike" //Crescent Pike
  | "DarkIronSword" //Dark Iron Sword
  | "Deathmatch" //Deathmatch
  | "DebateClub" //Debate Club
  | "DodocoTales" //Dodoco Tales
  | "DragonsBane" //Dragon's Bane
  | "DragonspineSpear" //Dragonspine Spear
  | "DullBlade" //Dull Blade
  | "ElegyForTheEnd" //Elegy for the End
  | "EmeraldOrb" //Emerald Orb
  | "EngulfingLightning" //Engulfing Lightning
  | "EverlastingMoonglow" //Everlasting Moonglow
  | "EyeOfPerception" //Eye of Perception
  | "FavoniusCodex" //Favonius Codex
  | "FavoniusGreatsword" //Favonius Greatsword
  | "FavoniusLance" //Favonius Lance
  | "FavoniusSword" //Favonius Sword
  | "FavoniusWarbow" //Favonius Warbow
  | "FerrousShadow" //Ferrous Shadow
  | "FesteringDesire" //Festering Desire
  | "FilletBlade" //Fillet Blade
  | "FreedomSworn" //Freedom-Sworn
  | "Frostbearer" //Frostbearer
  | "HakushinRing" //Hakushin Ring
  | "Halberd" //Halberd
  | "Hamayumi" //Hamayumi
  | "HarbingerOfDawn" //Harbinger of Dawn
  | "HuntersBow" //Hunter's Bow
  | "IronPoint" //Iron Point
  | "IronSting" //Iron Sting
  | "KatsuragikiriNagamasa" //Katsuragikiri Nagamasa
  | "KitainCrossSpear" //Kitain Cross Spear
  | "LionsRoar" //Lion's Roar
  | "LithicBlade" //Lithic Blade
  | "LithicSpear" //Lithic Spear
  | "LostPrayerToTheSacredWinds" //Lost Prayer to the Sacred Winds
  | "LuxuriousSeaLord" //Luxurious Sea-Lord
  | "MagicGuide" //Magic Guide
  | "MappaMare" //Mappa Mare
  | "MemoryOfDust" //Memory of Dust
  | "Messenger" //Messenger
  | "MistsplitterReforged" //Mistsplitter Reforged
  | "MitternachtsWaltz" //Mitternachts Waltz
  | "MouunsMoon" //Mouun's Moon
  | "OldMercsPal" //Old Merc's Pal
  | "OtherworldlyStory" //Otherworldly Story
  | "PocketGrimoire" //Pocket Grimoire
  | "PolarStar" //Polar Star
  | "Predator" //Predator
  | "PrimordialJadeCutter" //Primordial Jade Cutter
  | "PrimordialJadeWingedSpear" //Primordial Jade Winged-Spear
  | "PrototypeAmber" //Prototype Amber
  | "PrototypeArchaic" //Prototype Archaic
  | "PrototypeCrescent" //Prototype Crescent
  | "PrototypeRancour" //Prototype Rancour
  | "PrototypeStarglitter" //Prototype Starglitter
  | "Rainslasher" //Rainslasher
  | "RavenBow" //Raven Bow
  | "RecurveBow" //Recurve Bow
  | "RedhornStonethresher" //Redhorn Stonethresher
  | "RoyalBow" //Royal Bow
  | "RoyalGreatsword" //Royal Greatsword
  | "RoyalGrimoire" //Royal Grimoire
  | "RoyalLongsword" //Royal Longsword
  | "RoyalSpear" //Royal Spear
  | "Rust" //Rust
  | "SacrificialBow" //Sacrificial Bow
  | "SacrificialFragments" //Sacrificial Fragments
  | "SacrificialGreatsword" //Sacrificial Greatsword
  | "SacrificialSword" //Sacrificial Sword
  | "SeasonedHuntersBow" //Seasoned Hunter's Bow
  | "SerpentSpine" //Serpent Spine
  | "SharpshootersOath" //Sharpshooter's Oath
  | "SilverSword" //Silver Sword
  | "SkyriderGreatsword" //Skyrider Greatsword
  | "SkyriderSword" //Skyrider Sword
  | "SkywardAtlas" //Skyward Atlas
  | "SkywardBlade" //Skyward Blade
  | "SkywardHarp" //Skyward Harp
  | "SkywardPride" //Skyward Pride
  | "SkywardSpine" //Skyward Spine
  | "Slingshot" //Slingshot
  | "SnowTombedStarsilver" //Snow-Tombed Starsilver
  | "SolarPearl" //Solar Pearl
  | "SongOfBrokenPines" //Song of Broken Pines
  | "StaffOfHoma" //Staff of Homa
  | "SummitShaper" //Summit Shaper
  | "SwordOfDescension" //Sword of Descension
  | "TheAlleyFlash" //The Alley Flash
  | "TheBell" //The Bell
  | "TheBlackSword" //The Black Sword
  | "TheCatch" //"The Catch"
  | "TheFlute" //The Flute
  | "TheStringless" //The Stringless
  | "TheUnforged" //The Unforged
  | "TheViridescentHunt" //The Viridescent Hunt
  | "TheWidsith" //The Widsith
  | "ThrillingTalesOfDragonSlayers" //Thrilling Tales of Dragon Slayers
  | "ThunderingPulse" //Thundering Pulse
  | "TravelersHandySword" //Traveler's Handy Sword
  | "TwinNephrite" //Twin Nephrite
  | "VortexVanquisher" //Vortex Vanquisher
  | "WasterGreatsword" //Waster Greatsword
  | "WavebreakersFin" //Wavebreaker's Fin
  | "WhiteIronGreatsword" //White Iron Greatsword
  | "WhiteTassel" //White Tassel
  | "Whiteblind" //Whiteblind
  | "WindblumeOde" //Windblume Ode
  | "WineAndSong" //Wine and Song
  | "WolfsGravestone"; //Wolf's Gravestone
