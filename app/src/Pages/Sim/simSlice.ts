import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { Character, maxStatLength, Talent, Weapon } from "~/src/types";
import { characterKeyToICharacter } from "~src/Components/Character";
import { AppThunk } from "~src/store";
import { ascLvlMin, maxLvlToAsc } from "~src/util";
import { WorkerPool } from "~src/WorkerPool";

let pool: WorkerPool = new WorkerPool();

type RunStats = {
  progress: number;
  result: number;
  time: number;
};

export interface Sim {
  team: Character[];
  edit_index: number;
  ready: number;
  workers: number;
  cfg: string;
  run: RunStats;
}

const initialState: Sim = {
  team: [],
  edit_index: -1,
  ready: 0,
  workers: 8,
  cfg: `options debug=true iteration=100 duration=90 workers=24;
bennett char lvl=70/80 cons=2 talent=6,8,8; 
bennett add weapon="favoniussword" refine=1 lvl=90/90;
bennett add set="noblesseoblige" count=4;
bennett add stats hp=4780 atk=311 er=0.518 pyro%=0.466 cr=0.311 ; #main
bennett add stats hp=926 hp%=0.21 atk=121 atk%=0.47800000000000004 def=60 em=42 er=0.052000000000000005 cr=0.186 cd=0.327 ; #subs

raidenshogun char lvl=90/90 cons=1 talent=10,10,10; 
raidenshogun add weapon="engulfinglightning" refine=1 lvl=90/90;
raidenshogun add set="emblemofseveredfate" count=4;
raidenshogun add stats hp=4780 atk=311 er=0.518 electro%=0.466 cr=0.311 ; #main
raidenshogun add stats hp=299 hp%=0.053 atk=101 atk%=0.192 def%=0.073 em=42 er=0.148 cr=0.261 cd=1.119 ; #subs

xiangling char lvl=80/90 cons=6 talent=6,9,10; 
xiangling add weapon="staffofhoma" refine=1 lvl=90/90;
xiangling add set="crimsonwitchofflames" count=2;
xiangling add set="gladiatorsfinale" count=2;
xiangling add stats hp=4780 atk=311 er=0.518 pyro%=0.466 cr=0.311 ; #main
xiangling add stats hp=478 hp%=0.047 atk=65 atk%=0.152 def=76 def%=0.051 em=63 er=0.16199999999999998 cr=0.264 cd=0.9960000000000001 ; #subs

xingqiu char lvl=80/90 cons=6 talent=1,9,10; 
xingqiu add weapon="sacrificialsword" refine=5 lvl=90/90;
xingqiu add set="noblesseoblige" count=2;
xingqiu add set="gladiatorsfinale" count=2;
xingqiu add stats hp=4780 atk=311 atk%=0.466 hydro%=0.466 cr=0.311 ; #main
xingqiu add stats hp=299 hp%=0.08199999999999999 atk=78 atk%=0.449 def=63 def%=0.073 em=94 er=0.065 cr=0.15899999999999997 cd=0.831 ; #subs


##Default Enemy
target lvl=100 resist=.1;
# target lvl=100 resist=.1;

##Actions List
active raidenshogun;

# HP particle simulation. Per srl:
# it adds 1 particle randomly, uniformly distributed between 200 to 300 frames after the last time an energy drops
# so in the case above, it adds on avg one particle every 250 frames in effect
# so over 90s of combat that's 90 * 60 / 250 = 21.6 on avg
energy every interval=200,300 amount=1;

raidenshogun attack,attack,attack,attack,dash,attack,attack,attack,attack,dash,attack,attack,attack,attack,dash,attack,attack,charge  +if=.status.raidenburst>0;

# Additional check to reset at the start of the next rotation
raidenshogun skill  +if=.status.xianglingburst==0&&.energy.xingqiu>70&&.energy.xiangling>70;
raidenshogun skill  +if=.status.raidenskill==0;

# Skill is required before burst to activate Kageuchi. Otherwise ER is barely not enough
# For rotations #2 and beyond, need to ensure that Guoba is ready to go. Guoba timing is about 300 frames after XQ fires his skill
xingqiu skill[orbital=1],burst[orbital=1],attack  +if=.cd.xiangling.skill<300;

# Bennett burst goes after XQ burst for uptime alignment. Attack to proc swords
bennett burst,attack,skill  +if=.status.xqburst>0&&.cd.xiangling.burst<180;

# Only ever want to XL burst in Bennett buff and after XQ burst for uptime alignment
xiangling burst,attack,skill,attack,attack  +if=.status.xqburst>0&&.status.btburst>0;
# Second set of actions needed in case Guoba CD comes off while pyronado is spinning
xiangling burst,attack  +if=.status.xqburst>0&&.status.btburst>0;
xiangling skill ;

# Raiden must burst after all others. Requires an attack to allow Bennett buff to apply
raidenshogun burst  +if=.status.xqburst>0&&.status.xianglingburst>0&&.status.btburst>0;

# Funnelling
bennett attack,skill  +if=.status.xqburst>0&&.energy.xiangling<70 +swap_to=xiangling;
bennett skill  +if=.energy.xiangling<70 +swap_to=xiangling;
bennett skill  +if=.energy.xingqiu<80 +swap_to=xingqiu;
bennett attack,skill  +if=.status.xqburst>0 +if=.energy.raidenshogun<90 +swap_to=raidenshogun;

xingqiu attack  +if=.status.xqburst>0;
xiangling attack  +is_onfield;
bennett attack  +is_onfield;
xingqiu attack  +is_onfield;
raidenshogun attack  +is_onfield;

  `,
  run: {
    progress: -1,
    result: -1,
    time: -1,
  },
};

