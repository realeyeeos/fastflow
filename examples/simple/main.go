package main

import (
	"fmt"
	"log"
	"time"

	"github.com/realeyeeos/fastflow"
	mongoKeeper "github.com/realeyeeos/fastflow/keeper/mongo"
	"github.com/realeyeeos/fastflow/pkg/entity/run"
	"github.com/realeyeeos/fastflow/pkg/mod"
	mongoStore "github.com/realeyeeos/fastflow/store/mongo"
)

type PrintAction struct {
}

// Name define the unique action identity, it will be used by Task
func (a *PrintAction) Name() string {
	return "PrintAction"
}
func (a *PrintAction) Run(ctx run.ExecuteContext, params interface{}) error {
	fmt.Println("action start: ", time.Now())
	return nil
}

func main() {
	// Register action
	fastflow.RegisterAction([]run.Action{
		&PrintAction{},
	})

	// init keeper, it used to e
	keeper := mongoKeeper.NewKeeper(&mongoKeeper.KeeperOption{
		Key: "worker-1",
		// if your mongo does not set user/pwd, you should remove it
		ConnStr:  "mongodb://root:pwd@127.0.0.1:27017/fastflow?authSource=admin",
		Database: "mongo-demo",
		Prefix:   "test",
	})
	if err := keeper.Init(); err != nil {
		log.Fatal(fmt.Errorf("init keeper failed: %w", err))
	}

	// init store
	st := mongoStore.NewStore(&mongoStore.StoreOption{
		// if your mongo does not set user/pwd, you should remove it
		ConnStr:  "mongodb://root:pwd@127.0.0.1:27017/fastflow?authSource=admin",
		Database: "mongo-demo",
		Prefix:   "test",
	})
	if err := st.Init(); err != nil {
		log.Fatal(fmt.Errorf("init store failed: %w", err))
	}

	go createDagAndInstance()

	// start fastflow
	if err := fastflow.Start(&fastflow.InitialOption{
		Keeper: keeper,
		Store:  st,
		// use yaml to define dag
		ReadDagFromDir: "./",
	}); err != nil {
		panic(fmt.Sprintf("init fastflow failed: %s", err))
	}
}

func createDagAndInstance() {
	// wait fast start completed
	time.Sleep(time.Second)

	// run some dag instance
	for i := 0; i < 10; i++ {
		_, err := mod.GetCommander().RunDag("test-dag", nil)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second * 10)
	}
}
