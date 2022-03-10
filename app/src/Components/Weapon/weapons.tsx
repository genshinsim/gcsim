import { MenuItem } from "@blueprintjs/core";
import { ItemPredicate, ItemRenderer } from "@blueprintjs/select";

export const weaponKeyToName = {
  English: {
    akuoumaru: "Akuoumaru",
    alleyhunter: "Alley Hunter",
    amenomakageuchi: "Amenoma Kageuchi",
    amosbow: "Amos' Bow",
    apprenticesnotes: "Apprentice's Notes",
    aquilafavonia: "Aquila Favonia",
    beginnersprotector: "Beginner's Protector",
    blackcliffagate: "Blackcliff Agate",
    blackclifflongsword: "Blackcliff Longsword",
    blackcliffpole: "Blackcliff Pole",
    blackcliffslasher: "Blackcliff Slasher",
    blackcliffwarbow: "Blackcliff Warbow",
    blacktassel: "Black Tassel",
    bloodtaintedgreatsword: "Bloodtainted Greatsword",
    calamityqueller: "Calamity Queller",
    cinnabarspindle: "Cinnabar Spindle",
    compoundbow: "Compound Bow",
    coolsteel: "Cool Steel",
    crescentpike: "Crescent Pike",
    darkironsword: "Dark Iron Sword",
    deathmatch: "Deathmatch",
    debateclub: "Debate Club",
    dodocotales: "Dodoco Tales",
    dragonsbane: "Dragon's Bane",
    dragonspinespear: "Dragonspine Spear",
    dullblade: "Dull Blade",
    elegyfortheend: "Elegy for the End",
    emeraldorb: "Emerald Orb",
    engulfinglightning: "Engulfing Lightning",
    everlastingmoonglow: "Everlasting Moonglow",
    eyeofperception: "Eye of Perception",
    favoniuscodex: "Favonius Codex",
    favoniusgreatsword: "Favonius Greatsword",
    favoniuslance: "Favonius Lance",
    favoniussword: "Favonius Sword",
    favoniuswarbow: "Favonius Warbow",
    ferrousshadow: "Ferrous Shadow",
    festeringdesire: "Festering Desire",
    filletblade: "Fillet Blade",
    freedomsworn: "Freedom-Sworn",
    frostbearer: "Frostbearer",
    hakushinring: "Hakushin Ring",
    halberd: "Halberd",
    hamayumi: "Hamayumi",
    harbingerofdawn: "Harbinger of Dawn",
    huntersbow: "Hunter's Bow",
    ironpoint: "Iron Point",
    ironsting: "Iron Sting",
    kagurasverity: "Kagura's Verity",
    katsuragikirinagamasa: "Katsuragikiri Nagamasa",
    kitaincrossspear: "Kitain Cross Spear",
    lionsroar: "Lion's Roar",
    lithicblade: "Lithic Blade",
    lithicspear: "Lithic Spear",
    lostprayertothesacredwinds: "Lost Prayer to the Sacred Winds",
    luxurioussealord: "Luxurious Sea-Lord",
    magicguide: "Magic Guide",
    mappamare: "Mappa Mare",
    memoryofdust: "Memory of Dust",
    messenger: "Messenger",
    mistsplitterreforged: "Mistsplitter Reforged",
    mitternachtswaltz: "Mitternachts Waltz",
    mouunsmoon: "Mouun's Moon",
    oathsworneye: "Oathsworn Eye",
    oldmercspal: "Old Merc's Pal",
    otherworldlystory: "Otherworldly Story",
    pocketgrimoire: "Pocket Grimoire",
    polarstar: "Polar Star",
    predator: "Predator",
    primordialjadecutter: "Primordial Jade Cutter",
    primordialjadewingedspear: "Primordial Jade Winged-Spear",
    prototypeamber: "Prototype Amber",
    prototypearchaic: "Prototype Archaic",
    prototypecrescent: "Prototype Crescent",
    prototyperancour: "Prototype Rancour",
    prototypestarglitter: "Prototype Starglitter",
    rainslasher: "Rainslasher",
    ravenbow: "Raven Bow",
    recurvebow: "Recurve Bow",
    redhornstonethresher: "Redhorn Stonethresher",
    royalbow: "Royal Bow",
    royalgreatsword: "Royal Greatsword",
    royalgrimoire: "Royal Grimoire",
    royallongsword: "Royal Longsword",
    royalspear: "Royal Spear",
    rust: "Rust",
    sacrificialbow: "Sacrificial Bow",
    sacrificialfragments: "Sacrificial Fragments",
    sacrificialgreatsword: "Sacrificial Greatsword",
    sacrificialsword: "Sacrificial Sword",
    seasonedhuntersbow: "Seasoned Hunter's Bow",
    serpentspine: "Serpent Spine",
    sharpshootersoath: "Sharpshooter's Oath",
    silversword: "Silver Sword",
    skyridergreatsword: "Skyrider Greatsword",
    skyridersword: "Skyrider Sword",
    skywardatlas: "Skyward Atlas",
    skywardblade: "Skyward Blade",
    skywardharp: "Skyward Harp",
    skywardpride: "Skyward Pride",
    skywardspine: "Skyward Spine",
    slingshot: "Slingshot",
    snowtombedstarsilver: "Snow-Tombed Starsilver",
    solarpearl: "Solar Pearl",
    songofbrokenpines: "Song of Broken Pines",
    staffofhoma: "Staff of Homa",
    summitshaper: "Summit Shaper",
    swordofdescension: "Sword of Descension",
    thealleyflash: "The Alley Flash",
    thebell: "The Bell",
    theblacksword: "The Black Sword",
    thecatch: '"The Catch"',
    theflute: "The Flute",
    thestringless: "The Stringless",
    theunforged: "The Unforged",
    theviridescenthunt: "The Viridescent Hunt",
    thewidsith: "The Widsith",
    thrillingtalesofdragonslayers: "Thrilling Tales of Dragon Slayers",
    thunderingpulse: "Thundering Pulse",
    travelershandysword: "Traveler's Handy Sword",
    twinnephrite: "Twin Nephrite",
    vortexvanquisher: "Vortex Vanquisher",
    wastergreatsword: "Waster Greatsword",
    wavebreakersfin: "Wavebreaker's Fin",
    whiteblind: "Whiteblind",
    whiteirongreatsword: "White Iron Greatsword",
    whitetassel: "White Tassel",
    windblumeode: "Windblume Ode",
    wineandsong: "Wine and Song",
    wolfsgravestone: "Wolf's Gravestone",
  },
  Chinese: {
    akuoumaru: "恶王丸",
    alleyhunter: "暗巷猎手",
    amenomakageuchi: "天目影打刀",
    amosbow: "阿莫斯之弓",
    apprenticesnotes: "学徒笔记",
    aquilafavonia: "风鹰剑",
    beginnersprotector: "新手长枪",
    blackcliffagate: "黑岩绯玉",
    blackclifflongsword: "黑岩长剑",
    blackcliffpole: "黑岩刺枪",
    blackcliffslasher: "黑岩斩刀",
    blackcliffwarbow: "黑岩战弓",
    blacktassel: "黑缨枪",
    bloodtaintedgreatsword: "沐浴龙血的剑",
    calamityqueller: "息灾",
    cinnabarspindle: "辰砂之纺锤",
    compoundbow: "钢轮弓",
    coolsteel: "冷刃",
    crescentpike: "流月针",
    darkironsword: "暗铁剑",
    deathmatch: "决斗之枪",
    debateclub: "以理服人",
    dodocotales: "嘟嘟可故事集",
    dragonsbane: "匣里灭辰",
    dragonspinespear: "龙脊长枪",
    dullblade: "无锋剑",
    elegyfortheend: "终末嗟叹之诗",
    emeraldorb: "翡玉法球",
    engulfinglightning: "薙草之稻光",
    everlastingmoonglow: "不灭月华",
    eyeofperception: "昭心",
    favoniuscodex: "西风秘典",
    favoniusgreatsword: "西风大剑",
    favoniuslance: "西风长枪",
    favoniussword: "西风剑",
    favoniuswarbow: "西风猎弓",
    ferrousshadow: "铁影阔剑",
    festeringdesire: "腐殖之剑",
    filletblade: "吃虎鱼刀",
    freedomsworn: "苍古自由之誓",
    frostbearer: "忍冬之果",
    hakushinring: "白辰之环",
    halberd: "钺矛",
    hamayumi: "破魔之弓",
    harbingerofdawn: "黎明神剑",
    huntersbow: "猎弓",
    ironpoint: "铁尖枪",
    ironsting: "铁蜂刺",
    kagurasverity: "神乐之真意",
    katsuragikirinagamasa: "桂木斩长正",
    kitaincrossspear: "喜多院十文字",
    lionsroar: "匣里龙吟",
    lithicblade: "千岩古剑",
    lithicspear: "千岩长枪",
    lostprayertothesacredwinds: "四风原典",
    luxurioussealord: "衔珠海皇",
    magicguide: "魔导绪论",
    mappamare: "万国诸海图谱",
    memoryofdust: "尘世之锁",
    messenger: "信使",
    mistsplitterreforged: "雾切之回光",
    mitternachtswaltz: "幽夜华尔兹",
    mouunsmoon: "曚云之月",
    oathsworneye: "证誓之明瞳",
    oldmercspal: "佣兵重剑",
    otherworldlystory: "异世界行记",
    pocketgrimoire: "口袋魔导书",
    polarstar: "冬极白星",
    predator: "掠食者",
    primordialjadecutter: "磐岩结绿",
    primordialjadewingedspear: "和璞鸢",
    prototypeamber: "试作金珀",
    prototypearchaic: "试作古华",
    prototypecrescent: "试作澹月",
    prototyperancour: "试作斩岩",
    prototypestarglitter: "试作星镰",
    rainslasher: "雨裁",
    ravenbow: "鸦羽弓",
    recurvebow: "反曲弓",
    redhornstonethresher: "赤角石溃杵",
    royalbow: "宗室长弓",
    royalgreatsword: "宗室大剑",
    royalgrimoire: "宗室秘法录",
    royallongsword: "宗室长剑",
    royalspear: "宗室猎枪",
    rust: "弓藏",
    sacrificialbow: "祭礼弓",
    sacrificialfragments: "祭礼残章",
    sacrificialgreatsword: "祭礼大剑",
    sacrificialsword: "祭礼剑",
    seasonedhuntersbow: "Seasoned 猎弓",
    serpentspine: "螭骨剑",
    sharpshootersoath: "神射手之誓",
    silversword: "银剑",
    skyridergreatsword: "飞天大御剑",
    skyridersword: "飞天御剑",
    skywardatlas: "天空之卷",
    skywardblade: "天空之刃",
    skywardharp: "天空之翼",
    skywardpride: "天空之傲",
    skywardspine: "天空之脊",
    slingshot: "弹弓",
    snowtombedstarsilver: "雪葬的星银",
    solarpearl: "匣里日月",
    songofbrokenpines: "松籁响起之时",
    staffofhoma: "护摩之杖",
    summitshaper: "斫峰之刃",
    swordofdescension: "降临之剑",
    thealleyflash: "暗巷闪光",
    thebell: "钟剑",
    theblacksword: "黑剑",
    thecatch: '「渔获」',
    theflute: "笛剑",
    thestringless: "绝弦",
    theunforged: "无工之剑",
    theviridescenthunt: "苍翠猎弓",
    thewidsith: "流浪乐章",
    thrillingtalesofdragonslayers: "讨龙英杰谭",
    thunderingpulse: "飞雷之弦振",
    travelershandysword: "旅行剑",
    twinnephrite: "甲级宝珏",
    vortexvanquisher: "贯虹之槊",
    wastergreatsword: "训练大剑",
    wavebreakersfin: "断浪长鳍",
    whiteblind: "白影剑",
    whiteirongreatsword: "白铁大剑",
    whitetassel: "白缨枪",
    windblumeode: "风花之颂",
    wineandsong: "暗巷的酒与诗",
    wolfsgravestone: "狼的末路",
  },
};

