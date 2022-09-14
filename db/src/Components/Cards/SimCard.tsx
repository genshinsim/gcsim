import { Metadata } from "Types/stats";
import { CharacterCard } from "./CharacterCard";

type SimCardProps = {
  meta: Metadata;
  onCharacterClick?: (char: string) => void;
  onClick?: () => void;
  summary?: React.ReactNode;
  actions?: React.ReactNode;
};

export function SimCard(props: SimCardProps) {
  const chars = props.meta.char_details.map((char) => {
    return (
      <div className="rounded-md relative p-1">
        <CharacterCard
          char={char.name}
          onClick={
            props.onCharacterClick
              ? () => {
                  props.onCharacterClick!(char.name);
                }
              : undefined
          }
        />
        <div className=" absolute top-0 right-[5px] text-sm font-semibold text-white">{`${char.cons}`}</div>
      </div>
    );
  });
  return (
    <div className="flex flex-row flex-wrap sm:flex-nowrap gap-y-1 w-full m-2 p-2 rounded-md bg-gray-700 place-items-center">
      <div className="flex flex-col sm:basis-1/4 xs:basis-full">
        <div className="grid grid-cols-4">{chars}</div>
        <div className="hidden basis-0 lg:block md:flex-1"></div>
      </div>
      <div className=" flex-1 overflow-hidden mb-auto pl-2 hidden lg:block"></div>
      <div className="ml-auto flex flex-col mr-4 md:basis-60 basis-full">
        {props.summary ? props.summary : null}
      </div>
      <div>{props.actions ? props.actions : null}</div>
    </div>
  );
}
