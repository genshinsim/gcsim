import { AuthContext } from "../SharedComponents/Management.context";
import { Database } from "./Database";

export default function Management() {
  return (
    <AuthContext.Provider value={{ isAdmin: true }}>
      <Database />
    </AuthContext.Provider>
  );
}