export const weapons = {
  English: Object.keys(weaponKeyToName.English).map((k) => ({
    key: k,
    name: weaponKeyToName.English[k],
  })),
  Chinese: Object.keys(weaponKeyToName.Chinese).map((k) => ({
    key: k,
    name: weaponKeyToName.Chinese[k],
  })),
};

export interface IWeapon {
  key: string;
  name: string;
}

export const renderWeapon: ItemRenderer<IWeapon> = (
  weapon,
  { handleClick, modifiers, query }
) => {
  if (!modifiers.matchesPredicate) {
    return null;
  }
  return (
    <MenuItem
      active={modifiers.active}
      disabled={modifiers.disabled}
      label={weapon.name}
      key={weapon.key}
      onClick={handleClick}
      text={highlightText(weapon.name, query)}
    />
  );
};

export const filterWeapon: ItemPredicate<IWeapon> = (
  query,
  weapon,
  _index,
  exactMatch
) => {
  const normalizedQuery = query.toLowerCase();

  if (exactMatch) {
    return weapon.key === normalizedQuery;
  } else {
    return `${weapon.key} ${weapon.name}`.indexOf(normalizedQuery) >= 0;
  }
};

function escapeRegExpChars(text: string) {
  return text.replace(/([.*+?^=!:${}()|\[\]\/\\])/g, "\\$1");
}

function highlightText(text: string, query: string) {
  let lastIndex = 0;
  const words = query
    .split(/\s+/)
    .filter((word) => word.length > 0)
    .map(escapeRegExpChars);
  if (words.length === 0) {
    return [text];
  }
  const regexp = new RegExp(words.join("|"), "gi");
  const tokens: React.ReactNode[] = [];
  while (true) {
    const match = regexp.exec(text);
    if (!match) {
      break;
    }
    const length = match[0].length;
    const before = text.slice(lastIndex, regexp.lastIndex - length);
    if (before.length > 0) {
      tokens.push(before);
    }
    lastIndex = regexp.lastIndex;
    tokens.push(<strong key={lastIndex}>{match[0]}</strong>);
  }
  const rest = text.slice(lastIndex);
  if (rest.length > 0) {
    tokens.push(rest);
  }
  return tokens;
}

export const weaponSelectProps = {
  English: {
    itemPredicate: filterWeapon,
    itemRenderer: renderWeapon,
    items: weapons.English,
  },
  Chinese: {
    itemPredicate: filterWeapon,
    itemRenderer: renderWeapon,
    items: weapons.Chinese,
  },
};
