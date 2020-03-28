package led

import (
	"sync"
)

// Mode 灯的模式
type Mode byte

// 模式定义
const (
	ModeOff Mode = 1 << iota
	ModeOn
	ModeToggle
	ModeBlink
	ModeFlash
)

// 参数定义
const (
	BlinkContinuesTodo = 0
	BlinkTodo          = 1
	BlinkDutyCycle     = 5
	FlashCycleTime     = 1000 // ms
)

// Element led对象
type Element struct {
	mu        sync.Mutex
	mode      Mode
	todo      int
	onPct     int
	cycle     int
	next      int
	preStatus bool // 保存闪烁前的状态
	curStatus bool // 当前状态
	onOffFunc func(bool)
}

// Control led控制器
type Control struct {
	list []*Element
}

// NewLedControl 创建个新的led控制器
func NewLedControl() *Control {
	return &Control{}
}

// AddElement 注册led对象
func (sf *Control) AddElement(elements ...*Element) *Control {
	sf.list = append(sf.list, elements...)
	return sf
}

// Set 设置对象的模式
// ModeBlink: 闪烁1次,周期1s,占空比5%
// ModeFlash: 持烁闪烁,周期1s,占空比5%
// ModeOn,ModeOff,ModeToggle
func (sf *Element) Set(mode Mode) *Element {
	switch mode {
	case ModeBlink:
		sf.SetBlink(BlinkTodo, BlinkDutyCycle, FlashCycleTime)
	case ModeFlash:
		sf.SetBlink(BlinkContinuesTodo, BlinkDutyCycle, FlashCycleTime)
	case ModeOn:
		fallthrough
	case ModeOff:
		fallthrough
	case ModeToggle:
		sf.mu.Lock()
		if mode != ModeToggle {
			sf.mode = mode
		} else {
			sf.mode ^= ModeOn
		}
		sf.onOff(sf.mode == ModeOn)
		sf.mu.Unlock()
	default:
	}
	return sf
}

// NewElement 创建一个led对象
func NewElement(f func(bool)) *Element {
	return &Element{onOffFunc: f}
}

// SetBlink 设置闪烁
func (sf *Element) SetBlink(numBlinkTodo, duty, period int) *Element {
	if duty <= 0 || period <= 0 {
		sf.Set(ModeOff)
		return sf
	}
	if duty < 100 {
		sf.mu.Lock()
		if sf.mode < ModeBlink {
			// 保存闪烁之前的状态，闪烁完后恢复回去
			sf.preStatus = sf.curStatus
		}
		sf.mode = ModeOff
		sf.todo = numBlinkTodo
		sf.onPct = duty
		sf.cycle = period
		sf.next = 0
		if numBlinkTodo <= 0 {
			sf.mode |= ModeFlash
		}
		sf.mode |= ModeBlink
		sf.mu.Unlock()
	} else {
		sf.Set(ModeOff)
	}
	return sf
}

func (sf *Element) onOff(val bool) {
	sf.curStatus = val
	sf.onOffFunc(val)
}

// RunInterval 间隔运行,interval 为运行间隔的值
func (sf *Control) RunInterval(interval int) {
	var pct int

	for _, v := range sf.list {
		v.mu.Lock()
		if v.mode&ModeBlink == ModeBlink {
			if interval >= v.next {
				if v.mode&ModeOn == ModeOn {
					pct = 100 - v.onPct /* Percentage of cycle for off */
					v.mode &= ^ModeOn   /* Say it's not on */
					v.onOff(false)      /* Turn it off */
					if !(v.mode&ModeFlash == ModeFlash) {
						v.todo-- /* Not continuous, reduce count */
					}
				} else if v.todo == 0 && !(v.mode&ModeFlash == ModeFlash) {
					v.mode ^= ModeBlink /* No more blinks */
				} else {
					pct = v.onPct    /* Percentage of cycle for on */
					v.mode |= ModeOn /* Say it's on */
					v.onOff(true)    /* Turn it on */
				}

				if v.mode&ModeBlink == ModeBlink {
					v.next = v.cycle * pct / 100
				} else {
					/* no more blink, no more wait */
					v.next = 0
					/* After blinking, set the LED back to the state before it blinks */
					v.onOff(v.preStatus)
					/* Clear the status */
					v.preStatus = false
				}
			} else {
				v.next -= interval
			}
		}
		v.mu.Unlock()
	}
}
