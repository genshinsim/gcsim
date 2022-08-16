import { useAppSelector } from '~src/store';

function CharCard({ char }: { char: string }) {
  return (
    <div className="p-2 hover:bg-gray-600 rounded-md hover:cursor-pointer">
      <img
        src={'/images/avatar/' + char + '.png'}
        alt={char}
        className="ml-auto h-32 wide:h-auto "
      />
    </div>
  );
}

//TODO: place holder for now
export function DatabaseCharacters() {
  const characters = useAppSelector((state) => state.db.characters);

  const cards = characters
    .filter((c) => true)
    .sort((a, b) => a.avatar_name.localeCompare(b.avatar_name))
    .map((c) => <CharCard char={c.avatar_name} key={c.avatar_name} />);

  return <div className="mt-2 grid grid-cols-8">{cards}</div>;
}
