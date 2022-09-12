import { useLocation } from "wouter";
import { CharEntry } from "./PageDatabase";

function CharCard({ charEntry }: { charEntry: CharEntry }) {
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
    "chongyun",
    "diona",
    "fischl",
    "gorou",
    "kaeya",
    "kujousara",
    "lisa",
    "kuki",
    "ningguang",
    "noelle",
    "razor",
    "heizou",
    "rosaria",
    "sucrose",
    "sayu",
    "thoma",
    "xiangling",
    "xinyan",
    "xingqiu",
    "yanfei",
    "yunjin",
  ];

  const [_, setLocation] = useLocation();
  const [shortName, name] = charEntry;

  const legendaryCss = "bg-opacity-60 bg-[#FFB13F] ";
  const rareCss = "bg-opacity-60 bg-[#D28FD6]";
  return (
    <div className=" border-gray-700 border-2  rounded-md">
      <div className="hover:opacity-50">
        <div
          className={rareCharNames.includes(shortName) ? rareCss : legendaryCss}
          onClick={() => setLocation(`/db/${shortName}`)}
        >
          <img
            src={`/images/avatar/${shortName}.png`}
            alt={name}
            className="margin-auto"
          />
        </div>
        <div>
          {tooLongNames.includes(name) ? (
            <div className="text-xs flex items-center justify-center text-center h-8 bg-slate-600 ">
              {name}
            </div>
          ) : (
            <div className="text-md flex items-center justify-center text-center h-8 bg-slate-600 ">
              {name}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export function CharsView({ characters }: { characters: CharEntry[] }) {
  return (
    <div className="p-4">
      <div className="grid grid-cols-3 gap-2 wide:grid-cols-12">
        {characters.map((entry) => (
          <div key={entry[0]}>
            <CharCard charEntry={entry} />
          </div>
        ))}
      </div>
    </div>
  );
}
