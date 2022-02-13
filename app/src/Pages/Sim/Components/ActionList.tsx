import { TextArea } from "@blueprintjs/core";
import { useAppSelector, RootState, useAppDispatch } from "~src/store";
import { simActions } from "..";

export function ActionList() {
  const { cfg } = useAppSelector((state: RootState) => {
    return {
      cfg: state.sim.cfg,
    };
  });
  const dispatch = useAppDispatch();
  return (
    <div className="p-1 md:p-2">
      <TextArea
        rows={10}
        fill
        value={cfg}
        onChange={(e) => {
          dispatch(simActions.setCfg(e.target.value));
        }}
      ></TextArea>
    </div>
  );
}
