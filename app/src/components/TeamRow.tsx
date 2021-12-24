import React from "react";
import CopyIcon from "../content_copy_white_24dp.svg";
// import EyeIcon from "../visibility_white_24dp.svg";
import charData from "../data/character_images.json";

export default function TeamRow({
  title,
  version,
  author,
  description,
  characters,
  config,
}: {
  title: string;
  version: string;
  author: string;
  description: string;
  characters: string[];
  config: string;
}) {
  const [open, setOpen] = React.useState<boolean>(false);

  return (
    <div className="m-2 p-2 rounded-md bg-gray-600 gap-1 items-center grid lg:grid-cols-3 md:grid-cols-2 sm:grid-cols-1 lg:text-base md:text-sm sm:text-xs">
      <div>
        <div className="font-bold xl:text-lg mb-1">{title}</div>
        <div className="grid grid-cols-4">
          {characters.map((c: string) => {
            console.log(c);
            // @ts-ignore: Unreachable code error
            let image = charData[c];
            return (
              <div key={c} className="h-24">
                <img src={image} alt={c} className="object-contain h-full" />
              </div>
            );
          })}
        </div>
      </div>

      <div
        className="flex-grow flex flex-row items-center col-span-2"
        style={{ position: "relative" }}
      >
        <div className="flex flex-col">
          <div>
            <strong>Author: </strong>
            {author}
          </div>
          <div>
            <strong>Version: </strong>
            {version}
          </div>
          <div>
            <strong>Description: </strong>
            {description}
          </div>
        </div>

        <img
          src={CopyIcon}
          alt="copy"
          className="p-1 rounded-md hover:bg-gray-500"
          onClick={() => {
            navigator.clipboard.writeText(config).then(
              () => {
                alert("Configuration copied! Paste it in gcsim to run.");
              },
              () => {
                alert("Error copying :( Not sure what went wrong");
              }
            );
          }}
          style={{
            position: "absolute",
            top: "0",
            right: "0",
          }}
        />
      </div>
    </div>
  );
}

const customStyles = {
  content: {
    top: "50%",
    left: "50%",
    right: "auto",
    bottom: "auto",
    marginRight: "-50%",
    transform: "translate(-50%, -50%)",
  },
};
