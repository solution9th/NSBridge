package app

import (
	"go-agent/blueware"

	"github.com/solution9th/NSBridge/internal/utils"
)

func runOneAPM() (blueware.Application, error) {
	// dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	// if err != nil {
	// 	utils.Error("[apm] error:", err)
	// 	return nil, err
	// }
	// iniFilePath := dir + "/blueware-agent.ini"

	iniFilePath, ok := utils.FindFile("./blueware-agent.ini", "/etc/ns_bridge/blueware-agent.ini")
	if !ok {
		panic("not found agent ini")
	}

	cfg := blueware.NewConfig(iniFilePath)
	OneAPP, err := blueware.NewApplication(cfg)
	if nil != err {
		utils.Error("[apm] error:", err)
	}
	return OneAPP, err
}
