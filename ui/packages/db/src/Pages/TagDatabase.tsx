import { initialFilter } from "SharedComponents/FilterComponents/Filter.utils";
import { Database } from "./Database";

export default function TagDatabase({ tag }: { tag: string }) {
  const filter = {
    ...initialFilter,
    tags: [parseInt(tag)],
  };

  return <Database initialFilter={filter} />;
}
