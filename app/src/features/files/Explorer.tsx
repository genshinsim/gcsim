import React from "react";
import produce from "immer";

import { useDispatch, useSelector } from "react-redux";
import { sendMessage, setCurrentPath } from "app/appSlice";

import {
  Breadcrumbs,
  Button,
  ButtonGroup,
  Card,
  Dialog,
  BreadcrumbProps,
  Classes,
  InputGroup,
  TreeNodeInfo,
  Tree,
  Callout,
  Intent,
} from "@blueprintjs/core";
import { setConfig, setHasChange } from "features/sim/simSlice";
import { RootState } from "app/store";

interface IExplorerState {
  currentPath: string;
  isLoading: boolean;
  files: IFile[];
  selected: string;
}

interface IFile {
  name: string;
  is_dir: boolean;
  path: string;
}

type Action =
  | { type: "loading"; payload: boolean }
  | { type: "files"; payload: IFile[] }
  | { type: "path"; payload: string }
  | { type: "selected"; payload: string }
  | { type: "showNewFile"; payload: boolean };

const initialState: IExplorerState = {
  currentPath: "",
  isLoading: false,
  files: [],
  selected: "",
};

function reducer(state: IExplorerState, action: Action) {
  return produce(state, (next) => {
    switch (action.type) {
      case "files":
        next.files = action.payload;
        return;
      case "loading":
        next.isLoading = action.payload;
        return;
      case "path":
        next.currentPath = action.payload;
        return;
      case "selected":
        next.selected = action.payload;
        return;
      default:
        return;
    }
  });
}

