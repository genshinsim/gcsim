import React from "react";
import { Callout, Intent, Button } from "@blueprintjs/core";
import { Trans } from "react-i18next";
import { useAppDispatch, useAppSelector } from "../../Stores/store";
import { userActions } from "../../Stores/userSlice";

export const ActionListTooltip = () => {
  const settings = useAppSelector((state) => state.user.data.settings);
  const dispatch = useAppDispatch();
  const toggleTips = () => {
    dispatch(
      userActions.setUserSettings({
        showTips: !settings.showTips,
        showBuilder: settings.showBuilder,
      })
    );
  };

  if (!settings.showTips) {
    return null;
  }

  return (
    <div className="pl-2 pr-2 pt-2">
      <Callout intent={Intent.PRIMARY} className="flex flex-col">
        <p>
          <Trans>simple.discord_pre</Trans>
          <a href="https://discord.gg/W36ZwwhEaG" target="_blank" rel="noreferrer">
            Discord
          </a>
          <Trans>simple.discord_post</Trans>
        </p>
        <p>
          <Trans>simple.documentation_pre</Trans>
          <a href="https://docs.gcsim.app/guides" target="_blank" rel="noreferrer">
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
  const [openAddCharHelp, setOpenAddCharHelp] = React.useState<boolean>(false);

  const settings = useAppSelector((state) => state.user.data.settings);
  const dispatch = useAppDispatch();
  const toggleTips = () => {
    dispatch(
      userActions.setUserSettings({
        showTips: !settings.showTips,
        showBuilder: settings.showBuilder,
      })
    );
  };

  if (!settings.showTips) {
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
    </div>
  );
};
