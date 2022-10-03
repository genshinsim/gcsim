import React from "react";
import { useLocation } from "wouter";
import { useAppDispatch, useAppSelector } from "Store";
import { dbActions } from "../../Store/dbSlice";
import axios from "axios";
import pako from "pako";
import { DBAvatarSimDetails } from "Types/database";
import { CharacterCard, Spinner } from "Components";

type CharacterViewProps = {
  char: string;
};

type statusType = "idle" | "loading" | "done" | "error";

export function Teams({ char }: CharacterViewProps) {
  const [status, setStatus] = React.useState<statusType>("idle");
  const [errMsg, setErrMsg] = React.useState<string>("");
  const charSims = useAppSelector((state) => state.db.charSims);
  const dispatch = useAppDispatch();
  const [, setLocation] = useLocation();

  React.useEffect(() => {
    if (status === "idle") {
      axios
        .get(`/api/db/${char}`)
        .then((resp) => {
          console.log(resp.data);
          let next: DBAvatarSimDetails[] = [];
          resp.data.forEach((e: any) => {
            const binaryStr = Uint8Array.from(window.atob(e.config_hash), (v) =>
              v.charCodeAt(0)
            );
            const restored = pako.inflate(binaryStr, { to: "string" });
            next.push({
              ...e,
              metadata: JSON.parse(e.metadata),
              create_time: Math.floor(new Date(e.create_time).getTime() / 1000),
              config: restored,
            });
          });
          console.log(next);
          dispatch(dbActions.setCharSimList({ char: char, data: next }));
          setStatus("done");
        })
        .catch((err) => {
          setStatus("error");
          setErrMsg(`Error encountered loading sims for ${char}: ${err}`);
        });
    }
  }, [status, char, dispatch]);

  switch (status) {
    case "loading":
    case "idle":
      return (
        <div className="flex flex-row place-content-center mt-2">
          <Spinner />
        </div>
      );
    case "error":
      return (
        <div className="flex flex-row place-content-center mt-2">{errMsg}</div>
      );
  }

  if (!(char in charSims) || charSims[char].length === 0) {
    return (
      <div className="flex flex-row place-content-center mt-2">
        Sorry, this character does not have any sims submitted :(
      </div>
    );
  }

  const teams = charSims[char].reduce<{ [key in string]: number }>(
    (next, s) => {
      const key = s.metadata.char_names
        .map((x: string) => x)
        .sort()
        .join("-");
      next[key]++;
      return next;
    },
    {}
  );

  const rows = Object.keys(teams)
    .sort()
    .map((e) => {
      const col = e
        .split("-")
        .map((char) => (
          <CharacterCard char={char} key={char} custStyle="m-1" />
        ));
      return (
        <div
          key={e}
          className="grid grid-cols-4 bg-gray-700 rounded-md p-2 hover:cursor-pointer border-2 border-gray-700 hover:border-white"
          onClick={() => setLocation(`/db/${char}/${e}`)}
        >
          {col}
        </div>
      );
    });

  return (
    <main className="flex flex-col h-full m-2 w-full xs:w-full sm:w-[640px] hd:w-full wide:w-[1160px] ml-auto mr-auto ">
      <div className="text-white font-bold mb-2 text-xl">
        <span style={{ textTransform: "capitalize" }}>{char}</span> Teams
      </div>
      <div className="grid grid-cols-2 gap-3">{rows}</div>
    </main>
  );
}
