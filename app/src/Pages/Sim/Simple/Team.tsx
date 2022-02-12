import { Callout, Intent, Button, Card } from "@blueprintjs/core";
import React from "react";
import {
  CharacterCard,
  CharacterEdit,
  ConsolidateCharStats,
} from "~src/Components/Character";
import { SectionDivider } from "~src/Components/SectionDivider";
import { CharacterCardView } from "../Components";
import { useAppDispatch, useAppSelector } from "~src/store";
import { RootState } from "~src/store";
import { simActions } from "~src/Pages/Sim/simSlice";
import { Character } from "~src/types";

type Props = {};

export function Team(props: Props) {
  const { team, edit_index } = useAppSelector((state: RootState) => {
    return {
      team: state.sim.team,
      edit_index: state.sim.edit_index,
    };
  });
  const dispatch = useAppDispatch();
  const myRef = React.useRef<HTMLSpanElement>(null);
  React.useEffect(() => {
    executeScroll();
  }, [edit_index]);

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
  const handleChange = (index: number) => {
    return (char: Character) => {
      if (index > -1 && index < team.length) {
        dispatch(simActions.setCharacter({ char: char, index: index }));
      }
    };
  };
  const executeScroll = () => {
    if (myRef.current) {
      myRef.current.scrollIntoView({ behavior: "smooth" });
    }
  };

  const teamStats = ConsolidateCharStats(team);

  const cards = team.map((c, index) => {
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
        className="basis-full md:basis-1/2 wide:basis-1/4 pt-2 pr-2 pb-2"
      />
    );
  });

  return (
    <div className="flex flex-col">
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
      <span ref={myRef} />
      <div className={edit_index > -1 ? "hidden" : "mt-2"}>
        <div className="flex flex-row flex-wrap pl-2">{cards}</div>
      </div>
      {edit_index > -1 ? (
        <Card className="m-2">
          <CharacterEdit
            char={team[edit_index]}
            onChange={handleChange(edit_index)}
          />
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
      ) : null}
    </div>
  );
}
