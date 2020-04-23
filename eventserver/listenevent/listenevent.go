package listenevent

import (
	"os"

	"github.com/zhj0811/fabric/common/metadata"

	"github.com/op/go-logging"
)

var (
	logger = logging.MustGetLogger(metadata.LogModule)
)

//If the blocks'number listened from peer is not consequent, the program will exit to avoid losing some block.
func isConsequentNum(preNum, curNum uint64) {
	if 0 == preNum {
		logger.Debug("preNum is 0, it is the first listent block, no need to check!")
		return
	}

	if preNum+1 != curNum {
		logger.Errorf("preNum is %d and curNum is %d, which are not consequent!")
		logger.Errorf("The program exit to protect itself!")
		os.Exit(1)
	}

	logger.Debug("preNum and curNum is consequent.")
	return
}
