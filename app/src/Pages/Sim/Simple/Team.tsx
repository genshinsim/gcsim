import { Callout, Intent, Button, Card, Tabs, Tab } from "@blueprintjs/core";
import React from "react";
import { CharDetail } from "~src/Components/Character";
import { charTestConfig } from "..";
import {
  CharacterCardView,
  CharacterStats,
  SectionDivider,
} from "../Components";

type Props = {
  chars: CharDetail[];
};

export function Team(props: Props) {
  const [showTeamEdit, setShowTeamEdit] = React.useState<boolean>(false);

  const tabs = props.chars.map((c) => {
    return (
      <Tab
        key={c.name}
        id={c.name}
        title={c.name}
        className={"focus:outline-none"}
        panel={
          <CharacterStats
            char={c}
            onChange={(index, value) => {
              console.log(index, value);
            }}
          />
        }
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
      <div className={showTeamEdit ? "hidden" : "mt-2"}>
        <CharacterCardView chars={charTestConfig} />
      </div>
      <div className={showTeamEdit ? "" : "hidden"}>
        <Card className="m-2">
          <Tabs className="capitalize">{tabs}</Tabs>
        </Card>
      </div>
      <div className="ml-auto mr-2">
        <Button icon="edit" onClick={() => setShowTeamEdit(!showTeamEdit)}>
          {showTeamEdit ? "Done" : "Edit"}
        </Button>
      </div>
    </div>
  );
}
