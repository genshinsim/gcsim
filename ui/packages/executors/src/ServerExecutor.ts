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
        .catch(function (error) {
          reject(error.message);
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
          console.log(resp);
          resolve(resp.data);
        })
        .catch(function (error) {
          reject(error.message);
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
            const simres = JSON.parse(resp.data.result);
            updateResult(simres, resp.data.hash);
            //if not done make another update request
            if (resp.data.done) {
              c.is_running = false;
              if (resp.data.error !== "") {
                reject(resp.data.error);
              } else {
                resolve(true);
              }
            } else {
              setTimeout(() => {
                update();
              }, 100);
            }
          })
          .catch(function (resp) {
            //this should be either 404 or 500 if something went wrong
            console.log("something went wrong fetch updated results", resp);
            reject(resp.message);
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
          console.log(error);
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
