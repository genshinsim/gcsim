package compute

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/genshinsim/gcsim/pkg/model"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

type testResultHandler struct {
	res chan testRes
}

type testRes struct {
	key string
	err error
	res *model.SimulationResult
}

func (t *testResultHandler) Post(key string, callbackURL string, res *model.SimulationResult, err error) {
	t.res <- testRes{
		key: key,
		res: res,
		err: err,
	}

}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestComputeServer(t *testing.T) {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()

	//listener for res that blocks
	tt := testResultHandler{
		res: make(chan testRes),
	}

	server, err := New(Config{
		ResultHandler: &tt,
		APIKey:        "test",
		Timeout:       time.Minute * 10,
	})

	if err != nil {
		t.Fatal(err)
	}

	RegisterComputeServer(s, server)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	client := &Client{}
	client.cfg.APIKey = "test"

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client.srvClient = NewComputeClient(conn)

	err = client.Run("first", testCfg, context.Background())

	if err != nil {
		t.Fatal(err)
	}

	//fire off second request, resource should be busy

	err = client.Run("second", testCfg, context.Background())

	if err == nil {
		t.Errorf("expecting resource to be busy, got nil err")
	}

	r := <-tt.res

	if r.err != nil {
		t.Error(r.err)
	}

	fmt.Printf(
		"Simulation completed %v iterations\n",
		r.res.GetStatistics().GetIterations(),
	)

}

const testCfg = `
options swap_delay=12 debug=true iteration=1000 workers=50 mode=sl;

raiden char lvl=90/90 cons=0 talent=9,9,9;
raiden add weapon="dragonsbane" refine=3 lvl=90/90;
raiden add set="gildeddreams" count=5;
raiden add stats hp=4780 atk=311 em=559.5 ; #main
raiden add stats def%=0.124 def=39.36 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.0992 er=0.1102 em=118.92 cr=0.3972 cd=0.5296;

fischl char lvl=90/90 cons=6 talent=9,9,9;
fischl add weapon="thestringless" refine=3 lvl=90/90;
fischl add set="tf" count=5;
fischl add stats hp=4780 atk=311 em=187 electro%=0.466 cr=0.311 ; #main
fischl add stats def%=0.124 def=39.36 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.0992 er=0.1102 em=79.28 cr=0.331 cd=0.7944;

xingqiu char lvl=90/90 cons=6 talent=9,9,9;
xingqiu add weapon="lionsroar" refine=3 lvl=90/90;
xingqiu add set="emblemofseveredfate" count=4;
xingqiu add stats hp=4780 atk=311 atk%=0.466 hydro%=0.466 cr=0.311 ; #main
xingqiu add stats def%=0.124 def=39.36 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.0992 er=0.4408 em=39.64 cr=0.2648 cd=0.662;

nahida char lvl=90/90 cons=0 talent=9,9,9;
nahida add weapon="thewidsith" refine=3 lvl=90/90;
nahida add set="deepwoodmemories" count=4;
nahida add stats hp=4780 atk=311 em=186.5 dendro%=0.466 cr=0.311 ; #main
nahida add stats def%=0.124 def=39.36 hp=507.88 hp%=0.0992 atk=33.08 atk%=0.0992 er=0.1102 em=79.28 cr=0.331 cd=0.7944;


active raiden;
target lvl=100 resist=.1 hp=999999999999;
# target lvl=100 pos=1,0.5 resist=.1 hp=999999999999;
energy every interval=480,720 amount=1;

raiden skill;
for let i=0; i<4; i=i+1 {
    nahida skill, burst;
    fischl skill, attack;
    xingqiu burst, attack;
    nahida attack:2, skill,
            attack:3, dash,
            attack:3, dash,
            attack:3, charge;
    raiden attack, skill;
    fischl attack:2, burst;
    xingqiu attack, skill, dash, attack:2;
    nahida attack:2, skill,
            attack:3, dash,
            attack:3, dash,
            attack:3, dash,
            attack:3;
    print("done");
}

wait(12);

# changelog
# - pre-update dev: possum, RF
#   - stat solver: hessey
# - initial sim: vigne
# - 4tf 23.5s rotation: lettuce hunt
# - raiden E opener, N2/N3 optimization: skippi
# - lionsroar xingqiu: xardas (90k record)
# - sim wait() cutoff: skippi
# - wait(12) fischl burst: charlie
# - N1Q fischl, raiden before fischl: lettuce hunt

`