const defWep: { [key in string]: string } = {
  bow: "dullblade",
  catalyst: "dullblade",
  claymore: "dullblade",
  sword: "dullblade",
  polearm: "dullblade",
};

const newChar = (name: string): Character => {
  const c = characterKeyToICharacter[name];
  //default weapons
  return {
    name: name,
    level: 80,
    max_level: 90,
    element: c.element,
    cons: 0,
    weapon: {
      name: defWep[c.weapon_type],
      refine: 1,
      level: 1,
      max_level: 20,
    },
    talents: {
      attack: 6,
      skill: 6,
      burst: 6,
    },
    stats: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    snapshot: [
      0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
    ],
    sets: {},
  };
};

export function loadWorkers(): AppThunk {
  return function (dispatch, getState) {
    //do nothing if ready
    const state = getState();
    if (!state.sim.ready) {
      pool.load(24, (count: number) => {
        //call back for ready
        dispatch(simActions.setWorkerReady(count));
      });
      pool.setMaxWorker(state.sim.workers);
    }
  };
}

const regex = /iteration=(\d+)/;

export function runSim(): AppThunk {
  return function (dispatch, getState) {
    let cfg = getState().sim.cfg;
    //do processing here
    // console.log(cfg);
    // console.log(pool);

    //extract the number of iterations from the config file
    let iters = 1;
    let m = regex.exec(cfg);
    console.log(m);

    if (m) {
      iters = parseInt(m[1]);

      if (isNaN(iters)) {
        console.warn("no iteration found in settings: ", m);
        iters = 1000;
      }
    } else {
      console.log(cfg);
      console.log("regex failed");
    }
    console.log("parsed iters: " + iters);

    //run the sim
    dispatch(
      simActions.setRunStats({
        progress: 0,
        result: -1,
        time: -1,
      })
    );
    let queued = 0;
    let done = 0;
    let avg = 0;
    let progress = 0;
    const startTime = window.performance.now();
    console.time("sim");
    pool.setCfg(cfg);

    const cbFunc = (val: any) => {
      //parse the result
      const res = JSON.parse(val);
      // console.log(
      //   "finish a run, result: " + res.dps,
      //   "done",
      //   done,
      //   "queue",
      //   queued,
      //   "iters",
      //   iters
      // );
      avg += res.dps;

      done++;

      if (done === iters) {
        //done now
        console.log("done running");
        console.timeEnd("sim");
        const end = window.performance.now();
        // setRuntime(end - startTime);
        avg = avg / iters;

        dispatch(
          simActions.setRunStats({
            progress: -1,
            result: avg,
            time: end - startTime,
          })
        );
        return;
      }

      //check progress
      const per = Math.floor((20 * done) / iters);
      if (per > progress) {
        dispatch(
          simActions.setRunStats({
            progress: per,
            result: -1,
            time: -1,
          })
        );
        progress = per;
        // console.log(per);
      }

      //otherwise check if we need to queue more
      if (queued < iters) {
        //queue another worker
        queued++;
        pool.queue({ cmd: "run", cb: cbFunc });
      }
      //otherwise do nothing
    };
    const debugCB = (val: any) => {
      const res = JSON.parse(val);
      if (res.err) {
        console.error(res.err);
        return;
      }
      console.log("finish debug run: ", res);
    };

    pool.queue({ cmd: "debug", cb: debugCB });

    //queue up 20 out of iters
    let count = 20;
    if (count > iters) {
      count = iters;
    }
    for (; queued < count; queued++) {
      pool.queue({ cmd: "run", cb: cbFunc });
    }
  };
}

