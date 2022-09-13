package processbar

import (
	"errors"
	"fmt"
	"golang.org/x/term"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

/*
	Description % =======> | speed count rate time
*/

type processBar struct {
	config pbConfig
	lock   sync.Mutex
	state  state
}

type Theme struct {
	Saucer        string
	SaucerRight   string
	SaucerPadding string
	BarStart      string
	BarEnd        string
}

type showState struct {
	currentPercent float64
	currentBytes   float64
	timeSince      float64
	timeLeft       float64
	speedPerSecond float64
}

type state struct {
	// 当前状态
	currentNum        int64
	currentBytes      float64
	currentPercent    int
	currentSaucerSize int

	// 历史状态
	lastPercent int
	lastShown   time.Time

	counterNumSinceLast int64
	counterTime         time.Time
	counterLastTenRates []float64
	startTime           time.Time

	// 状态信息
	maxLineWidth int

	finished bool

	rendered string
}

type Option func(p *processBar)

type pbConfig struct {
	// 开关
	enable bool
	// 显示选项
	cleanAfterFinished bool
	max                int64
	writer             io.Writer
	fullWidth          bool
	showFileName       bool
	url                string

	// 输出格式
	customPrefix string
	theme        Theme

	// 自适应参数
	throttleDuration time.Duration
	width            int
}

func getBasicState() state {
	now := time.Now()
	return state{
		startTime:   now,
		lastShown:   now,
		counterTime: now,
	}
}

func (pb *processBar) Add64(num int64) error {
	pb.lock.Lock()
	defer pb.lock.Unlock()

	if pb.config.max <= 0 || pb.config.enable == true {
		return nil
	}

	pb.state.currentNum += num
	if pb.state.currentNum > pb.config.max {
		fmt.Println(pb.state.currentNum, pb.config.max)
		return errors.New("超过限制长度")
	}

	// 速率计算
	pb.state.counterNumSinceLast += num
	if time.Since(pb.state.counterTime).Seconds() > 0.5 {
		pb.state.counterLastTenRates = append(pb.state.counterLastTenRates, float64(pb.state.counterNumSinceLast)/time.Since(pb.state.counterTime).Seconds())
		if len(pb.state.counterLastTenRates) > 10 {
			pb.state.counterLastTenRates = pb.state.counterLastTenRates[1:]
		}
		pb.state.counterTime = time.Now()
		pb.state.counterNumSinceLast = 0
	}

	// 百分比计算
	percent := float64(pb.state.currentNum) / float64(pb.config.max)
	pb.state.currentPercent = int(percent * 100)
	pb.state.currentSaucerSize = int(percent * float64(pb.config.width))
	pb.state.currentBytes += float64(num)

	needUpdate := pb.state.lastPercent != pb.state.currentPercent && pb.state.currentNum > 0
	pb.state.lastPercent = pb.state.currentPercent

	// 渲染条件判断
	if needUpdate {
		return pb.render()
	}

	return nil
}

func (pb *processBar) render() error {
	// 渲染内容、配置填充
	if time.Since(pb.state.lastShown).Nanoseconds() < pb.config.throttleDuration.Nanoseconds() && pb.state.currentNum < pb.config.max {
		return nil
	}

	// 清理历史processBar
	clearProcessBar(pb.config, pb.state)

	// 检查是否finished
	if !pb.state.finished && pb.state.currentNum >= pb.config.max {
		pb.state.finished = true
		if !pb.config.cleanAfterFinished {
			// 执行渲染
			renderProcessBar(pb.config, &pb.state)
		}
	}

	if pb.state.finished {
		if pb.config.cleanAfterFinished {
			writeString(pb.config, "\r")
		} else {
			writeString(pb.config, "\n")
		}
		return nil
	}

	// 未完成的任务执行渲染
	w, err := renderProcessBar(pb.config, &pb.state)
	if err != nil {
		return err
	}
	if w > pb.state.maxLineWidth {
		pb.state.maxLineWidth = w
	}
	pb.state.lastShown = time.Now()
	return nil
}

func (pb *processBar) Add(n int) error {
	return pb.Add64(int64(n))
}

func (pb *processBar) Read(p []byte) (int, error) {
	n := len(p)
	return n, pb.Add(n)
}

func (pb *processBar) Write(p []byte) (int, error) {
	n := len(p)
	return n, pb.Add(n)
}

func (pb *processBar) Finish() error {
	pb.lock.Lock()
	pb.state.currentNum = pb.config.max
	pb.lock.Unlock()
	return nil
}

func (pb *processBar) Close() error {
	return pb.Finish()
}

func (pb *processBar) String() string {
	return pb.state.rendered
}

// settings

func OptionsCleanAfterFinish(enable bool) Option {
	return func(pb *processBar) {
		pb.config.cleanAfterFinished = enable
	}
}

func OptionsEnables(enable bool) Option {
	return func(pb *processBar) {
		pb.config.enable = enable
	}
}

func OptionsPrefix(prefix string) Option {
	return func(pb *processBar) {
		pb.config.customPrefix = prefix
	}
}

func OptionsFullWidth(enable bool) Option {
	return func(pb *processBar) {
		pb.config.fullWidth = enable
	}
}

func OptionsShowFileName(enable bool) Option {
	return func(pb *processBar) {
		pb.config.showFileName = enable
	}
}

func OptionsUrl(url string) Option {
	return func(pb *processBar) {
		pb.config.url = url
	}
}

// 初始化

func NewProcessBar(lens int64, opts ...Option) *processBar {
	pr := &processBar{
		state: getBasicState(),
		config: pbConfig{
			max:                lens,
			writer:             os.Stdout,
			cleanAfterFinished: true,
			customPrefix:       "Loading",
			width:              defaultWidth,
			theme:              defaultTheme,
			fullWidth:          true,
		},
	}
	for _, opt := range opts {
		opt(pr)
	}
	if pr.config.showFileName && pr.config.url != "" {
		filename := filepath.Base(pr.config.url)
		if filename != "" {
			pr.config.customPrefix = filename
		}
	}
	return pr
}

var defaultTheme = Theme{"-", ">", " ", "[", "]"}

const (
	defaultWidth = 80
)

// output

func getSize() (int, int) {
	if term.IsTerminal(0) {
		width, height, err := term.GetSize(int(os.Stdout.Fd()))
		if err == nil {
			return width, height
		}
	}
	return defaultWidth, 1
}

func getWidth() int {
	w, _ := getSize()
	return w
}

func renderProcessBar(c pbConfig, s *state) (int, error) {
	leftBrac := ""
	rightBrac := ""
	saucer := ""

	averageRate := average(s.counterLastTenRates)
	if len(s.counterLastTenRates) == 0 || s.finished {
		if t := time.Since(s.startTime).Seconds(); t > 0 {
			averageRate = s.currentBytes / t
		} else {
			averageRate = 0
		}
	}

	// 左侧
	leftBrac = fmt.Sprintf("%s %s %3.2f%%", s.startTime.Format("2006-01-02 15:04:05"), c.customPrefix, float64(s.currentNum*100)/float64(c.max))
	// 右侧
	currentHumanize, currentSuffix := humanizeBytes(s.currentBytes)
	speedHumanize, speedSuffix := humanizeBytes(averageRate)
	rightBrac += fmt.Sprintf("%s%s %s%s/s [%s]", currentHumanize, currentSuffix, speedHumanize, speedSuffix,
		(time.Duration(time.Since(s.startTime).Seconds()) * time.Second).String())
	// 中间
	if c.fullWidth {
		width := getWidth()
		c.width = width - 5 - len(leftBrac) - len(rightBrac)
		s.currentSaucerSize = int(float64(s.currentPercent) / 100.0 * float64(c.width))
	}

	if s.currentSaucerSize > 0 {
		saucer = c.theme.BarStart + strings.Repeat(c.theme.Saucer, s.currentSaucerSize-1) + c.theme.SaucerRight + strings.Repeat(c.theme.SaucerPadding, c.width-s.currentSaucerSize) + c.theme.BarEnd
	}

	// 拼接
	str := fmt.Sprintf("\r%s %s %s", leftBrac, saucer, rightBrac)

	return len(str), writeString(c, str)
}

func writeString(c pbConfig, str string) error {
	if _, err := io.WriteString(c.writer, str); err != nil {
		return err
	}
	if f, ok := c.writer.(*os.File); ok {
		f.Sync()
	}
	return nil
}

func clearProcessBar(c pbConfig, s state) error {
	str := fmt.Sprintf("\r%s\r", strings.Repeat(" ", s.maxLineWidth))
	return writeString(c, str)
}

func average(ss []float64) float64 {
	total := 0.0
	for _, v := range ss {
		total += v
	}
	return total / float64(len(ss))
}

func humanizeBytes(s float64) (string, string) {
	sizes := []string{" B", " kB", " MB", " GB", " TB", " PB", " EB"}
	base := 1024.0
	if s < 10 {
		return fmt.Sprintf("%2.0f", s), sizes[0]
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := math.Floor(float64(s)/math.Pow(base, e)*10+0.5) / 10
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}

	return fmt.Sprintf(f, val), suffix
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}
