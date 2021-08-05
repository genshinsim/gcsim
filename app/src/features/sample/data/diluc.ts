import { PremadeConfig } from "../SampleConfig";

export const dilucvape: PremadeConfig = {
  name: "diluc vape",
  description: "Diluc vape",
  characters: ["diluc", "bennett", "xingqiu", "sucrose"],
  tags: ["vape"],
  data: `

  char+=diluc ele=pyro lvl=80 hp=12068.167 atk=311.310 def=728.818 cr=0.242 cr=.05 cd=0.5 cons=0 talent=8,8,8;
  weapon+=diluc label="prototype archaic" atk=564.784 atk%=0.276 refine=1;
  art+=diluc label="crimson witch of flames" count=4;
  stats+=diluc label=main hp=4780 atk=311 pyro%=0.466 atk%=0.466 cr=0.311;
  stats+=diluc label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
  
  char+=xingqiu ele=hydro lvl=80 hp=9514.469 atk=187.803 def=705.132 atk%=0.240 cr=.05 cd=0.5 cons=6 talent=6,8,8;
  weapon+=xingqiu label="sacrificial sword" atk=454.363 er=0.613 refine=1;
  art+=xingqiu label="heart of depth" count=2;
  art+=xingqiu label="noblesse oblige" count=2;
  stats+=xingqiu label=main hp=4780 atk=311 hydro%=0.466 er=0.518 cr=0.311;
  stats+=xingqiu label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
  
  char+=bennett ele=pyro lvl=80 hp=11538.824 atk=177.919 def=717.837 er=0.267 cr=.05 cd=0.5 cons=1 talent=6,8,8;
  weapon+=bennett label="festering desire" atk=509.606 er=0.459 refine=5;
  art+=bennett label="noblesse oblige" count=4;
  stats+=bennett  label=main hp=4780 atk=311 pyro%=0.466 er=0.518 cr=0.311;
  stats+=bennett label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
  
  char+=sucrose ele=anemo lvl=80 hp=8603.509 atk=158.150 def=654.311 anemo%=0.240 cr=.05 cd=0.5 cons=1 talent=1,8,8;
  weapon+=sucrose label="sacrificial fragments" atk=454.363 em=220.512 refine=1;
  art+=sucrose label="viridescent venerer" count=4;
  stats+=sucrose  label=main hp=4780 atk=311 em=187 em=187 em=187;
  stats+=sucrose label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
  
  
  target+="test dummy" lvl=88 pyro=0.1 dendro=0.1 hydro=0.1 electro=0.1 geo=0.1 anemo=0.1 physical=.1;
  active+=xingqiu;
  ##ROTATION
  actions+=sequence_strict target=xingqiu exec=skill[orbital=1],burst[orbital=1] lock=100;
  actions+=skill[orbital=1] target=xingqiu if=.status.xingqiu.energy<80 lock=100;
  actions+=burst[orbital=1] target=xingqiu;
  actions+=burst target=bennett;
  
  actions+=skill target=sucrose if=.element.pyro==1&&.debuff.res.vvpyro==0 label=vv;
  
  actions+=burst[dot=2,explode=0] target=diluc;
  actions+=sequence_strict target=diluc exec=attack,skill,attack,skill,attack,skill,attack if=.status.dilucq==1;
  actions+=attack target=diluc if=.status.dilucq==1;
  
  actions+=sequence_strict target=diluc exec=attack,skill,attack,skill,attack,skill,attack if=.cd.diluc.burst>100||.status.diluc.energy<20;
  actions+=skill target=bennett;
  
  actions+=attack target=diluc;
  actions+=attack target=xingqiu active=xingqiu;
  actions+=attack target=bennett active=bennett;
  actions+=attack target=sucrose active=sucrose;
  
  
`,
};
