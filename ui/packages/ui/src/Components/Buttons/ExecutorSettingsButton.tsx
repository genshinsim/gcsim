import { Button } from "@blueprintjs/core";
import { appActions } from "../../Stores/appSlice";
import { useAppDispatch } from "../../Stores/store";

// TODO: translation
export default ({}) => {
  const dispatch = useAppDispatch();

  return (
    <Button
        icon="cog"
        text="Settings"
        onClick={() => dispatch(appActions.setSettingsOpen(true))} />
  );
};