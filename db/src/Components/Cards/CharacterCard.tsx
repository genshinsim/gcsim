type CharacterCardProps = {
  char: string;
  custStyle?: string;
  longName?: string;
  onClick?: () => void;
};
export function CharacterCard({
  char,
  custStyle = "",
  longName = "",
  onClick,
}: CharacterCardProps) {
  var cls = "rounded-md bg-opacity-70";
  const legendaryCss = " bg-[#FFB13F] ";
  const rareCss = " bg-[#D28FD6]";
  cls += rareCharNames.includes(char) ? rareCss : legendaryCss;
  cls += onClick ? " hover:opacity-60 hover:cursor-pointer" : "";
  cls += custStyle !== "" ? " " + custStyle : "";

  return (
    <div className={cls} onClick={onClick ? onClick : () => {}}>
      <img
        src={"/api/assets/avatar/" + char + ".png"}
        alt={char}
        className="margin-auto"
      />
      {longName !== "" ? (
        <div>
          {tooLongNames.includes(char) ? (
            <div className="text-xs flex items-center justify-center text-center h-8 bg-slate-600 text-white">
              {longName}
            </div>
          ) : (
            <div className="text-xs flex items-center justify-center text-center h-8 bg-slate-600 text-white">
              {longName}
            </div>
          )}
        </div>
      ) : null}
    </div>
  );
}

const tooLongNames = [
  "Kaedehara Kazuha",
  "Kamisato Ayaka",
  "Raiden Shogun",
  "Sangonomiya Kokomi",
  "Kamisato Ayato",
  "Shikanoin Heizou",
  "Traveler (Anemo)",
  "Traveler (Geo)",
  "Traveler (Electro)",
  "Traveler (Pyro)",
];
const rareCharNames = [
  "amber",
  "barbara",
  "beidou",
  "bennett",
  "candace",
  "chongyun",
  "collei",
  "diona",
  "dori",
  "fischl",
  "gorou",
  "kaeya",
  "lisa",
  "kuki",
  "ningguang",
  "noelle",
  "razor",
  "heizou",
  "rosaria",
  "sara",
  "sucrose",
  "sayu",
  "thoma",
  "xiangling",
  "xinyan",
  "xingqiu",
  "yanfei",
  "yunjin",
];
