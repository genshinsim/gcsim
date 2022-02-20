import { TextArea } from "@blueprintjs/core";
import { useAppSelector, RootState, useAppDispatch } from "~src/store";
import { simActions } from "..";

type Props = {
  cfg: string;
  onChange: (v: string) => void;
};

export function ActionList(props: Props) {
  const dispatch = useAppDispatch();
  return (
    <div className="p-1 md:p-2">
      <TextArea
        rows={30}
        fill
        value={props.cfg}
        onChange={(e) => {
          props.onChange(e.target.value);
        }}
      ></TextArea>
    </div>
  );
}
