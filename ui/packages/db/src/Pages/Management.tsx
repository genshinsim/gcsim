import { createContext } from "react";
import { Database } from "./Database";

export const AuthContext = createContext({
  isAdmin: false,
});
export default function Management() {
  return (
    <AuthContext.Provider value={{ isAdmin: true }}>
      <Database />
    </AuthContext.Provider>
  );
}
