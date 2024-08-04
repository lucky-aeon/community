package constant

/*
*
状态流转：

通过后由审核中转为报名中
会议中和已完成自动流转
*/
const (
	Reviewing   = "审核中"  // 会议申请默认该状态
	Registering = "报名中"  // 会议被审核通过后为报名中,报名需要有截止时间，不可将房间信息直接放出来
	Preparing   = "筹备中"  // 报名时间结束后进入筹备中
	InMeeting   = "会议中"  // 开会过程中的状态
	Completed   = "已完成"  // 会议结束
	Pass        = "PASS" // 会议不被通过
	Success     = "审核通过"
)
