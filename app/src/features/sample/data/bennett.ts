import { PremadeConfig } from "../SampleConfig";

export const bennett4tf: PremadeConfig = {
  name: "4tf bennett",
  description: "4TF Bennett with Fiscl (Aluminum#5462)",
  characters: ["bennett", "beidou", "fischl", "xingqiu"],
  tags: ["4tf"],
  data: `

  char+=xingqiu ele=hydro lvl=80 hp=9514.469 atk=187.803 def=705.132 atk%=0.240 cr=.05 cd=0.5 cons=6 talent=6,8,8;
  weapon+=xingqiu label="sacrificial sword" atk=454.363 er=0.613 refine=1;
  art+=xingqiu label="heart of depth" count=2;
  art+=xingqiu label="noblesse oblige" count=2;
  stats+=xingqiu label=main hp=4780 atk=311 hydro%=0.466 er=0.518 cr=0.311;
  stats+=xingqiu label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
  
  char+=bennett ele=pyro lvl=80 hp=11538.824 atk=177.919 def=717.837 er=0.267 cr=.05 cd=0.5 cons=1 talent=6,8,8;
  weapon+=bennett label="festering desire" atk=509.606 er=0.459 refine=5;
  art+=bennett label="thundering fury" count=4;
  stats+=bennett  label=main hp=4780 atk=311 pyro%=0.466 er=0.518 cr=0.311;
  stats+=bennett label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
  
  char+=fischl ele=electro lvl=80 hp=8552.900 atk=227.341 def=552.671 atk%=0.240 cr=.05 cd=0.5 cons=4 talent=6,8,6;
  weapon+=fischl label="favonius warbow" atk=454.363 er=0.613 refine=1;
  art+=fischl label="gladiator's finale" count=2;
  art+=fischl label="thundering fury" count=2;
  stats+=fischl label=main hp=4780 atk=311 electro%=0.466 atk%=0.466 cr=0.311;
  stats+=fischl  label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
  
  char+=beidou ele=electro lvl=80 hp=11565 atk=200 def=575 cr=0.05 cd=0.50 electro%=.18 cons=6 talent=6,8,8;
  weapon+=beidou label="serpent spine" atk=510 refine=5 cr=.276;
  art+=beidou label="noblesse oblige" count=4;
  stats+=beidou label=main hp=4780 atk=311 electro%=0.466 er=0.518 cr=0.311;
  stats+=beidou label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
  
  target+="dummy" lvl=88 pyro=0.1 dendro=0.1 hydro=0.1 electro=0.1 geo=0.1 anemo=0.1 physical=.1;
  active+=bennett;
  
  actions+=sequence_strict target=xingqiu exec=skill,burst,attack;
  actions+=sequence_strict target=xingqiu exec=skill,attack,attack if=.status.xingqiu.energy<80;
  actions+=burst target=xingqiu;
  actions+=burst target=bennett;
  actions+=sequence_strict target=beidou exec=skill[counter=2],burst;
  actions+=burst target=fischl if=.status.fischloz==0;
  actions+=skill target=fischl if=.status.fischloz==0;
  actions+=sequence_strict target=bennett exec=skill,attack;
  
  
  actions+=attack target=bennett;
  actions+=attack target=fischl active=fischl;
  actions+=attack target=beidou active=beidou;
  actions+=attack target=xingqiu active=xingqiu;
`,
};

export const bennett4tfzajef: PremadeConfig = {
  name: "4tf bennett",
  description: "4TF Bennett with Lisa (Aluminum#5462)",
  characters: ["bennett", "beidou", "lisa", "xingqiu"],
  tags: ["4tf"],
  data: `

  char+=xingqiu ele=hydro lvl=80 hp=9514.469 atk=187.803 def=705.132 atk%=0.240 cr=.05 cd=0.5 cons=6 talent=6,8,8;
  weapon+=xingqiu label="sacrificial sword" atk=454.363 er=0.613 refine=1;
  art+=xingqiu label="heart of depth" count=2;
  art+=xingqiu label="noblesse oblige" count=2;
  stats+=xingqiu label=main hp=4780 atk=311 hydro%=0.466 er=0.518 cr=0.311;
  stats+=xingqiu label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
  
  char+=bennett ele=pyro lvl=80 hp=11538.824 atk=177.919 def=717.837 er=0.267 cr=.05 cd=0.5 cons=1 talent=6,8,8;
  weapon+=bennett label="festering desire" atk=509.606 er=0.459 refine=5;
  art+=bennett label="thundering fury" count=4;
  stats+=bennett  label=main hp=4780 atk=311 pyro%=0.466 er=0.518 cr=0.311;
  stats+=bennett label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
  
  char+=lisa ele=electro lvl=80 hp=8907.162 atk=215.480 def=533.613 em=96.000 cr=.05 cd=0.5 cons=4 talent=1,8,8;
  weapon+=lisa label="favonius sword" atk=509.606 er=0.459 refine=1;  
  art+=lisa label="gladiator's finale" count=2;
  art+=lisa label="thundering fury" count=2;
  stats+=lisa label=main hp=4780 atk=311 electro%=0.466 atk%=0.466 cr=0.311;
  stats+=lisa label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
  
  char+=beidou ele=electro lvl=80 hp=11565 atk=200 def=575 cr=0.05 cd=0.50 electro%=.18 cons=6 talent=6,8,8;
  weapon+=beidou label="serpent spine" atk=510 refine=5 cr=.276;
  art+=beidou label="noblesse oblige" count=4;
  stats+=beidou label=main hp=4780 atk=311 electro%=0.466 er=0.518 cr=0.311;
  stats+=beidou label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
  
  target+="dummy" lvl=88 pyro=0.1 dendro=0.1 hydro=0.1 electro=0.1 geo=0.1 anemo=0.1 physical=.1;
  active+=bennett;
  
  actions+=sequence_strict target=xingqiu exec=skill,burst,attack;
  actions+=sequence_strict target=xingqiu exec=skill,attack,attack if=.status.xingqiu.energy<80;
  actions+=burst target=xingqiu;
  actions+=burst target=bennett;
  actions+=sequence_strict target=beidou exec=skill[counter=2],burst;
  actions+=burst target=lisa;
  actions+=sequence_strict target=bennett exec=skill,attack;
  
  
  actions+=attack target=bennett;
  actions+=skill target=lisa active=lisa;
  actions+=attack target=lisa active=lisa;
  actions+=attack target=beidou active=beidou;
  actions+=attack target=xingqiu active=xingqiu;
`,
};
