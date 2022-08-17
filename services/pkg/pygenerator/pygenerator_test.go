package pygenerator

import (
	"fmt"
	"testing"

	"github.com/genshinsim/gcsim/services/pkg/store"
)

func TestGenerate(t *testing.T) {

	t.Setenv("ASSETS_PATH", "../../../assets")

	g := &Generator{scriptPath: "../../embed/scripts/embed.py"}

	err := g.Generate(store.Simulation{Metadata: sample}, "result.png")

	if err != nil {
		fmt.Println("generate failed with error: ", err)
		t.FailNow()
	}

}

const sample = `{"char_names":["raiden","xingqiu","bennett","xiangling"],"dps":{"min":10835.347014196617,"max":13356.594095165565,"mean":12233.043251874682,"sd":395.4411340863557},"sim_duration":{"min":0,"max":105.35,"mean":105.35000000000113,"sd":0},"dps_by_target":{"1":{"min":0,"max":0,"mean":12234.97855376552,"sd":395.5036939873349}},"iter":1000,"runtime":2962000000,"char_details":[{"name":"raiden","element":"electro","level":90,"max_level":90,"cons":0,"weapon":{"name":"favoniuslance","refine":3,"level":90,"max_level":90},"talents":{"attack":9,"skill":9,"burst":9},"sets":{"tenacityofthemillelith":4},"stats":[0,0.124,39.36,5287.88,0.0992,344.08,0.1488,0.6833,39.64,0.642,0.7944,0,0,0,0,0.466,0,0,0,0,0,0],"snapshot":[0,0.124,39.36,5287.88,0.2992,344.08,0.3988,1.3095681723624013,39.64,0.6920000000000001,1.2944,0,0,0,0,0.986,0,0,0,0,0,0]},{"name":"xingqiu","element":"hydro","level":90,"max_level":90,"cons":6,"weapon":{"name":"harbingerofdawn","refine":5,"level":90,"max_level":90},"talents":{"attack":9,"skill":9,"burst":9},"sets":{"noblesseoblige":4},"stats":[0,0.124,39.36,5287.88,0.0992,344.08,0.5652,0.3306,39.64,0.5758,0.7944,0,0,0.466,0,0,0,0,0,0,0,0],"snapshot":[0,0.124,39.36,5287.88,0.0992,344.08,1.055199994635582,0.3306,39.64,0.9058,1.7629879772300723,0,0,0.666,0,0,0,0,0,0,0,0]},{"name":"bennett","element":"pyro","level":90,"max_level":90,"cons":6,"weapon":{"name":"thealleyflash","refine":1,"level":90,"max_level":90},"talents":{"attack":9,"skill":9,"burst":9},"sets":{"instructor":4},"stats":[0,0.124,39.36,4078.88,0.0992,265.08,0.0992,0.1102,226.64,0.5299,0.4634,0,0.348,0,0,0,0,0,0,0,0,0],"snapshot":[0,0.124,39.36,4078.88,0.0992,265.08,0.3492,0.3768999995708466,361.76794139671495,0.5799000000000001,0.9634,0,0.348,0,0,0,0,0,0,0,0,0.12]},{"name":"xiangling","element":"pyro","level":90,"max_level":90,"cons":6,"weapon":{"name":"thecatch","refine":5,"level":90,"max_level":90},"talents":{"attack":9,"skill":9,"burst":9},"sets":{"emblemofseveredfate":4},"stats":[0,0.124,39.36,5287.88,0.0992,344.08,0.0992,0.1102,266.28,0.642,0.7944,0,0.466,0,0,0,0,0,0,0,0,0],"snapshot":[0,0.124,39.36,5287.88,0.0992,344.08,0.3492,0.7695999931126831,362.28,0.6920000000000001,1.2944,0,0.466,0,0,0,0,0,0,0,0,0]}],"num_targets":1}`
