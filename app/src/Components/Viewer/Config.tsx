import { SimResults } from "./DataType";

export function Config({ data }: { data: SimResults }) {
  return (
    <div className="flex flex-col">
      <div className="m-2 p-2 rounded-md bg-gray-600">
        <pre className="whitespace-pre-wrap">{data.config_file}</pre>
      </div>
    </div>
  );
}
