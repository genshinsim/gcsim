import { Callout, Intent, Button, Card, ButtonGroup } from "@blueprintjs/core";
import React from "react";
import {
  CharacterCard,
  characterKeyToICharacter,
  CharacterSelect,
  ConsolidateCharStats,
  ICharacter,
  newChar,
} from "~src/Components/Character";
import { SectionDivider } from "~src/Components/SectionDivider";
import { RootState, useAppDispatch, useAppSelector } from "~src/store";
import { simActions } from "~src/Pages/Sim/simSlice";
import { CharacterEdit } from "./CharacterEdit";
import { LoadGOOD, VideoPlayer } from "../Components";
import { Trans, useTranslation } from "react-i18next";

export function Team() {
  useTranslation();

  const { team, edit_index, showTips, imported } = useAppSelector(
    (state: RootState) => {
      return {
        team: state.sim.team,
        edit_index: state.sim.edit_index,
        showTips: state.sim.showTips,
        imported: state.user_data.GOODImport,
      };
    }
  );
  const dispatch = useAppDispatch();
  const [open, setOpen] = React.useState<boolean>(false);
  const [openImport, setOpenImport] = React.useState<boolean>(false);
  const [openAddCharHelp, setOpenAddCharHelp] = React.useState<boolean>(false);
  const myRef = React.useRef<HTMLSpanElement>(null);
  React.useEffect(() => {
    if (edit_index > -1) {
      executeScroll();
    }
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
  const executeScroll = () => {
    if (myRef.current) {
      myRef.current.scrollIntoView({ behavior: "smooth" });
    }
  };
  const handleAddCharacter = (character: ICharacter) => {
    setOpen(false);
    //check if this is from GOOD
    if (character.notes) {
      dispatch(simActions.addCharacter({ character: imported[character.key] }));
      return;
    }
    //else it's new
    const c = newChar(character.key);
    dispatch(simActions.addCharacter({ character: c }));
  };

  const hideTips = () => {
    dispatch(simActions.setShowTips(false));
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

  const additionalChars = Object.keys(imported).map((k) => {
    let x = Object.assign({}, characterKeyToICharacter[k]);
    x.notes = `Imported on ${imported[k].date_added}`;
    return x;
  });

  return (
    <div className="flex flex-col">
      <span ref={myRef} />
      <SectionDivider>
        <Trans>simple.team</Trans>
      </SectionDivider>
      {showTips ? (
        <div className="pl-2 pr-2">
          <Callout intent={Intent.PRIMARY} className="flex flex-col">
            <span>
              <Trans>simple.video_pre</Trans>
              <a onClick={() => setOpenAddCharHelp(true)}>
                <Trans>simple.video</Trans>
              </a>
              <Trans>simple.video_post</Trans>
            </span>
            <div className="ml-auto">
              <Button small onClick={hideTips}>
                <Trans>simple.hide_all_tips</Trans>
              </Button>
            </div>
          </Callout>
        </div>
      ) : null}
      {team.length == 0 ? (
        <div className="p-4 bg-gray-700 rounded-md mt-2 ml-2 mr-2 text-center font-bold">
          <Trans>simple.start_by_adding</Trans>
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
            <Trans>simple.done</Trans>
          </Button>
        </Card>
      ) : (
        <div className="mt-2 pl-2 pr-2">
          <ButtonGroup fill>
            <Button onClick={() => setOpenImport(true)}>Import Data</Button>
            <Button
              icon="add"
              intent="primary"
              onClick={() => setOpen(true)}
              disabled={team.length >= 4}
            >
              <Trans>simple.add_character</Trans>
            </Button>
          </ButtonGroup>
        </div>
      )}
      <CharacterSelect
        disabled={disabled}
        onClose={() => setOpen(false)}
        onSelect={handleAddCharacter}
        isOpen={open}
        additionalOptions={additionalChars}
      />
      <LoadGOOD isOpen={openImport} onClose={() => setOpenImport(false)} />
      <VideoPlayer
        url="/videos/add-character.webm"
        isOpen={openAddCharHelp}
        onClose={() => setOpenAddCharHelp(false)}
        title="Adding a character"
      />
    </div>
  );
}
