import { model } from "@gcsim/types";
import { AvatarPortrait } from "../AvatarPortait/AvatarPortrait";

type Props = {
  //array of characters to display
  chars?: (model.Character | null)[];
  invalid?: string[];
  className?: string;
  onImageLoaded?: () => void;
};

export const AvatarCard = ({
  chars,
  invalid,
  className = "",
  onImageLoaded = () => {},
}: Props) => {
  if (chars?.length === 0) {
    return <div> no data</div>;
  }
  return (
    <div
      className={
        "flex flex-row flex-wrap justify-center gap-2" +
        (className == "" ? "" : " " + className)
      }
    >
      {chars?.map((c, i) => {
        let name = "";
        if (c != null) {
          name = c.name ?? "";
        }
        return (
          <AvatarPortrait
            key={"char-" + i}
            i={i}
            char={c}
            className="min-[420px]:max-w-[96px] min-[420px]:basis-1/4 basis-[45%]"
            invalid={invalid != null && invalid?.includes(name)}
            onImageLoaded={onImageLoaded}
          />
        );
      })}
    </div>
  );
};
