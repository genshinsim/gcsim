import { PremadeConfig } from "../SampleConfig";

export const hutao4TF: PremadeConfig = {
  name: "4tf hutao",
  description: "4TF Hutao (Aluminum#5462)",
  characters: ["ganyu", "diona", "zhongli", "fill"],
  tags: ["hutao", "4tf"],
  data: `
char+=hutao ele=pyro lvl=90 hp=15552 atk=106 def=876 cr=0.05 cd=0.884 cons=1 talent=8,8,8 starthp=1;
weapon+=hutao label="staff of homa" atk=608 refine=1 cd=.662;
art+=hutao label="thundering fury" count=4;
stats+=hutao label=main hp=4780 atk=311 pyro%=0.466 em=187 cr=0.311;
stats+=hutao label=subs atk=50 hp%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 atk%=.149 def=59 def%=.186;

char+=beidou ele=electro lvl=80 hp=11565 atk=200 def=575 cr=0.05 cd=0.50 electro%=.18 cons=6 talent=6,8,7;
weapon+=beidou label="serpent spine" atk=510 refine=5 cr=.276;
art+=beidou label="gladiator's finale" count=2;
art+=beidou label="thundering fury" count=2;
stats+=beidou label=main hp=4780 atk=311 electro%=0.466 er=0.518 cr=0.311;
stats+=beidou label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;

char+=fischl ele=electro lvl=80 hp=8144 atk=216 def=526 cr=0.05 cd=0.50 atk%=.18 cons=6 talent=6,8,7;
weapon+=fischl label="skyward harp" atk=674 refine=1 cr=.221;
art+=fischl label="gladiator's finale" count=2;
art+=fischl label="thundering fury" count=2;
stats+=fischl label=main hp=4780 atk=311 electro%=0.466 er=0.518 cr=0.311;
stats+=fischl label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;

char+=xingqiu ele=hydro lvl=80 hp=9060 atk=179 def=671 cr=0.05 cd=0.50 atk%=.18 cons=6 talent=6,8,8;
weapon+=xingqiu label="favonius sword" atk=454 refine=4 er=.613;
art+=xingqiu label="noblesse oblige" count=2;
art+=xingqiu label="heart of depth" count=2;
stats+=xingqiu label=main hp=4780 atk=311 hydro%=0.466 er=0.518 cr=0.311;
stats+=xingqiu label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;

target+="primo geovishap" lvl=95 pyro=0.1 dendro=0.1 hydro=0.1 electro=0.1 geo=0.5 anemo=0.1 physical=.3;

active+=fischl;

actions+=skill target=fischl if=.status.fischloz==0;
actions+=burst target=fischl if=.cd.fischl.skill<1050;

actions+=sequence_strict target=beidou exec=skill[counter=2],burst lock=30;

actions+=sequence_strict target=xingqiu exec=skill[orbital=1],burst[orbital=1] lock=30;

actions+=skill target=hutao;
#actions+=burst target=hutao if=.element.hydro==1||.element.ec==1;

actions+=sequence_strict target=hutao if=.cd.hutao.skill>600 exec=
attack,attack,charge,dash,
attack,charge,dash,
attack,attack,charge,dash,
attack,charge,dash,
attack,attack,charge,dash,
attack,charge,jump,
attack,attack,charge,jump,
attack,charge,jump,
attack,attack,charge,jump;
`,
};
