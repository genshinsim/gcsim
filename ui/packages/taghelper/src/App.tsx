import { Button, DBCard, Toaster, useToast } from "@gcsim/components";
import { Card } from "@gcsim/components/src/common/ui/card";
import "@gcsim/components/src/index.css";
import { db } from "@gcsim/types";
import axios from "axios";
import React from "react";
import { Route, Switch } from "wouter";

function App({ id }: { id: string }) {
  const [main, setMain] = React.useState<db.IEntry | null>(null);
  const [data, setData] = React.useState<db.IEntry[]>([]);
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

  const rows = data
    .filter((e) => e["_id"] !== id)
    .map((e, i) => {
      return (
        <DBCard
          key={e.id ?? "entry-" + i}
          entry={e}
          footer={
            <a
              href={"https://gcsim.app/db/" + e["_id"]}
              className=" text-red-600 font-extrabold ml-auto mr-auto"
            >
              Result Viewer
            </a>
          }
        />
      );
    });

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

  return (
    <div className="flex flex-col place-items-center m-4">
      <Card className="m-2 bg-green-900 text-white p-2 w-full flex flex-col place-items-center">
        Showing entries with the same team for id: {id}
        {main !== null ? (
          <DBCard
            entry={main}
            footer={
              <a
                href={"https://gcsim.app/db/" + id}
                className=" text-red-600 font-extrabold ml-auto mr-auto"
              >
                Result Viewer
              </a>
            }
          />
        ) : null}
        <div className="w-full flex flex-row place-content-center">
          <Button className=" bg-red-600" onClick={() => copy("reject")}>
            Copy Reject
          </Button>
          <Button className="ml-4 bg-blue-600" onClick={() => copy("approve")}>
            Copy Approve
          </Button>
        </div>
      </Card>
      <Card className="m-2 p-2 bg-blue-900 text-white w-full flex flex-col place-items-center">
        <span>Existing sims with same characters</span>
        {rows.length > 0 ? rows : <div>Nothing found</div>}
      </Card>
      <Toaster />
    </div>
  );
}

function Routes() {
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
}

export default Routes;
