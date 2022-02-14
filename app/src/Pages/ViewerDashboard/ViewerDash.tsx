import { Viewer } from "~src/Components/Viewer";
import { Viewport } from "~src/Components/Viewport";
import { useAppSelector, RootState, useAppDispatch } from "~src/store";
import Dropzone from "./Dropzone";
import Shared from "./Shared";
import { viewerActions } from "./viewerSlice";

type Props = {
  path: string;
};

export function ViewerDash({ path }: Props) {
  const { data, selected } = useAppSelector((state: RootState) => {
    return {
      data: state.viewer.data,
      selected: state.viewer.selected,
    };
  });
  const dispatch = useAppDispatch();

  //if path is not "/" then load the shared view
  if (path !== "/") {
    //need a check here to make sure this doesn't already exists
    return <Shared path={path} />;
  }
  //show viewer if selected != -1
  if (selected !== "") {
    return (
      <Viewport className="flex-grow">
        <Viewer
          data={data[selected]}
          className="h-full"
          handleClose={() => {
            dispatch(viewerActions.setSelected(""));
          }}
        />
      </Viewport>
    );
  }

  let rows: JSX.Element[] = [];

  for (const key in data) {
    rows.push(
      <div
        className="p-2 bg-gray-600 rounded-md m-1 hover:bg-gray-800 hover:cursor-pointer"
        key={key}
        onClick={() => {
          dispatch(viewerActions.setSelected(key));
        }}
      >
        {key}
      </div>
    );
  }

  //show dash board otherwise
  return (
    <Viewport className="flex flex-col p-1">
      <div className="font-bold">Upload a file</div>
      <Dropzone />
      <div className="font-bold mb-2">
        Or select from the following previously opened files:
      </div>
      {rows}
    </Viewport>
  );
}
