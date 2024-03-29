export const sampleConfig = `
# Bennett Chevreuse Fischl Yae. ~18s Rot. TF Bennett on-field DPS generates partcles for DPS build Chevreuse. Chevreuse Q3E every rotation in Bennett's burst. Fischl snapshots Bennett's buff. No Burst Yae. Bennett EDC>ADC>3EM.

# - Wepeaon/Set options ordered by DPS
bennett char lvl=90/90 cons=6 talent=9,9,9;
bennett add weapon="alleyflash" refine=1 lvl=90/90;
#bennett add weapon="wolffang" refine=1 lvl=90/90; # Personal DPS at the expense of team DPS
#bennett add weapon="umbrella" refine=5 lvl=90/90;
#bennett add weapon="ironsting" refine=5 lvl=90/90;
bennett add set="tf" count=4;
#bennett add stats hp=4780 atk=311 em=187 pyro%=0.466 cd=0.622;
#bennett add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.0992 er=0.1102 em=79.28 cr=0.3972 cd=0.662; #EM/CD
#bennett add stats hp=4780 atk=311 atk%=0.466 pyro%=0.466 cd=0.622 ; #main
#bennett add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.1984 er=0.1102 em=39.64 cr=0.3972 cd=0.662; #ATK/CD
bennett add stats hp=4780 atk=311 atk%=0.466 pyro%=0.466 cr=0.311 ; #main
bennett add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.0992 er=0.1102 em=79.28 cr=0.331 cd=0.7944; #EM/CR
#bennett add stats hp=4780 atk=311 atk%=0.466 pyro%=0.466 cr=0.311 ; #main
#bennett add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.1984 er=0.1102 em=39.64 cr=0.331 cd=0.7944; #AM/CR
#bennett add stats hp=4780 atk=311 em=187 em=187 em=187 ; #main
#bennett add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.0992 er=0.1102 em=118.92 cr=0.3972 cd=0.5296; #3EM

chev char lvl=90/90 cons=6 talent=9,9,9;
chev add weapon="deathmatch" refine=1 lvl=90/90;
chev add set="gt" count=4;
#chev add set="cw" count=4;
#chev add set="emblem" count=4;
#chev add stats hp=4780 atk=311 er=0.518 pyro%=0.466 cd=0.622 ; #main
#chev add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.0992 er=0.2204 em=39.64 cr=0.3972 cd=0.662 ; # Most ER
chev add stats hp=4780 atk=311 er=0.518 pyro%=0.466 cd=0.622; #main
chev add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.1488 er=0.1653 em=39.64 cr=0.3972 cd=0.662; #
#chev add stats hp=4780 atk=311 er=0.518 pyro%=0.466 cd=0.622 ; #main
#chev add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.1984 er=0.1102 em=39.64 cr=0.3972 cd=0.662; # Least ER

fischl char lvl=90/90 cons=6 talent=9,9,9;
fischl add weapon="viridescenthunt" refine=1 lvl=90/90;
#fischl add weapon="alleyhunter" refine=1 lvl=90/90 +params=[stacks=10];
#fischl add weapon="songofstillness" refine=5 lvl=90/90;
#fischl add weapon="stringless" refine=3 lvl=90/90;
fischl add set="gt" count=5;
fischl add stats hp=4780 atk=311 atk%=0.466 electro%=0.466 cd=0.622 ; #main
fischl add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.1984 er=0.1102 em=39.64 cr=0.3972 cd=0.662; # w/ VH
#fischl add stats hp=4780 atk=311 atk%=0.466 electro%=0.466 cr=0.311; #main
#fischl add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.1984 er=0.1102 em=39.64 cr=0.331 cd=0.7944; # w/ others

yae char lvl=90/90 cons=0 talent=9,9,9; 
yae add weapon="widsith" refine=3 lvl=90/90;
#yae add weapon="flowingpurity" refine=5 lvl=90/90;
yae add set="gt" count=4;
#yae add set="tom" count=4;
yae add stats hp=4780 atk=311 atk%=0.466 electro%=0.466 cr=0.311 ; #main
yae add stats def=39.36 def%=0.124 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.1984 er=0.1102 em=39.64 cr=0.331 cd=0.7944;

# Options
options iteration=1000; 
options swap_delay=12;
energy every interval=480,720 amount=1;

# Targets
target lvl=100 resist=0.1 radius=2 pos=0,2.4 hp=100000000;
#target lvl=100 resist=0.1 radius=2 pos=3.4,0 hp=100000000; 

# Rotation
active yae;

for let i=0; i<4; i=i+1 {
  yae skill:3; # No burst
  
  bennett skill, burst;
  
  chevreuse burst, skill[hold=1]:3; # 3 is always ideal
  
  if .fischl.skill.ready {
    fischl attack:1, skill;
  } else {
    fischl attack:2, burst;
  }
  
  bennett skill, attack;
  while .bennett.mods.bennett-field {
    if .bennett.skill.ready {
      bennett skill;
    } else {
      bennett attack;
    }
  }
  
}
`;
