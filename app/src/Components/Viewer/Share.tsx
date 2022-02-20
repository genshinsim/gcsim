import {
  Dialog,
  Classes,
  Button,
  Spinner,
  SpinnerSize,
  ButtonGroup,
  Callout,
  Checkbox,
  FormGroup,
  InputGroup,
  Position,
  Toaster,
} from "@blueprintjs/core";
import axios, { AxiosRequestHeaders } from "axios";
import pako from "pako";
import React from "react";
import { bytesToBase64 } from "./base64";
import { SimResults } from "./DataType";

export interface ShareProps {
  // isOpen: boolean;
  // handleClose: () => void;
  data: SimResults;
}

const disabled = false;

const ak = "api-key";

export const AppToaster = Toaster.create({
  position: Position.BOTTOM_RIGHT,
});

export default function Share(props: ShareProps) {
  const [loading, setIsLoading] = React.useState<boolean>(false);
  const [errMsg, setErrMsg] = React.useState<string>("");
  const [url, setURL] = React.useState<string>("");
  const [isPerm, setIsPerm] = React.useState<boolean>(false);
  const [perm, setPerm] = React.useState<boolean>(false);
  const [apiKey, setAPIKey] = React.useState<string>("");
  const [viewPass, setViewPass] = React.useState<boolean>(false);

  React.useEffect(() => {
    let key = localStorage.getItem(ak);
    if (key !== null && key !== "") {
      setAPIKey(key);
    }
  }, []);

  const handleUpload = () => {
    //encode data
    let compressed = pako.deflate(JSON.stringify(props.data));

    // const restored = JSON.parse(pako.inflate(compressed, { to: "string" }));

    // console.log(restored);

    let s = bytesToBase64(compressed);

    // console.log(s);
    //"{\"author\":\"anon\",\"description\":\"none\",\"data\":\"stuff\"}"
    setIsLoading(true);
    setIsPerm(false);
    setURL("");
    let h: AxiosRequestHeaders = {
      "Access-Control-Allow-Origin": "*",
    };

    if (perm) {
      h[ak] = apiKey;
    }

    axios({
      method: "post",
      url: "https://viewer.gcsim.workers.dev/upload",
      headers: h,
      data: {
        author: "anon",
        description: "none",
        data: s,
      },
    })
      .then((response) => {
        console.log(response);
        if (response.data.id) {
          setErrMsg("");
          setURL(response.data.id);
          setIsPerm(response.data.perm);
          setIsLoading(false);
        } else {
          setErrMsg("upload failed");
          setURL("");
          setIsLoading(false);
        }
      })
      .catch((error) => {
        console.log(error);
        setErrMsg("error encountered: " + error);
        setURL("");
        setIsLoading(false);
      });
  };

  const handleCopy = () => {
    navigator.clipboard.writeText(`https://gcsim.app/viewer/share/${url}`).then(
      () => {
        AppToaster.show({ message: "Copied to clipboard", intent: "success" });
      },
      () => {
        AppToaster.show({
          message: "Error copying :( Not sure what went wrong",
          intent: "danger",
        });
      }
    );
  };

  return (
    <div className="wide:w-[70rem] ml-auto mr-auto bg-gray-600 rounded-md p-4 flex flex-col gap-2">
      <div>
        <div className="font-bold text-lg mb-2">Create a shareable link</div>
        <div>
          Note that by default shareable links are <b>only valid for 7 days</b>.
          This is done in order to keep server storage usage at a reasonable
          level. Contributors and Ko-Fi supporters can enter in a private key to
          make their links permanent.
        </div>
      </div>
      <div className="flex flex-col place-items-center">
        <FormGroup label="Make link permanent?" inline>
          <Checkbox checked={perm} onClick={() => setPerm(!perm)} />
        </FormGroup>
        <FormGroup label="API Key (Supporters only)" inline>
          <InputGroup
            type={viewPass ? "text" : "password"}
            value={apiKey}
            onChange={(v) => {
              const val = v.target.value;
              setAPIKey(val);
              localStorage.setItem(ak, val);
            }}
            rightElement={
              <Button
                icon={viewPass ? "unlock" : "lock"}
                intent="warning"
                onClick={() => setViewPass(!viewPass)}
              />
            }
          />
        </FormGroup>
      </div>
      <div>
        Please note that this is not an attempt to get people to donate money.
        This is a simple way to gate the amount of data being shared while
        providing a small thank you for those that either supported this project
        financially or by contributing. If you do need a permanent link and you
        don't have the access key, simply hop over to our Discord and ask
        someone to do it for you.
      </div>
      <ButtonGroup fill className="mb-4">
        <Button
          intent="primary"
          onClick={handleUpload}
          disabled={loading || disabled}
        >
          Upload
        </Button>
      </ButtonGroup>
      {loading ? <Spinner size={SpinnerSize.LARGE} /> : null}
      {errMsg === "" ? (
        url !== "" ? (
          <Callout intent="success">
            <div className="flex flex-col gap-2 place-items-center">
              <span className="text-lg">
                Upload successful.{" "}
                {isPerm
                  ? "Link is permanent. "
                  : "Link will expire in 7 days. "}{" "}
                Results can be viewed at:
              </span>
              <div className="p-2 rounded-md bg-green-700">
                <pre>{`https://gcsim.app/viewer/share/${url}`}</pre>
              </div>
              <Button
                intent="success"
                onClick={handleCopy}
                disabled={url === ""}
              >
                Copy Link
              </Button>
            </div>
          </Callout>
        ) : null
      ) : (
        <Callout intent="warning">{errMsg}</Callout>
      )}
    </div>
  );
}