export const simSlice = createSlice({
  name: "sim",
  initialState: initialState,
  reducers: {
    setRunStats: (state, action: PayloadAction<RunStats>) => {
      state.run = action.payload;
      return state;
    },
    setWorkers: (state, action: PayloadAction<number>) => {
      pool.setMaxWorker(action.payload);
      state.workers = action.payload;
      return state;
    },
    setWorkerReady: (state, action: PayloadAction<number>) => {
      state.ready = action.payload;
      return state;
    },
    setCfg: (state, action: PayloadAction<string>) => {
      state.cfg = action.payload;
      return state;
    },
    setCharacterNameAndEle: (
      state,
      action: PayloadAction<{ name: string; ele: string }>
    ) => {
      state.team[state.edit_index].name = action.payload.name;
      state.team[state.edit_index].element = action.payload.ele;
      return state;
    },
    setCharacterLvl: (state, action: PayloadAction<{ val: number }>) => {
      state.team[state.edit_index].level = action.payload.val;
      return state;
    },
    setCharacterMaxLvl: (state, action: PayloadAction<{ val: number }>) => {
      let m = action.payload.val;
      let l = state.team[state.edit_index].level;
      let asc = maxLvlToAsc(m);
      if (l > m) {
        l = m;
      } else if (l < ascLvlMin(asc)) {
        l = ascLvlMin(asc);
      }

      state.team[state.edit_index].max_level = m;
      state.team[state.edit_index].level = l;
      return state;
    },
    setCharacterCon: (state, action: PayloadAction<{ val: number }>) => {
      state.team[state.edit_index].cons = action.payload.val;
      return state;
    },
    setCharacterSetBonus: (
      state,
      action: PayloadAction<{ set: string; val: number }>
    ) => {
      state.team[state.edit_index].sets[action.payload.set] =
        action.payload.val;
      return state;
    },
    addCharacterSet: (state, action: PayloadAction<{ set: string }>) => {
      state.team[state.edit_index].sets[action.payload.set] = 0;
      return state;
    },
    deleteCharacterSet: (state, action: PayloadAction<{ set: string }>) => {
      delete state.team[state.edit_index].sets[action.payload.set];
      return state;
    },
    setCharacterWeapon: (state, action: PayloadAction<{ val: Weapon }>) => {
      let w = action.payload.val;
      let asc = maxLvlToAsc(w.max_level);
      if (w.level > w.max_level) {
        w.level = w.max_level;
      } else if (w.level < ascLvlMin(asc)) {
        w.level = ascLvlMin(asc);
      }
      state.team[state.edit_index].weapon = w;
      return state;
    },
    setCharacterTalent: (state, action: PayloadAction<{ val: Talent }>) => {
      state.team[state.edit_index].talents = action.payload.val;
      return state;
    },
    setCharacterStats: (
      state,
      action: PayloadAction<{ index: number; val: number }>
    ) => {
      if (action.payload.index < 0 || action.payload.index > maxStatLength) {
        return state;
      }
      state.team[state.edit_index].stats[action.payload.index] =
        action.payload.val;
      return state;
    },
    addCharacter: (state, action: PayloadAction<{ name: string }>) => {
      if (state.team.length >= 4) return state;
      state.team.push(newChar(action.payload.name));
      return state;
    },
    deleteCharacter: (state, action: PayloadAction<{ index: number }>) => {
      if (
        action.payload.index < 0 ||
        action.payload.index >= state.team.length
      ) {
        return state;
      }
      state.team.splice(action.payload.index, 1);
      return state;
    },
    editCharacter: (state, action: PayloadAction<{ index: number }>) => {
      //if index < 0 then it means it's closed; but it shouldn't be bigger than max length
      if (action.payload.index >= state.team.length) {
        return state;
      }
      state.edit_index = action.payload.index;
      return state;
    },
  },
});

export const simActions = simSlice.actions;

export type SimSlice = {
  [simSlice.name]: ReturnType<typeof simSlice["reducer"]>;
};
