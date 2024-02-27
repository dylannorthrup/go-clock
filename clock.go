package main

// Use progress bars to show a clock
import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

// Things the clock models rely on
type tickMsg time.Time

const (
	timeout  = time.Second * 5
	padding  = 2
	maxWidth = 180
)

func tickEach100msCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg { return tickMsg(t) })
}

/// clockProgress: Setup progress bars for the clock

type clockProgress struct {
	progBar  progress.Model
	maxValue int
	value    int
	name     string
	barType  string
}

func (cp clockProgress) Init() tea.Cmd {
	return tickEach100msCmd()
}

func (cp clockProgress) View() string {
	pad := strings.Repeat(" ", padding)
	pct := float64(cp.value) / float64(cp.maxValue)
	cp.progBar.SetPercent(pct)
	return "\n" +
		pad + cp.progBar.ViewAs(pct) + "\n" +
		pad + cp.progBar.View()
}

func (cp clockProgress) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		cp.progBar.Width = msg.Width - padding*2 - 4
		if cp.progBar.Width > maxWidth {
			cp.progBar.Width = maxWidth
		}
		return cp, nil

	case time.Time:
		cmd := cp.progBar.SetPercent(float64(cp.value) / float64(cp.maxValue))
		return cp, tea.Batch(tickEach100msCmd(), cmd)

	default:
		fmt.Printf("Progress Message Type: %s\n", msg)
		if cp.value >= cp.maxValue {
			cp.value = 0
		}
		return cp, tickEach100msCmd()
	}
}

func NewClockProgress(name string, maxValue int) clockProgress {
	pb := progress.New(progress.WithWidth(maxWidth), progress.WithDefaultGradient(), progress.WithoutPercentage())
	m := clockProgress{
		progBar:  pb,
		maxValue: maxValue,
		name:     name,
		barType:  "clock",
	}
	return m
}

/// clockModel: Setup the whole of the clock

type clockModel struct {
	hour      clockProgress
	minute    clockProgress
	second    clockProgress
	modelType string
}

func (cp clockModel) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (cm clockModel) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + fmt.Sprintf("%d:%d:%d\n", cm.hour.value, cm.minute.value, cm.second.value) +
		cm.hour.View() +
		cm.minute.View() +
		cm.second.View()
}

func (cm clockModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		return tea.Model(cm), tea.Quit

	case tea.WindowSizeMsg:
		m := msg.(tea.WindowSizeMsg)
		cm.hour.progBar.Width = int(float64(m.Width) * 0.9)
		cm.minute.progBar.Width = int(float64(m.Width) * 0.9)
		cm.second.progBar.Width = int(float64(m.Width) * 0.9)
		return cm, nil

	case time.Time:
		h, m, s := time.Now().Clock()
		cm.hour.value = h
		cm.hour.progBar.Update(msg)
		cm.minute.value = m
		cm.second.value = s
		return cm, tickEach100msCmd()

	case nil:
		fmt.Printf("Model Message Type: %s :: %+v\n", msgType, msg)
		return cm, nil

	default:
		return cm, nil
	}
}

func NewClockModel() tea.Model {
	m := clockModel{
		hour:      NewClockProgress("hour", 24),
		minute:    NewClockProgress("minute", 60),
		second:    NewClockProgress("second", 60),
		modelType: "clock",
	}
	return tea.Model(m)
}

func RunClock() {
	m := NewClockModel()
	p := tea.NewProgram(m)
	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			p.Send(time.Now())
		}
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println("Oh no!", err)
		return
	}
}
