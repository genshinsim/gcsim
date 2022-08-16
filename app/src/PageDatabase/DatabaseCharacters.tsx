import { useLocation } from 'wouter';
import { useAppSelector } from '~src/store';
import { CharacterCard } from './CharacterCard';

//TODO: place holder for now
export function DatabaseCharacters() {
  const characters = useAppSelector((state) => state.db.characters);
  const [_, setLocation] = useLocation();

  const cards = characters
    .map((x) => x)
    .sort((a, b) => a.avatar_name.localeCompare(b.avatar_name))
    .map((c) => (
      <div
        key={c.avatar_name}
        onClick={() => {
          setLocation(`/db/${c.avatar_name}`);
        }}
      >
        <CharacterCard char={c.avatar_name} />
      </div>
    ));

  console.log(characters);

  return <div className="mt-2 grid grid-cols-8">{cards}</div>;
}
