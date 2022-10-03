import { CharacterCard } from "Components";
import { useLocation } from "wouter";

export function CharsGrid({ characters }: { characters: string[][] }) {
  const [, setLocation] = useLocation();
  return (
    <div className="p-4">
      <div className="grid grid-cols-3 gap-2 sm:grid-cols-4 md:grid-cols-8 wide:grid-cols-12">
        {characters.map((entry) => (
          <div key={entry[0]}>
            <CharacterCard
              custStyle="border-gray-700 border-2"
              char={entry[0]}
              longName={entry[1]}
              onClick={() => setLocation(`/db/${entry[0]}`)}
            />
          </div>
        ))}
      </div>
    </div>
  );
}
