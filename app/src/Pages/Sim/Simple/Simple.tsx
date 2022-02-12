import { Button, Callout, Card, Collapse, Intent } from "@blueprintjs/core";
import React from "react";
import { SectionDivider } from "~src/Components/SectionDivider";
import { useAppSelector } from "~src/store";
import { RootState } from "~src/store";
import { Main, ActionList } from "../Components";
import { Team } from "./Team";

export function Simple() {
  const [showActionList, setShowActionList] = React.useState<boolean>(true);
  const [showOptions, setShowOptions] = React.useState<boolean>(false);
  return (
    <Main className="flex flex-col gap-2">
      <div className="flex flex-col">
        <Team />
        <SectionDivider>Action List</SectionDivider>
        <div className="ml-auto mr-2">
          <Button
            icon="edit"
            onClick={() => setShowActionList(!showActionList)}
          >
            {showActionList ? "Hide" : "Show"}
          </Button>
        </div>
        <div className="pl-2 pr-2 pt-2">
          <Callout intent={Intent.PRIMARY} className="flex flex-col">
            Enter action list here. For more detailed on action list, see ??
            <br />
            <div className="ml-auto">
              <Button small>Hide all tips</Button>
            </div>
          </Callout>
        </div>
        <Collapse
          isOpen={showActionList}
          keepChildrenMounted
          className="basis-full flex flex-col"
        >
          <ActionList />
        </Collapse>
        <SectionDivider>Sim Options</SectionDivider>
        <div className="ml-auto mr-2">
          <Button icon="edit" onClick={() => setShowOptions(!showOptions)}>
            {showOptions ? "Hide" : "Show"}
          </Button>
        </div>
        <Collapse
          isOpen={showOptions}
          keepChildrenMounted
          className="basis-full flex flex-col"
        >
          <Card className="m-2">
            <div>character stat details in here</div>
          </Card>
        </Collapse>
      </div>
      <div className="sticky bottom-0 bg-bp-bg p-2 wide:ml-2 wide:mr-2 flex flex-row flex-wrap">
        <div className="ml-auto basis-full wide:basis-1/3">
          <Button icon="play" fill intent="primary">
            Run
          </Button>
        </div>
      </div>
    </Main>
  );
}
