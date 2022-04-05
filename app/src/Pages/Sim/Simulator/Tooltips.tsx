import { Callout, Intent, Button } from "@blueprintjs/core";
import React from "react";
import { Trans } from "react-i18next";
import { useAppSelector, RootState, useAppDispatch } from "~src/store";
import { VideoPlayer } from "../Components";
import { simActions } from "../simSlice";

export const ActionListTooltip = () => {
  const { showTips } = useAppSelector((state: RootState) => {
    return {
      showTips: state.sim.showTips,
    };
  });
  const dispatch = useAppDispatch();

  const toggleTips = () => {
    dispatch(simActions.setShowTips(!showTips));
  };

  if (!showTips) {
    return null;
  }

  return (
    <div className="pl-2 pr-2 pt-2">
      <Callout intent={Intent.PRIMARY} className="flex flex-col">
        <p>
          <Trans>simple.discord_pre</Trans>
          <a href="https://discord.gg/W36ZwwhEaG" target="_blank">
            Discord
          </a>
          <Trans>simple.discord_post</Trans>
        </p>
        <p>
          <Trans>simple.documentation_pre</Trans>
          <a
            href="https://docs.gcsim.app/guide/sequential_mode"
            target="_blank"
          >
            <Trans>simple.documentation</Trans>
          </a>
          <Trans>simple.documentation_post</Trans>
        </p>
        <div className="ml-auto">
          <Button small onClick={toggleTips}>
            <Trans>simple.hide_all_tips</Trans>
          </Button>
        </div>
      </Callout>
    </div>
  );
};

export const TeamBuilderTooltip = () => {
  const { showTips } = useAppSelector((state: RootState) => {
    return {
      showTips: state.sim.showTips,
    };
  });
  const [openAddCharHelp, setOpenAddCharHelp] = React.useState<boolean>(false);

  const dispatch = useAppDispatch();

  const toggleTips = () => {
    dispatch(simActions.setShowTips(!showTips));
  };

  if (!showTips) {
    return null;
  }

  return (
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
          <Button small onClick={toggleTips}>
            <Trans>simple.hide_all_tips</Trans>
          </Button>
        </div>
      </Callout>
      <VideoPlayer
        url="/videos/add-character.webm"
        isOpen={openAddCharHelp}
        onClose={() => setOpenAddCharHelp(false)}
        title="Adding a character"
      />
    </div>
  );
};
