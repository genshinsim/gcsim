import type { Sim } from "@gcsim/types";
import type { Executor } from "./executor.js";

export class ExecutorError extends Error {
  constructor(
    message: string,
    public readonly code: "NETWORK" | "SERVER" | "PARSE" | "UNKNOWN",
    public readonly cause?: unknown,
  ) {
    super(message);
    this.name = "ExecutorError";
  }
}

export class ServerExecutor implements Executor {
  private ipaddr: string;
  private id: string;
  private isRunning: boolean;
  private readyCache: boolean | undefined;

  constructor(ipaddr: string) {
    this.ipaddr = ipaddr;
    this.id = `id${Date.now()}`;
    this.isRunning = false;
  }

  public setUrl(ipaddr: string): void {
    this.ipaddr = ipaddr;
    this.readyCache = undefined;
  }

  public async ready(): Promise<boolean> {
    if (this.readyCache !== undefined) {
      return this.readyCache;
    }

    try {
      const resp = await fetch(`${this.ipaddr}/ready/${this.id}`);
      this.readyCache = resp.ok;
      return resp.ok;
    } catch {
      this.readyCache = false;
      return false;
    }
  }

  public running(): boolean {
    return this.isRunning;
  }

  public async validate(cfg: string): Promise<Sim.ParsedResult> {
    let resp: Response;
    try {
      resp = await fetch(`${this.ipaddr}/validate/${this.id}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ config: cfg }),
      });
    } catch (err) {
      throw new ExecutorError(
        "Network error encountered communicating with server",
        "NETWORK",
        err,
      );
    }

    const data = await resp.json();
    if (typeof data === "string") {
      throw new ExecutorError(data, "SERVER");
    }

    return {
      characters: data.characters,
      errors: data.error_msgs,
      player_initial_pos: data.initial_player_pos,
    };
  }

  public async sample(cfg: string, seed: string): Promise<Sim.Sample> {
    let resp: Response;
    try {
      resp = await fetch(`${this.ipaddr}/sample/${this.id}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ config: cfg, seed: parseInt(seed, 10) }),
      });
    } catch (err) {
      throw new ExecutorError(
        "Network error encountered communicating with server",
        "NETWORK",
        err,
      );
    }

    return resp.json();
  }

  public async run(
    cfg: string,
    updateResult: (result: Sim.SimResults, hash: string) => void,
  ): Promise<boolean | void> {
    try {
      await fetch(`${this.ipaddr}/run/${this.id}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ config: cfg }),
      });
    } catch (err) {
      this.isRunning = false;
      throw new ExecutorError(
        "Network error encountered communicating with server",
        "NETWORK",
        err,
      );
    }

    this.isRunning = true;
    return this.pollResults(updateResult);
  }

  private async pollResults(
    updateResult: (result: Sim.SimResults, hash: string) => void,
  ): Promise<boolean> {
    while (this.isRunning) {
      let resp: Response;
      try {
        resp = await fetch(`${this.ipaddr}/results/${this.id}`);
      } catch (err) {
        this.isRunning = false;
        throw new ExecutorError(
          "Network error encountered communicating with server",
          "NETWORK",
          err,
        );
      }

      const data = await resp.json();

      if (data.error !== "") {
        this.isRunning = false;
        throw new ExecutorError(data.error, "SERVER");
      }

      if (data.result === "") {
        this.isRunning = false;
        throw new ExecutorError("Unexpected response from server: blank result", "SERVER");
      }

      let simres: Sim.SimResults;
      try {
        simres = JSON.parse(data.result);
      } catch (err) {
        this.isRunning = false;
        throw new ExecutorError(`Could not parse sim result: ${err}`, "PARSE", err);
      }

      updateResult(simres, data.hash);

      if (data.done) {
        this.isRunning = false;
        return true;
      }

      await new Promise((resolve) => setTimeout(resolve, 100));
    }

    return true;
  }

  private async sendCancel(): Promise<void> {
    try {
      await fetch(`${this.ipaddr}/cancel/${this.id}`, { method: "POST" });
    } catch {
      // Ignore cancel errors
    }
  }

  public cancel(): void {
    this.sendCancel();
    this.isRunning = false;
  }

  public buildInfo(): { hash: string; date: string } {
    return { hash: "", date: "" };
  }
}
