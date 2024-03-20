import "@fontsource/fira-mono";
import AceEditor from "react-ace";
//THESE IMPORTS NEEDS TO BE AFTER IMPORTING AceEditor
import "ace-builds/src-noconflict/ext-language_tools";
import "ace-builds/src-noconflict/theme-tomorrow_night";
import "../../../util/mode-gcsim.js";

type Props = {
  cfg: string;
  onChange: (v: string) => void;
};

export function ActionList(props: Props) {
  // const t = () => {
  //   Prism.highlightElement();
  // };
  return (
    <div className="p-1 md:p-2">
      <AceEditor
        mode="gcsim"
        theme="tomorrow_night"
        width="100%"
        onChange={props.onChange}
        value={props.cfg}
        name="UNIQUE_ID_OF_DIV"
        editorProps={{
          $blockScrolling: true,
        }}
        setOptions={{
          maxLines: Infinity,
          fontSize: 14,
          tabSize: 2,
        }}
      />
    </div>
  );
}
