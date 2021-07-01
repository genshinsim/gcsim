import { PremadeConfig } from "../SampleConfig";

export const ganyuAimOnly: PremadeConfig = {
  name: "ganyu w/ zhongli shield solo",
  description: "Ganyu spamming aimed shots w/ Zhongli shield",
  characters: ["ganyu", "diona", "zhongli", "fill"],
  tags: ["ganyu", "f2p"],
  data: `
char+=ganyu ele=cryo lvl=90 hp=9796.729 atk=334.850 def=630.215 cd=0.884 cr=.05 cons=0 talent=9,8,8;
weapon+=ganyu label="prototype crescent" atk=509.606 atk%=0.413 refine=1;
art+=ganyu label="blizzard strayer" count=4;
stats+=ganyu label=main hp=4780 atk=311 cryo%=0.466 atk%=0.466 cd=0.622;
stats+=ganyu label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;

char+=diona ele=cryo lvl=80 hp=8481.262 atk=188.235 def=532.293 cryo%=0.180 cr=.05 cd=0.5 cons=0 talent=1,1,1;
weapon+=diona label="sacrificial bow" atk=564.784 er=0.306 refine=1;

char+=zhongli ele=geo lvl=80 hp=13662.076 atk=233.483 def=685.947 geo%=0.288 cr=.05 cd=0.5 cons=0 talent=6,8,8;
weapon+=zhongli label="black tassel" atk=354.379 hp%=0.469 refine=5;
art+=zhongli  label="noblesse oblige" count=4;
stats+=zhongli label=main hp=4780 atk=311 hp%=.466 hp%=.466 hp%=.466;
stats+=zhongli label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;

char+=sucrose ele=anemo lvl=80 hp=8192.127 atk=150.588 def=623.025 anemo%=0.180 cr=.05 cd=0.5 cons=0 talent=1,1,1;
weapon+=sucrose label="sacrificial fragments" atk=454.363 em=220.512 refine=1;

target+="dummy" lvl=88 pyro=0.1 dendro=0.1 hydro=0.1 electro=0.1 geo=0.1 anemo=0.1 physical=.1;
active+=zhongli;

actions+=skill[hold=1] target=zhongli if=.tags.zhongli.shielded==0;

actions+=skill target=ganyu;
actions+=burst target=ganyu;
actions+=aim target=ganyu;
`,
};

export const ganyuAimVV: PremadeConfig = {
  name: "ganyu w/ vv",
  description: "Ganyu spamming aimed shots vv no shielding",
  characters: ["ganyu", "diona", "zhongli", "fill"],
  tags: ["ganyu", "f2p"],
  data: `
char+=ganyu ele=cryo lvl=90 hp=9796.729 atk=334.850 def=630.215 cd=0.884 cr=.05 cons=0 talent=9,8,8;
weapon+=ganyu label="prototype crescent" atk=509.606 atk%=0.413 refine=1;
art+=ganyu label="blizzard strayer" count=4;
stats+=ganyu label=main hp=4780 atk=311 cryo%=0.466 atk%=0.466 cd=0.622;
stats+=ganyu label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;

char+=diona ele=cryo lvl=80 hp=8481.262 atk=188.235 def=532.293 cryo%=0.180 cr=.05 cd=0.5 cons=0 talent=1,1,1;
weapon+=diona label="sacrificial bow" atk=564.784 er=0.306 refine=1;

char+=zhongli ele=geo lvl=80 hp=13662.076 atk=233.483 def=685.947 geo%=0.288 cr=.05 cd=0.5 cons=0 talent=6,8,8;
weapon+=zhongli label="black tassel" atk=354.379 hp%=0.469 refine=5;
art+=zhongli  label="noblesse oblige" count=4;
stats+=zhongli label=main hp=4780 atk=311 hp%=.466 hp%=.466 hp%=.466;
stats+=zhongli label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;

char+=sucrose ele=anemo lvl=80 hp=8192.127 atk=150.588 def=623.025 anemo%=0.180 cr=.05 cd=0.5 cons=0 talent=1,1,1;
weapon+=sucrose label="sacrificial fragments" atk=454.363 em=220.512 refine=1;

target+="dummy" lvl=88 pyro=0.1 dendro=0.1 hydro=0.1 electro=0.1 geo=0.1 anemo=0.1 physical=.1;
active+=zhongli;

actions+=skill[hold=1] target=zhongli if=.tags.zhongli.shielded==0;

actions+=skill target=ganyu;
actions+=burst target=ganyu;
actions+=aim target=ganyu;
  `,
};
