import { Button, Checkbox, Switch } from "@blueprintjs/core";
import { Character } from "~/src/types";

type Props = {
  char: Character;
  onChange: (char: Character) => void;
};

export function CharacterEditArtifactSets(props: Props) {
  const arts: JSX.Element[] = [];
  for (const key in props.char.sets) {
    arts.push(
      <div className="flex flex-row rounded-md" key={key}>
        <img key="key" src={`/images/artifacts/${key}_flower.png`} />
        <span className="text-center text-xs">{props.char.sets[key]}</span>
      </div>
    );
  }

  return (
    <div className="flex flex-col place-items-center">
      <div className="flex flex-row flex-wrap gap-2 justify-center w-full">
        <div className="basis-full sm:basis-320 rounded-md bg-gray-600 flex flex-row place-items-center pl-1 pr-2">
          <img
            key="key"
            src={`/images/artifacts/emblemofseveredfate_flower.png`}
            className="w-12"
          />
          <span className="font-bold">Set Bonus:</span>
          <div className="flex flex-row gap-2 flex-grow justify-center">
            <Checkbox large style={{ marginBottom: 0 }}>
              2
            </Checkbox>
            <Checkbox large style={{ marginBottom: 0 }}>
              4
            </Checkbox>
          </div>
          <Button className="ml-auto" icon="trash" intent="danger" />
        </div>
        <div className="basis-full sm:basis-320 rounded-md bg-gray-600 flex flex-row place-items-center pl-1 pr-2">
          <img
            key="key"
            src={`/images/artifacts/emblemofseveredfate_flower.png`}
            className="w-12"
          />
          <span className="font-bold">Set Bonus:</span>
          <div className="flex flex-row gap-2 flex-grow justify-center">
            <Checkbox large style={{ marginBottom: 0 }}>
              2
            </Checkbox>
            <Checkbox large style={{ marginBottom: 0 }}>
              4
            </Checkbox>
          </div>
          <Button className="ml-auto" icon="trash" intent="danger" />
        </div>
      </div>
      <div className="mt-2 w-full xs:w-[25rem]">
        <Button icon="add" fill intent="success">
          Add Set Bonus
        </Button>
      </div>
    </div>
  );
}
