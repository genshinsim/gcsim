import { SimResults } from "./DataType";
import { KeyboardEvent, ClipboardEvent, MouseEvent} from "react";

export function Config({ data }: { data: SimResults }) {
    function onKeyDown(e: KeyboardEvent) {
        if (e.metaKey || e.ctrlKey) {
            return true;
        }
        e.preventDefault()
        return false;
    }
    function preventDefault(e: ClipboardEvent) {
        e.preventDefault();
        return false;
    }

    function copyToClipboard(e: MouseEvent) {
        navigator.clipboard.writeText(data.config_file)
        // TODO: Need to add a blueprintjs Toaster for ephemeral confirmation box
    }
    return (
    <div>
        <button className="m-2 p-2 rounded-md bg-gray-600" onClick={ copyToClipboard }>Copy Config to Clipboard
        </button>
        <div className="m-2 p-2 rounded-md bg-gray-600" contentEditable="true" onKeyDown={ onKeyDown } onCut={ preventDefault }>
          <pre className="whitespace-pre-wrap">
              {data.config_file}
          </pre>
      </div>
    </div>
  );
}
