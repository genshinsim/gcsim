import {
  Button,
  DBCard,
  Separator,
  Toaster,
  useToast,
} from "@gcsim/components";
import "@gcsim/components/src/index.css";
import { Entry } from "@gcsim/types/src/generated/index.db";
import axios from "axios";
import React from "react";
import { Route, Switch } from "wouter";

function App({ id }: { id: string }) {
  const [main, setMain] = React.useState<Entry | null>(null);
  const [data, setData] = React.useState<Entry[]>([]);
  const { toast } = useToast();

  React.useEffect(() => {
    axios.get(`/api/db/id/${id}`).then((res) => {
      console.log(res);
      if (res.data) {
        setMain(res.data);
      }
    });
  }, [id, setData]);

  React.useEffect(() => {
    if (main === null) {
      return;
    }
    const includedChars = main.summary?.char_names;
    if (includedChars === null || includedChars == undefined) {
      return;
    }
    let q = {
      query: {},
      limit: 100,
      sort: {
        created_data: -1,
      },
    };
    if (includedChars.length > 0) {
      const and: unknown[] = [];
      const trav: { [key in string]: boolean } = {};
      includedChars.forEach((char) => {
        if (char.includes("aether") || char.includes("lumine")) {
          const ele = char.replace(/(aether|lumine)(.+)/, "$2");
          trav[ele] = true;
          return;
        }
        and.push({
          "summary.char_names": char,
        });
      });
      Object.keys(trav).forEach((ele) => {
        and.push({
          $or: [
            { "summary.char_names": `aether${ele}` },
            { "summary.char_names": `lumine${ele}` },
          ],
        });
      });
      if (and.length > 0) {
        q.query["$and"] = and;
      }
    }
    axios
      .get(`/api/db?q=${encodeURIComponent(JSON.stringify(q))}`)
      .then((res) => {
        console.log(res);
        if (res.data && res.data.data && res.data.data.length > 0) {
          setData(res.data.data);
        }
      });
  }, [main, setData]);

  if (main === null) {
    return <div className="text-gray-200">Loading, please wait...</div>;
  }

  const copy = (cmd: string) => {
    const s = `/${cmd} id:${id}`;
    navigator.clipboard.writeText(s).then(() => {
      console.log("copy ok");
      toast({
        title: "Copied to clipboard",
        description: `Copied ${s} to clipboard`,
      });
    });
  };

  const copyReplace = (from: string) => {
    const s = `/replace id:${from} link:${main.share_key}`;
    console.log("copying command: ", s);
    navigator.clipboard.writeText(s).then(() => {
      console.log("copy ok");
      toast({
        title: "Copied to clipboard",
        description: `Copied replace command ${s} to clipboard`,
      });
    });
  };

  const rows = data
    .filter((e) => e["_id"] !== id)
    .map((e, i) => {
      return (
        <DBCard
          className="border-0"
          key={"entry-" + i}
          entry={e}
          skipTags={-1}
          footer={
            <div className="flex flex-row flex-wrap place-content-end mr-2 gap-4">
              <Button
                className="bg-yellow-600"
                onClick={() => {
                  copyReplace(e._id ?? "");
                }}
              >
                Replace This
              </Button>
              <a
                href={"https://gcsim.app/db/" + e._id}
                target="_blank"
                rel="noopener noreferrer"
              >
                <Button className="bg-blue-600">Result Viewer</Button>
              </a>
            </div>
          }
        />
      );
    });

  return (
    <div className="flex flex-col place-items-center m-4">
      <div className="m-2 text-white p-2 w-full flex flex-col place-items-center">
        <span className="text-lg font-semibold">
          Showing entries with the same team for id: {id}
        </span>
        {main !== null ? (
          <DBCard
            className="bg-stone-800 border-0"
            entry={main}
            skipTags={-1}
            footer={
              <div className="flex flex-row flex-wrap place-content-end mr-2 gap-4">
                <Button className=" bg-red-600" onClick={() => copy("reject")}>
                  Copy Reject
                </Button>
                <Button
                  className="bg-green-600"
                  onClick={() => copy("approve")}
                >
                  Copy Approve
                </Button>
                <a
                  href={"https://gcsim.app/db/" + id}
                  className=""
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  <Button className="bg-blue-600">Result Viewer</Button>
                </a>
              </div>
            }
          />
        ) : null}
      </div>
      <Separator />
      <div className="m-2 p-2  text-white w-full flex flex-col place-items-center">
        <span className="text-lg font-semibold">
          Existing sims with same characters
        </span>
        {rows.length > 0 ? rows : <div>Nothing found</div>}
      </div>
      <Toaster />
    </div>
  );
}

const Routes = (): JSX.Element => {
  return (
    <>
      <Switch>
        <Route path="/">
          <div>nothing here</div>
        </Route>
        <Route path="/id/:id">{({ id }) => <App id={id} />}</Route>
      </Switch>
    </>
  );
};

export default Routes;
