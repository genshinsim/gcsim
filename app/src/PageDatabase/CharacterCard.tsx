type CharacterCardProps = {
  char: string;
  style?: string;
};
export function CharacterCard({ char, style = '' }: CharacterCardProps) {
  return (
    <div
      className={
        'p-2 hover:bg-gray-600 rounded-md' + style !== '' ? ' ' + style : ''
      }
    >
      <img
        src={'/api/assets/avatar/' + char + '.png'}
        alt={char}
        className="ml-auto h-32 wide:h-auto "
      />
    </div>
  );
}
