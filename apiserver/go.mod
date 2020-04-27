module github.com/zhj0811/fabric/apiserver

go 1.14

require (
	github.com/DeanThompson/ginpprof v0.0.0-20190408063150-3be636683586
	github.com/fvbock/endless v0.0.0-20170109170031-447134032cb6
	github.com/gin-gonic/gin v1.6.2
	github.com/mailru/easyjson v0.7.1 // indirect
	github.com/peerfintech/gohfc v1.0.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/viper v1.6.3
	github.com/zhj0811/fabric v0.0.0-20200426030925-cfb6083aa555
	go.uber.org/zap v1.15.0
)

replace github.com/zhj0811/fabric => ../

replace github.com/peerfintech/gohfc v1.0.0 => ../gohfc
