import { PremadeConfig } from "../SampleConfig";

export const ningandco: PremadeConfig = {
  name: "ning and co",
  description: "Ningguang and friends",
  characters: ["ningguang", "albedo", "zhongli", "xingqiu"],
  tags: ["ningguang", "geo", "gems"],
  data: `

  char+=zhongli ele=geo lvl=80 hp=13662.076 atk=233.483 def=685.947 geo%=0.288 cr=.05 cd=0.5 cons=0 talent=6,8,8;
  weapon+=zhongli label="black tassel" atk=354.379 hp%=0.469 refine=5;
  art+=zhongli  label="tenacity of millelith" count=4;
  stats+=zhongli label=main hp=4780 atk=311 hp%=.466 hp%=.466 hp%=.466;
  stats+=zhongli label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
  
  char+=bennett ele=pyro lvl=80 hp=11538.824 atk=177.919 def=717.837 er=0.267 cr=.05 cd=0.5 cons=1 talent=6,8,8;
  weapon+=bennett label="favonius sword" atk=454.363 er=0.613 refine=1;
  art+=bennett label="noblesse oblige" count=4;
  stats+=bennett  label=main hp=4780 atk=311 pyro%=0.466 er=0.518 cr=0.311;
  stats+=bennett label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
  
  char+=ningguang ele=geo lvl=80 hp=9109.598 atk=197.688 def=533.613 geo%=0.240 cr=.05 cd=0.5 cons=6 talent=8,8,8;
  weapon+=ningguang label="solar pearl" atk=509.606 cr=0.276 refine=1;
  art+=ningguang label="gladiator's finale" count=2;
  art+=ningguang label="archaic petra" count=2;
  stats+=ningguang label=main hp=4780 atk=311 geo%=.466 atk%=.466 cd=.622;
  stats+=ningguang label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
  
  char+=albedo ele=geo lvl=80 hp=12295.868 atk=233.483 def=814.562 geo%=0.288 cr=.05 cd=0.5 cons=0 talent=6,8,8;
  weapon+=albedo label="festering desire" atk=509.606 er=0.459 refine=5;
  art+=albedo label="noblesse oblige" count=4;
  stats+=albedo label=main hp=4780 atk=311 geo%=.466 def%=.583 cr=.311;
  stats+=albedo  label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
  
  target+="dummy" lvl=88 pyro=0.1 dendro=0.1 hydro=0.1 electro=0.1 geo=0.1 anemo=0.1 physical=.1;
  active+=zhongli;
  
  actions+=skill[hold=1] target=zhongli if=.tags.zhongli.shielded==0;
  actions+=skill target=albedo if=.tags.albedo.elevator==0;
  
  actions+=burst target=bennett;
  
  actions+=charge target=ningguang if=.tags.ningguang.jade==7;
  actions+=sequence_strict target=ningguang exec=skill,burst if=.status.noelleq==0;
  actions+=skill target=ningguang if=.status.ningparticles==0;
  
  actions+=burst target=albedo;
  
  actions+=burst target=ningguang;
  
  actions+=skill target=bennett actionlock=600;
  
  actions+=charge target=ningguang if=.tags.ningguang.jade>1;
  actions+=attack target=ningguang;
  
  actions+=attack target=zhongli active=zhongli;
  actions+=attack target=albedo active=albedo;
  actions+=attack target=bennett active=bennett;  
  
`,
};
