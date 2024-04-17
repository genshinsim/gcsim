import { Spinner } from "@blueprintjs/core";
import { Button, DBCard, Toaster, useToast } from "@gcsim/components";
import { db } from "@gcsim/types";

export function ListView({ data }: { data: db.Entry[] }) {
  const { toast } = useToast();
  if (!data) {
    return (
      <div>
        <Spinner />
      </div>
    );
  }

  const copyConfig = (cfg: string) => {
    if (cfg === "") {
      toast({
        title: "Failed",
        description: "Copied failed unexpected, no config found.",
      });
    }
    navigator.clipboard.writeText(cfg).then(() => {
      console.log("copy ok");
      toast({
        title: "Copied to clipboard",
        description: `Copied config to clipboard`,
      });
    });
  };

  return (
    <>
      <div className="flex flex-col gap-2 justify-center align-middle items-center ">
        {data.map((entry, index) => {
          return (
            <DBCard
              entry={entry}
              key={index}
              className="min-[1300px]:w-[970px] border-0 w-full max-w-[970px]"
              footer={
                <div className="flex flex-row flex-wrap place-content-end mr-2 gap-4">
                  <Button
                    className="bg-emerald-600"
                    onClick={() => {
                      copyConfig(entry.config ?? "");
                    }}
                  >
                    Copy Config
                  </Button>
                  <a
                    href={"https://gcsim.app/db/" + entry._id}
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    <Button className="bg-blue-600">Open in Viewer</Button>
                  </a>
                </div>
              }
            />
          );
        })}
        <Toaster />
      </div>
      {/* <div className="flex flex-col gap-2">
        {data.map((entry, index) => {
          return <DBEntryView dbEntry={entry} key={index} />;
        })}
      </div> */}
    </>
  );
}
