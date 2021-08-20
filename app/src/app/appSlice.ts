import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { AppThunk } from "./store";

var socket: WebSocket;
let cbm = new Map();

function msgID() {
  // Math.random should be unique because of its seeding algorithm.
  // Convert it to base 36 (numbers + letters), and grab the first 9 characters
  // after the decimal.
  return "_" + Math.random().toString(36).substr(2, 9);
}

export function sendMessage(
  type: string,
  target: string,
  payload: string,
  callback: (resp: any) => void
): AppThunk {
  return function () {
    //generate an id, try till ok
    var ok = false;
    var id = "";
    //this code should only run once but is here just in case of collision
    while (!ok) {
      id = msgID();
      //check if id already exists
      ok = !cbm.has(id);
      if (ok) {
        cbm.set(id, callback);
      }
    }
    var packet = {
      id: id,
      token: "",
      type: type,
      target: target,
      payload: payload,
    };
    console.log("socket - sending msg: ", packet);
    socket.send(JSON.stringify(packet));
    console.log("callbacks left: ", cbm);
    //register a call back for id?
  };
}

export function openSocket(): AppThunk {
  return function (dispatch, getState) {
    console.log("attempting to open socket");
    //we can call dispatch here to trigger changes on receive
    //attempt to open a socket
    // var loc = window.location,
    //   new_uri;
    // if (loc.protocol === "https:") {
    //   new_uri = "wss:";
    // } else {
    //   new_uri = "ws:";
    // }
    // new_uri += "//" + loc.hostname + "/ws";
    var socketstring = "ws://" + process.env.REACT_APP_HOST_ADDRESS + "/ws";
    console.log(socketstring);
    socket = new WebSocket(socketstring);
    //set up some handlers
    socket.onopen = (e) => {
      console.log("socket opened");
      dispatch(setSocketState(true));
    };
    socket.onerror = (e) => {
      console.log("socket got error");
      console.log(e);
    };
    socket.onclose = (e) => {
      console.log("socket closed");
      dispatch(setSocketState(false));
      dispatch(retryOpen());
    };
    socket.onmessage = (e) => {
      console.log("socket - received message: ", e);
      var response = JSON.parse(e.data);
      //check if id of result in callback
      if (cbm.has(response.id)) {
        console.log("socket - reponse received for request id: ", response.id);
        var cb = cbm.get(response.id);
        cb(response);
        cbm.delete(response.id);
        // console.log("callbacks left: ", cbm);
      }
      //we can dispatch stuff from here too
    };
  };
}

export function retryOpen(): AppThunk {
  return function (dispatch, getState) {
    console.log("retrying in 2 seconds...");
    setTimeout(() => {
      var connected = getState().app.isOpen;
      if (!connected) {
        dispatch(openSocket());
      }
    }, 2000);
  };
}

interface IAppState {
  isOpen: boolean;
  config: string;
  currentPath: string;
  isLoading: boolean;
}

const initialState: IAppState = {
  isOpen: false,
  config: "",
  currentPath: "",
  isLoading: false,
};

export const appSlice = createSlice({
  name: "app",
  initialState,
  reducers: {
    setSocketState: (state, action: PayloadAction<boolean>) => {
      state.isOpen = action.payload;
    },
    setConfig: (state, action: PayloadAction<string>) => {
      state.config = action.payload;
    },
    setCurrentPath: (state, action: PayloadAction<string>) => {
      state.currentPath = action.payload;
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.isLoading = action.payload;
    },
  },
});

export const { setSocketState, setConfig, setCurrentPath, setLoading } =
  appSlice.actions;
export default appSlice.reducer;
