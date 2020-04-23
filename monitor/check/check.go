package check

import (
	"os"
	"time"

	"github.com/zhj0811/fabric/common/metadata"
	"github.com/zhj0811/fabric/common/sdk"
	"github.com/zhj0811/fabric/eventserver/handle"

	"github.com/op/go-logging"
)

const MaxFail = 3

var logger = logging.MustGetLogger(metadata.LogModule)

func CheckRecover(checkTime time.Duration) {
	var checkFail, blockHeight, blockRecvHeight, preBlockRecvHeight uint64

	for {
		time.Sleep(checkTime)
		logger.Debug("start to check height in block and record file.")

		blockRecv, err := handle.GetBlockInfo()
		if nil != err {
			logger.Errorf("get record height form file failed. %s", err.Error())
			continue
		}
		blockRecvHeight = blockRecv.BlockNumber

		if preBlockRecvHeight != blockRecvHeight {
			logger.Debugf("eventserver has received %d block from peer, it is normal.", blockRecvHeight)
			preBlockRecvHeight = blockRecvHeight
			checkFail = 0
			continue
		}

		blockHeight, _ = sdk.GetBlockHeightByEndorserPeer()
		blockHeight--
		if blockRecvHeight == blockHeight {
			logger.Debugf("The height(%d) in block is same as eventserver records.", blockHeight)
			checkFail = 0
			continue
		} else {
			logger.Warningf("The height(%d) in block is not same as eventserver records(%d).", blockHeight, blockRecvHeight)
			checkFail++
		}

		logger.Errorf("The program has check failed for %d times.", checkFail)
		if checkFail >= MaxFail {
			logger.Critical("The failed time of checking has reached the max time!")
			os.Exit(1)
		}
	}
}
