import { TextArea } from "@blueprintjs/core";
import Editor from "react-simple-code-editor";

//@ts-ignore
import { highlight, languages } from "prismjs/components/prism-core";
import "prismjs/components/prism-gcsim";

// import Prism from "prismjs";
import "prismjs/themes/prism-tomorrow.css";

// Prism.highlight("stuff",Prism.languages)

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
      <Editor
        value={props.cfg}
        onValueChange={(code) => props.onChange(code)}
        highlight={(code) => highlight(code, languages.gcsim)}
        insertSpaces
        padding={10}
        style={{
          fontFamily: '"Fira code", "Fira Mono", monospace',
          fontSize: 14,
          backgroundColor: "rgb(75 85 99)",
        }}
      />
    </div>
  );
}
