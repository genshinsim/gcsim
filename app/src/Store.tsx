import React from "react";

type Action = { type: "set"; str: string };

function reducer(state: string, action: Action): string {
  switch (action.type) {
    case "set":
      return action.str;
    default:
      return "";
  }
}

export const AppContext = React.createContext<{
  state: string;
  dispatch: React.Dispatch<Action>;
}>({ state: "", dispatch: () => null });

const Store: React.FC = ({ children }) => {
  const [state, dispatch] = React.useReducer(reducer, "");

  return (
    <AppContext.Provider value={{ state, dispatch }}>
      {children}
    </AppContext.Provider>
  );
};

export default Store;
