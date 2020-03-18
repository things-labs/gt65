package module

// PressType 按下类型
type PressType byte

// 按下类型定义
const (
	PressDown PressType = iota
	PressUp
	PressLong
)

// KeyElement 按键对象
type KeyElement struct {
	FilterTime     int // 滤波时间
	LongTime       int // 长按时间,0表示不检测长按
	RepeatSpeed    int // 连击间隔,0表示不支持连击
	RepeatLongTime int // 如果不支持长按,但支持连击,有个长延时到连击的默认时间。
	IsDownFunc     func() bool
	CbFunc         func(PressType)

	longRepCnt  int // 长按连击共用计数器
	filterCount int // 滤波计数器
	state       int // 状态机
}

// KeyControl 按键控制器
type KeyControl struct {
	list []*KeyElement
}

const (
	stateStart = iota
	stateDown
	stateLong
	stateRepeat
	stateUp
)

// NewKeyControl 创建一个按键控制器
func NewKeyControl() *KeyControl {
	return &KeyControl{}
}

// RegisterKeyElement 注册一个对象
func (sf *KeyControl) RegisterKeyElement(elements ...*KeyElement) *KeyControl {
	sf.list = append(sf.list, elements...)
	return sf
}

// RunInterval 间隔时间运行
func (sf *KeyControl) RunInterval(interval int) {
	defer func() {
		_ = recover()
	}()

	for _, e := range sf.list {
		switch e.state {
		case stateStart:
			if e.IsDownFunc() {
				e.state = stateDown
				e.filterCount = 0
			}
		case stateDown:
			if e.filterCount += interval; e.filterCount >= e.FilterTime { // 滤波
				e.filterCount = 0
				if e.IsDownFunc() {
					if e.LongTime == 0 && e.RepeatSpeed == 0 { // 不支持长击和连击,直接到抬键状态
						e.CbFunc(PressDown)
						e.state = stateUp
					} else {
						e.state = stateLong
					}
				} else {
					e.state = stateStart
				}
			}
		case stateLong:
			if e.LongTime > 0 { // 支持长按
				if e.IsDownFunc() {
					if e.longRepCnt += interval; e.longRepCnt >= e.LongTime {
						e.CbFunc(PressLong)
						if e.RepeatSpeed == 0 { // 不支持连击，直接抬键
							e.state = stateUp
						} else {
							e.state = stateRepeat
						}
						e.longRepCnt = 0
					}
				} else { // 短按
					e.CbFunc(PressDown)
					e.state = stateUp
					e.longRepCnt = 0
				}
			} else { // 不支持长按
				if e.RepeatSpeed > 0 { // 支持连击
					if e.IsDownFunc() {
						if e.longRepCnt += interval; e.longRepCnt >= 50 {
							e.state = stateRepeat
							e.longRepCnt = 0
						}
					} else {
						e.CbFunc(PressDown)
						e.longRepCnt = 0
						e.state = stateUp
					}
				} else { // 不支持连击
					e.CbFunc(PressDown)
					e.state = stateUp
				}
			}
		case stateRepeat:
			if e.RepeatSpeed > 0 { // 支持连击
				if e.IsDownFunc() {
					if e.longRepCnt += interval; e.longRepCnt >= e.RepeatSpeed {
						e.CbFunc(PressDown)
						e.longRepCnt = 0
					}
				} else {
					e.longRepCnt = 0
					e.state = stateUp
				}
			} else {
				e.longRepCnt = 0
				e.state = stateUp
			}
		case stateUp:
			if e.filterCount += interval; e.filterCount >= e.FilterTime {
				if !e.IsDownFunc() {
					e.CbFunc(PressUp)
					e.state = stateStart
				}
				e.filterCount = 0
			}
		default:
			e.state = stateStart
		}
	}
}
