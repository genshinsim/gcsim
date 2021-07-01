import { PremadeConfig } from "../SampleConfig";

export const eulaBasic: PremadeConfig = {
  name: "eula",
  description: "Eula with Diona + Kaeya + Fischl battery (Aluminum#5462)",
  characters: ["eula", "diona", "kaeya", "fischl"],
  tags: ["eula"],
  data: `
  char+=eula ele=cryo lvl=80 hp=12295.868 atk=317.981 def=698.094 cd=0.884 cr=.05 cons=0 talent=8,8,8;
  weapon+=eula label="serpent spine" atk=509.606 cr=0.276 refine=1;
  art+=eula label="pale flame" count=4;
  stats+=eula label=main hp=4780 atk=311 phys%=0.583 atk%=0.466 cd=0.622;
  stats+=eula label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
  
  char+=kaeya ele=cryo lvl=80 hp=10830.300 atk=207.572 def=736.894 er=0.267 cr=.05 cd=0.5 cons=0 talent=1,6,6;
  weapon+=kaeya label="favonius sword" atk=454.363 er=0.613 refine=1;
  art+=kaeya label="noblesse oblige" count=4;
  stats+=kaeya label=main hp=4780 atk=311 cryo%=0.466 er=0.518 cr=0.311;
  stats+=kaeya label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
  
  char+=diona ele=cryo lvl=80 hp=8907.162 atk=197.688 def=559.023 cryo%=0.240 cr=.05 cd=0.5 cons=0 talent=1,6,6;
  weapon+=diona label="sacrificial bow" atk=564.784 er=0.306 refine=1;
  art+=diona label="noblesse oblige" count=4;
  stats+=diona label=main hp=4780 atk=311 hp%=0.466 er=0.518 cr=0.311;
  stats+=diona label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
  
  char+=fischl ele=electro lvl=80 hp=8552.900 atk=227.341 def=552.671 atk%=0.240 cr=.05 cd=0.5 cons=4 talent=6,8,6;
  weapon+=fischl label="favonius warbow" atk=454.363 er=0.613 refine=1;
  art+=fischl label="gladiator's finale" count=2;
  art+=fischl label="thundering fury" count=2;
  stats+=fischl label=main hp=4780 atk=311 electro%=0.466 atk%=0.466 cr=0.311;
  stats+=fischl  label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
  
  
  target+="dummy" lvl=88 pyro=0.1 dendro=0.1 hydro=0.1 electro=0.1 geo=0.1 anemo=0.1 physical=.1;
  active+=kaeya;
  
  actions+=sequence_strict target=kaeya exec=skill,burst;
  
  actions+=burst target=fischl if=.status.fischloz==0;
  actions+=skill target=fischl if=.status.fischloz==0;
  
  actions+=sequence_strict target=eula exec=skill,burst,attack,attack,attack,attack,skill[hold=1],attack,attack,attack,attack;

  actions+=skill[hold=1] target=diona swap=eula;

  actions+=skill target=kaeya swap=eula;
  
  actions+=skill target=eula;
  actions+=sequence_strict target=eula exec=attack,attack,attack,attack if=.cd.kaeya.burst<300 lock=150 actionlock=210;
`,
};

export const eulaBennett: PremadeConfig = {
  name: "eula",
  description: "Eula with Diona + Bennett + Fischl battery",
  characters: ["eula", "diona", "bennett", "fischl"],
  tags: ["eula"],
  data: `
    char+=eula ele=cryo lvl=80 hp=12295.868 atk=317.981 def=698.094 cd=0.884 cr=.05 cons=0 talent=8,8,8;
    weapon+=eula label="serpent spine" atk=509.606 cr=0.276 refine=1;
    art+=eula label="pale flame" count=4;
    stats+=eula label=main hp=4780 atk=311 phys%=0.583 atk%=0.466 cd=0.622;
    stats+=eula label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
    
    char+=bennett ele=pyro lvl=80 hp=11538.824 atk=177.919 def=717.837 er=0.267 cr=.05 cd=0.5 cons=1 talent=6,8,8;
    weapon+=bennett label="festering desire" atk=509.606 er=0.459 refine=5;
    art+=bennett label="noblesse oblige" count=4;
    stats+=bennett  label=main hp=4780 atk=311 pyro%=0.466 er=0.518 cr=0.311;
    stats+=bennett label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
    
    char+=diona ele=cryo lvl=80 hp=8907.162 atk=197.688 def=559.023 cryo%=0.240 cr=.05 cd=0.5 cons=0 talent=1,6,6;
    weapon+=diona label="sacrificial bow" atk=564.784 er=0.306 refine=1;
    art+=diona label="noblesse oblige" count=4;
    stats+=diona label=main hp=4780 atk=311 hp%=0.466 er=0.518 cr=0.311;
    stats+=diona label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
    
    char+=fischl ele=electro lvl=80 hp=8552.900 atk=227.341 def=552.671 atk%=0.240 cr=.05 cd=0.5 cons=4 talent=6,8,6;
    weapon+=fischl label="favonius warbow" atk=454.363 er=0.613 refine=1;
    art+=fischl label="gladiator's finale" count=2;
    art+=fischl label="thundering fury" count=2;
    stats+=fischl label=main hp=4780 atk=311 electro%=0.466 atk%=0.466 cr=0.311;
    stats+=fischl  label=subs atk=50 atk%=.249 cr=.198 cd=.396 em=99 er=.257 hp=762 hp%=.149 def=59 def%=.186;
    
    
    target+="dummy" lvl=88 pyro=0.1 dendro=0.1 hydro=0.1 electro=0.1 geo=0.1 anemo=0.1 physical=.1;
    active+=fischl;

    actions+=skill target=fischl if=.tags.fischl.oz==0;
    actions+=burst target=fischl if=.tags.fischl.oz==0;
    
    actions+=skill target=bennett actionlock=600;
    actions+=sequence_strict target=bennett exec=burst if=.cd.eula.burst<120 swap=eula;
    
    actions+=sequence_strict target=eula exec=skill,burst,attack,attack,attack,attack,skill[hold=1],attack,attack,attack,attack;
    actions+=skill[hold=1] target=eula if=.tags.eula.grimheart==2;
    actions+=skill target=eula;
    
    
    actions+=skill[hold=1] target=diona if=.status.eula.energy<70 swap=eula;
    
    actions+=sequence_strict target=eula exec=attack,attack,attack,attack;
    actions+=attack target=fischl active=fischl;
    actions+=attack target=bennett active=bennett;
  `,
};
