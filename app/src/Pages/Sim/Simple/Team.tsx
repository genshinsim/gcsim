import { Callout, Intent, Button, Card, useHotkeys } from "@blueprintjs/core";
import React from "react";
import {
  CharacterCard,
  CharacterSelect,
  ConsolidateCharStats,
  ICharacter,
} from "~src/Components/Character";
import { SectionDivider } from "~src/Components/SectionDivider";
import { useAppDispatch, useAppSelector } from "~src/store";
import { RootState } from "~src/store";
import { simActions } from "~src/Pages/Sim/simSlice";
import { CharacterEdit } from "./CharacterEdit";

export function Team() {
  const { team, edit_index } = useAppSelector((state: RootState) => {
    return {
      team: state.sim.team,
      edit_index: state.sim.edit_index,
    };
  });
  const dispatch = useAppDispatch();
  const [open, setOpen] = React.useState<boolean>(false);
  const myRef = React.useRef<HTMLSpanElement>(null);
  React.useEffect(() => {
    if (edit_index > -1) {
      executeScroll();
    }
  }, [edit_index]);

  const hotkeys = React.useMemo(
    () => [
      {
        combo: "Esc",
        global: true,
        label: "Exit edit",
        onKeyDown: () => {
          dispatch(simActions.editCharacter({ index: -1 }));
        },
      },
    ],
    []
  );
  useHotkeys(hotkeys);

  const handleEdit = (index: number) => {
    return () => {
      if (index > -1 && index < team.length) {
        dispatch(simActions.editCharacter({ index: index }));
      }
    };
  };
  const handleDelete = (index: number) => {
    return () => {
      if (index > -1 && index < team.length) {
        dispatch(simActions.deleteCharacter({ index: index }));
      }
    };
  };
  const executeScroll = () => {
    if (myRef.current) {
      myRef.current.scrollIntoView({ behavior: "smooth" });
    }
  };
  const handleAddCharacter = (w: ICharacter) => {
    setOpen(false);
    dispatch(simActions.addCharacter({ name: w.key }));
  };

  let disabled: string[] = [];
  let cards: JSX.Element[] = [];

  if (team.length > 0) {
    const teamStats = ConsolidateCharStats(team);

    // console.log(team);
    // console.log(teamStats);

    disabled = team.map((c) => c.name);

    cards = team.map((c, index) => {
      return (
        <CharacterCard
          key={c.name}
          char={c}
          stats={teamStats.stats[c.name]}
          statsRows={teamStats.maxRows}
          showDelete
          showEdit
          toggleEdit={handleEdit(index)}
          handleDelete={handleDelete(index)}
          className="basis-full sm:basis-1/2 hd:basis-1/4 pt-2 pr-2 pb-2"
        />
      );
    });
  }

  return (
    <div className="flex flex-col">
      <span ref={myRef} />
      <SectionDivider>Team</SectionDivider>
      <div className="pl-2 pr-2">
        <Callout intent={Intent.PRIMARY} className="flex flex-col">
          Enter your team information in this section
          <br />
          <div className="ml-auto">
            <Button small>Hide all tips</Button>
          </div>
        </Callout>
      </div>
      {team.length == 0 ? (
        <div className="p-4 bg-gray-700 rounded-md mt-2 ml-2 mr-2 text-center font-bold">
          Start by adding some team members
        </div>
      ) : null}
      <div className={edit_index > -1 ? "hidden" : "mt-2"}>
        <div className="flex flex-row flex-wrap pl-2">{cards}</div>
      </div>
      {edit_index > -1 ? (
        <Card className="m-2">
          <CharacterEdit />
          <Button
            fill
            intent="primary"
            icon="edit"
            onClick={() => {
              dispatch(simActions.editCharacter({ index: -1 }));
            }}
          >
            Done
          </Button>
        </Card>
      ) : (
        <div className={team.length >= 4 ? "hidden" : "mt-2 pl-2 pr-2"}>
          <Button
            fill
            icon="add"
            intent="primary"
            onClick={() => setOpen(true)}
          >
            Add Character
          </Button>
        </div>
      )}
      <CharacterSelect
        disabled={disabled}
        onClose={() => setOpen(false)}
        onSelect={handleAddCharacter}
        isOpen={open}
      />
    </div>
  );
}
