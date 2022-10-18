import Editor from "react-simple-code-editor";
import { NonIdealState, Spinner, SpinnerSize } from "@blueprintjs/core";

//@ts-ignore
import { highlight, languages } from "prismjs/components/prism-core";
import "prismjs/components/prism-gcsim";
import "prismjs/themes/prism-tomorrow.css";


type ConfigProps = {
  cfg: string | undefined;
};

export default ({ cfg }: ConfigProps) => {
  if (cfg === undefined) {
    return <NonIdealState icon={<Spinner size={SpinnerSize.LARGE} />} />;
  }

  return (
    <div className="w-full 2xl:mx-auto 2xl:container">
      <Editor
        value={cfg}
        disabled={true}
        onValueChange={() => {}}
        textareaId="codeArea"
        className="editor"
        highlight={(code) =>
          highlight(code, languages.gcsim)
            .split("\n")
            .map(
              (line: string, i: number) =>
                `<span class='editorLineNumber'>${i + 1}</span>${line}`
            )
            .join("\n")
        }
        insertSpaces
        padding={10}
        style={{
          fontFamily: '"Fira code", "Fira Mono", monospace',
          fontSize: 14,
          backgroundColor: "rgb(45 45 45)",
        }}
      />
    </div>
  );
};