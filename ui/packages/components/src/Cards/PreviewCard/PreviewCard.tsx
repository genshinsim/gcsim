import { model } from "@gcsim/types";
import { AvatarPortrait } from "../AvatarPortait/AvatarPortrait";
import { Graphs } from "./Graphs";
import { Metadata } from "./Metadata";

type PreviewCardProps = {
  data: model.SimulationResult;
  onImageLoaded?: () => void;
};

export const PreviewCard = ({
  data,
  onImageLoaded = () => {},
}: PreviewCardProps) => {
  return (
    <div className="!w-[540px] !h-[250px] bg-slate-800">
      <div className="flex flex-col h-full">
        <div className="grid grid-cols-4">
          {data.character_details?.map((c, i) => {
            return (
              <AvatarPortrait
                key={"char-" + i}
                i={i}
                char={c}
                invalid={
                  data.incomplete_characters?.includes(c.name ?? "") ?? false
                }
                className="m-1"
                onImageLoaded={onImageLoaded}
              />
            );
          })}
        </div>
        <Metadata data={data} />
        <Graphs data={data} className="grow" />
      </div>
    </div>
  );
};
