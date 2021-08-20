import { PremadeConfig } from "../SampleConfig";

export const xiaomango: PremadeConfig = {
  name: "xiao plunge spam",
  description: "xiao plunge spam (Mango#6990)",
  characters: ["xiao", "sucrose", "bennett", "fischl"],
  tags: ["xiao"],
  data: `
  #XIAO
  char+=xiao ele=anemo lvl=81 hp=11929.696 atk=327.099 def=748.709 cr=0.242 cd=0.5 cons=0 talent=8,8,8;
  weapon+=xiao label="primordial jade winged-spear" atk=674.335 cr=0.221 refine=1;
  art+=xiao label="gladiator's finale" count=2;
  art+=xiao label="viridescent venerer" count=2;
  stats+=xiao label=main hp=4780 atk=311 anemo%=0.466 atk%=0.466 cr=0.311;
  stats+=xiao label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
  
  #SUCROSE
  char+=sucrose ele=anemo lvl=60 hp=6501.184 atk=119.505 def=494.426 anemo%=0.120 cr=.05 cd=0.5 cons=6 talent=1,1,1;
  weapon+=sucrose label="sacrificial fragments" atk=454.363 em=220.512 refine=1;
  art+=sucrose label="viridescent venerer" count=4;
  stats+=sucrose label=main hp=4780 atk=311 em=187 em=187 em=187;
  stats+=sucrose label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
  
  #FISCHL
  char+=fischl ele=electro lvl=70 hp=7508 atk=200 def=485 cr=0.05 cd=0.50 atk%=.18 cons=6 talent=6,8,8;
  weapon+=fischl label="skyward harp" atk=674.335 cr=0.221 refine=1;
  art+=fischl label="gladiator's finale" count=2;
  art+=fischl label="thundering fury" count=2;
  stats+=fischl label=main hp=4780 atk=311 electro%=0.466 atk%=0.466 cd=0.622;
  stats+=fischl label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
  
  #BENNETT
  char+=bennett ele=pyro lvl=70 hp=10129 atk=156 def=630 cr=0.05 cd=0.50 er=.2 cons=2 talent=6,8,8;
  weapon+=bennett label="festering desire" atk=510 er=0.459 refine=5;
  art+=bennett label="noblesse oblige" count=4;
  stats+=bennett label=main hp=4780 atk=311 pyro%=0.466 hp%=0.466 cd=0.622;
  stats+=bennett label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
  
    target+="dummy" lvl=90 pyro=0.1 dendro=0.1 hydro=0.1 electro=0.1 geo=0.1 anemo=0.1 physical=0.1;
    active+=bennett;
  
  #bennett E and Q
  actions+=sequence_strict target=bennett exec=skill,burst if=.status.xiaoq==0;
  
  #fischl Q and E
  actions+=burst target=fischl if=.status.xiaoq==0&&(.tags.fischl.oz==0);
  actions+=skill target=fischl if=.status.xiaoq==0&&.tags.fischl.oz==0&&(.cd.fischl.burst<60);
  
  #sucrose battery
  actions+=skill target=sucrose if=.status.xiaoq==0&&.status.xiao.energy<70 swap=xiao;
  
  actions+=skill target=bennett if=.status.xiaoq==0&&.status.bennett.energy<60 lock=120;
  
  #xiao
  actions+=sequence_strict target=xiao exec=skill,skill,burst;
  actions+=sequence_strict target=xiao exec=skill,burst;
  actions+=burst target=xiao;
  
  actions+=high_plunge target=xiao if=.status.xiaoq==1;
  actions+=skill target=xiao if=.status.xiaoq==0&&.status.xiao.energy<70 lock=90;
  
  actions+=sequence_strict target=xiao exec=attack,charge;
  actions+=attack target=xiao;
  
  actions+=attack target=bennett active=bennett;
`,
};
