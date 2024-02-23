import { ParsedResult, Sample, SimResults } from "@gcsim/types";
import axios from "axios";
import { Executor } from "./Executor";

export class ServerExecutor implements Executor {
  private ipaddr: string;
  private id: string; //unique id for this instance
  private is_running: boolean;
  private ready_cache: boolean | undefined;

  constructor(ipaddr: string) {
    this.ipaddr = ipaddr;
    this.id = "id" + new Date().getTime();
    this.is_running = false;
  }

  public set_url(ipaddr: string) {
    this.ipaddr = ipaddr;
    console.log("updating url ", ipaddr);
    this.ready_cache = undefined;
  }

  public ready(): Promise<boolean> {
    if (this.ready_cache != undefined) {
      const ready = this.ready_cache;
      return new Promise((resolve) => resolve(ready));
    }

    this.ready_cache = undefined;
    const c = this;
    return new Promise((resolve) => {
      axios
        .get(`${this.ipaddr}/ready/${this.id}`)
        .then(function (resp) {
          c.ready_cache = resp.status == 200;
          resolve(resp.status == 200);
        })
        .catch(function (error) {
          c.ready_cache = false;
          resolve(false);
        });
    });
  }

  public running(): boolean {
    return this.is_running;
  }

  public validate(cfg: string): Promise<ParsedResult> {
    return new Promise((resolve, reject) => {
      axios
        .post(`${this.ipaddr}/validate/${this.id}`, {
          config: cfg,
        })
        .then(function (resp) {
          //resp should be json body?
          console.log(resp);
          if (typeof resp.data == "string") {
            resolve({
              characters: [],
              errors: [resp.data],
              player_initial_pos: {
                x: 0,
                y: 0,
                r: 0,
              },
            });
          } else {
            resolve({
              characters: resp.data.characters,
              errors: resp.data.error_msgs,
              player_initial_pos: resp.data.initial_player_pos,
            });
          }
        })
        .catch(function (resp) {
          console.log("something went wrong validating", resp);
          if (resp.code === "ERR_NETWORK") {
            reject("Network error encountered communicating with server");
          }
          {
            reject(
              "Unknown error encountered communicating with server: " +
                resp.message
            );
          }
        });
    });
  }

  public sample(cfg: string, seed: string): Promise<Sample> {
    const c = this;
    return new Promise((resolve, reject) => {
      axios
        .post(`${this.ipaddr}/sample/${this.id}`, {
          config: cfg,
          seed: parseInt(seed),
        })
        .then(function (resp) {
          //resp should be json body?
          console.log("sample resp", resp);
          resolve(resp.data);
        })
        .catch(function (resp) {
          console.log("something went wrong fetch sample", resp);
          if (resp.code === "ERR_NETWORK") {
            reject("Network error encountered communicating with server");
          }
          {
            reject(
              "Unknown error encountered communicating with server: " +
                resp.message
            );
          }
        });
    });
  }

  public run(
    cfg: string,
    updateResult: (result: SimResults, hash: string) => void
  ): Promise<boolean | void> {
    const c = this;
    return new Promise((resolve, reject) => {
      const update = () => {
        axios
          .get(`${this.ipaddr}/results/${this.id}`)
          .then(function (resp) {
            console.log("result resp", resp);
            //handle error first before attempting to parse since data could be empty
            if (resp.data.error !== "") {
              c.is_running = false;
              reject(resp.data.error);
              return;
            }
            //sanity check just in case result is blank; shouldn't happen though
            if (resp.data.result === "") {
              c.is_running = false;
              reject("unexpected response from server: blank result");
              return;
            }
            let simres: SimResults;
            try {
              simres = JSON.parse(resp.data.result);
              updateResult(simres, resp.data.hash);
            } catch (e) {
              c.is_running = false;
              console.log("error decoding sim result");
              reject("could not unmarshall sim result: " + e);
              return;
            }
            //end sim now
            if (resp.data.done) {
              c.is_running = false;
              resolve(true);
              return;
            }
            //otherwise keep making update polling
            setTimeout(() => {
              update();
            }, 100);
          })
          .catch(function (resp) {
            //this should be either 404 or 500 if something went wrong
            c.is_running = false;
            console.log("something went wrong fetch updated results", resp);
            if (resp.code === "ERR_NETWORK") {
              reject("Network error encountered communicating with server");
            }
            {
              reject(
                "Unknown error encountered communicating with server: " +
                  resp.message
              );
            }
          });
      };
      axios
        .post(`${this.ipaddr}/run/${this.id}`, {
          config: cfg,
        })
        .then(function (resp) {
          console.log("run resp", resp);
          c.is_running = true;
          update();
        })
        .catch(function (error) {
          //this should be bad requests
          console.log("error executing run", error);
          reject(error.message);
          c.is_running = false;
        });
    });
  }

  private async send_cancel() {
    const resp = await axios.post(`${this.ipaddr}/cancel/${this.id}`);
  }

  public cancel(): void {
    this.send_cancel();
    this.is_running = false;
  }

  public buildInfo(): { hash: string; date: string } {
    return {
      hash: "",
      date: "",
    };
  }
}
