import { useLocation } from "wouter";
import { Viewer } from "~src/Components/Viewer";
import { Viewport } from "~src/Components/Viewport";
import { RootState, useAppDispatch, useAppSelector } from "~src/store";
import Dropzone from "./Dropzone";
import Shared from "./Shared";
import { viewerActions } from "./viewerSlice";
import { Trans, useTranslation } from "react-i18next";

type Props = {
  path: string;
  next?: boolean;
};

export function ViewerDash({ path, next = false }: Props) {
  useTranslation();

  const [_, setLocation] = useLocation();
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
    return (
      <Shared
        next={next}
        path={path}
        handleClose={() => {
          dispatch(viewerActions.setSelected(""));
          setLocation("/viewer");
        }}
      />
    );
  }
  //show viewer if selected != -1
  if (selected !== "") {
    return (
      <div className="flex-grow">
        <Viewer
          data={data[selected]}
          className="h-full"
          handleClose={() => {
            dispatch(viewerActions.setSelected(""));
          }}
        />
      </div>
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
      <div className="font-bold">
        <Trans>viewerdashboard.upload_a_file</Trans>
      </div>
      <Dropzone />
      <div className="font-bold mb-2">
        <Trans>viewerdashboard.or_select_from</Trans>
      </div>
      {rows}
    </Viewport>
  );
}
