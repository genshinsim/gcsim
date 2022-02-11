import { Callout, Intent, Button, Card, Tabs, Tab } from "@blueprintjs/core";
import React from "react";
import { CharacterEdit, CharDetail } from "~src/Components/Character";
import { charTestConfig } from "..";
import { CharacterCardView, SectionDivider } from "../Components";

type Props = {
  chars: CharDetail[];
};

export function Team(props: Props) {
  const [showTeamEdit, setShowTeamEdit] = React.useState<boolean>(false);
  const [edit, setEdit] = React.useState<number>(-1);

  const handleEdit = (index: number) => {
    return () => {
      if (index > -1 && index < props.chars.length) {
        setEdit(index);
        setShowTeamEdit(true);
        console.log("editing: " + index);
      }
    };
  };

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
      <div className={showTeamEdit ? "hidden" : "mt-2"}>
        <CharacterCardView chars={charTestConfig} handleEdit={handleEdit} />
      </div>
      {showTeamEdit ? (
        <Card className="m-2">
          <CharacterEdit
            char={props.chars[edit]}
            onChange={(char) => console.log("editing " + char.name)}
          />
        </Card>
      ) : null}
      <div className={showTeamEdit ? "ml-auto mr-2" : "hidden"}>
        <Button
          intent="primary"
          icon="edit"
          onClick={() => setShowTeamEdit(false)}
        >
          Done
        </Button>
      </div>
    </div>
  );
}