function Explorer() {
  const { hasChange } = useSelector((state: RootState) => {
    return {
      hasChange: state.sim.hasChange,
    };
  });
  const storeDispatch = useDispatch();

  const [state, dispatch] = React.useReducer(reducer, initialState);
  const [selected, setSelected] = React.useState<string>("");
  const [showNewFile, setShowNewFile] = React.useState<boolean>(false);
  const [newFileName, setNewFileName] = React.useState<string>("");
  const [showNewFolder, setShowNewFolder] = React.useState<boolean>(false);
  const [newFolderName, setNewFolderName] = React.useState<string>("");

  const [showOverride, setShowOverride] = React.useState<boolean>(false);
  const [targetPath, setTargetPath] = React.useState<string>("");

  React.useEffect(() => {
    ls("");
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  //ls directory
  const ls = (path: string) => {
    const cb = (resp: any) => {
      //check resp code
      if (resp.status !== 200) {
        //do something here
        console.log("Error from server: ", resp.payload);
        return;
      }
      //update
      console.log("explorer/list received response");
      var data = JSON.parse(resp.payload);
      console.log(data);

      //expecting data to be array of IFile
      dispatch({ type: "files", payload: data ? data : [] });
      dispatch({ type: "loading", payload: false });
    };
    dispatch({ type: "loading", payload: true });
    storeDispatch(sendMessage("file", "ls", path, cb));
  };

  //next
  const handleCD = (next: string) => {
    return () => {
      setSelected("");
      dispatch({ type: "path", payload: next });
      ls(next);
    };
  };

  const handleEditNewFileName = (e: React.ChangeEvent<HTMLInputElement>) => {
    setNewFileName(e.target.value);
  };

  const openFile = (target: string) => {
    const cb = (resp: any) => {
      //check resp code
      if (resp.status !== 200) {
        //do something here
        console.log("Error from server: ", resp.payload);
        return;
      }
      //update
      console.log("explorer/open received response");
      var data = JSON.parse(resp.payload);
      console.log(data);
      storeDispatch(setConfig(data.data));
      storeDispatch(setCurrentPath(target));
      storeDispatch(setHasChange(false));
    };
    storeDispatch(sendMessage("file", "open/file", target, cb));
  };

  const newFile = () => {
    if (newFileName === "") {
      return;
    }
    const cb = (resp: any) => {
      setShowNewFile(false);
      //check resp code
      if (resp.status !== 200) {
        //do something here
        console.log("Error from server: ", resp.payload);
        return;
      }
      //update
      ls(state.currentPath);
    };
    storeDispatch(
      sendMessage("file", "new/file", state.currentPath + "/" + newFileName, cb)
    );
  };

  const newFolder = () => {
    if (newFolderName === "") {
      return;
    }
    const cb = (resp: any) => {
      setShowNewFolder(false);
      //check resp code
      if (resp.status !== 200) {
        //do something here
        console.log("Error from server: ", resp.payload);
        return;
      }
      //update
      ls(state.currentPath);
    };
    storeDispatch(
      sendMessage(
        "file",
        "new/folder",
        state.currentPath + "/" + newFolderName,
        cb
      )
    );
  };

  const handleNodeClick = (n: TreeNodeInfo) => {
    console.log(n);
    let path = n.id.toString();
    if (n.nodeData) {
      setSelected("");
      dispatch({ type: "path", payload: path });
      ls(path);
    } else {
      setSelected(path);
      if (!hasChange) {
        openFile(path);
      } else {
        setTargetPath(path);
        setShowOverride(true);
      }
    }
  };

  //build current path
  var paths = state.currentPath.split("/");

  var crumbs: BreadcrumbProps[] = [];
  crumbs.push({
    icon: "folder-close",
    text: "home",
    onClick: handleCD(""),
  });

  var pa = "";

  for (var i = 0; i < paths.length; i++) {
    if (paths[i] === "") continue;
    if (i !== 0) {
      pa += "/";
    }
    pa += paths[i];
    crumbs.push({
      icon: "folder-close",
      text: paths[i],
      onClick: handleCD(pa),
    });
  }

  let rows: TreeNodeInfo[] = state.files.map((e, i) => {
    return {
      id: e.path,
      hasCaret: false,
      label: e.name,
      icon: e.is_dir ? "folder-close" : "document",
      isSelected: selected === e.path,
      nodeData: e.is_dir,
    };
  });

  return (
    <div className="box">
      <Card>
        <Breadcrumbs items={crumbs}></Breadcrumbs>

        <Tree
          contents={rows}
          onNodeClick={handleNodeClick}
          className={Classes.ELEVATION_0}
        />

        <br />
      </Card>
      <ButtonGroup fill vertical>
        <Button
          icon="add"
          intent="primary"
          onClick={() => setShowNewFolder(true)}
        >
          New Folder
        </Button>
        <Button
          icon="add"
          intent="primary"
          onClick={() => setShowNewFile(true)}
        >
          New File
        </Button>
        <Button icon="delete" intent="danger" disabled>
          Delete
        </Button>
      </ButtonGroup>
      <Dialog
        isOpen={showOverride}
        onClose={() => setShowOverride(false)}
        title="New File"
      >
        <div>
          <div className={Classes.DIALOG_BODY}>
            <Callout intent="warning">
              You have unsaved changes. If you choose to continue you will lose
              your change.
            </Callout>
          </div>
          <div className={Classes.DIALOG_FOOTER}>
            <div className={Classes.DIALOG_FOOTER_ACTIONS}>
              <Button onClick={() => setShowOverride(false)}>Cancel</Button>
              <Button
                intent={Intent.WARNING}
                onClick={() => {
                  openFile(targetPath);
                  setShowOverride(false);
                }}
              >
                Continue
              </Button>
            </div>
          </div>
        </div>
      </Dialog>
      <Dialog
        isOpen={showNewFile}
        onClose={() => setShowNewFile(false)}
        title="New File"
      >
        <div className={Classes.DIALOG_BODY}>
          <InputGroup
            value={newFileName}
            placeholder="enter file name"
            onChange={handleEditNewFileName}
          ></InputGroup>
          <ButtonGroup fill>
            <Button
              intent="primary"
              disabled={newFileName === ""}
              onClick={() => {
                newFile();
              }}
            >
              Create
            </Button>
            <Button onClick={() => setShowNewFile(false)}>Cancel</Button>
          </ButtonGroup>
        </div>
      </Dialog>
      <Dialog
        isOpen={showNewFolder}
        onClose={() => setShowNewFolder(false)}
        title="New Folder"
      >
        <div className={Classes.DIALOG_BODY}>
          <InputGroup
            value={newFolderName}
            placeholder="enter folder name"
            onChange={(e) => setNewFolderName(e.target.value)}
          ></InputGroup>
          <ButtonGroup fill>
            <Button
              intent="primary"
              disabled={newFolderName === ""}
              onClick={() => {
                newFolder();
              }}
            >
              Create
            </Button>
            <Button onClick={() => setShowNewFolder(false)}>Cancel</Button>
          </ButtonGroup>
        </div>
      </Dialog>
    </div>
  );
}

export default Explorer;
