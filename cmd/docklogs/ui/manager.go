package ui

import (
	"github.com/cinarmert/doclogs/cmd/docklogs/container"
	"github.com/gdamore/tcell"
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
)

type LayoutManager struct {
	App      *tview.Application
	Grid     *tview.Grid
	Sessions []*container.Session
}

func NewManagerForSessions(sessions []*container.Session) (*LayoutManager, error) {
	if len(sessions) == 0 {
		return nil, errors.New("could not create layour manager: no sessions given")
	}

	lm := (&LayoutManager{}).
		SetApp().
		SetSessions(sessions).
		SetGrid()

	return lm, nil
}

func (lm *LayoutManager) SetApp() *LayoutManager {
	lm.App = tview.NewApplication()
	return lm
}

func (lm *LayoutManager) SetSessions(sessions []*container.Session) *LayoutManager {
	lm.Sessions = sessions
	return lm
}

func (lm *LayoutManager) SetGrid() *LayoutManager {
	n := len(lm.Sessions)
	rows := (n + 1) / 2
	cols := 2

	if n == 1 {
		cols = 1
	}

	grid := tview.NewGrid().
		SetRows(make([]int, rows)...).
		SetColumns(make([]int, cols)...).
		SetBorders(true)

	lm.Grid = grid
	return lm
}

func (lm *LayoutManager) createTextView(title string) *tview.TextView {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		ScrollToEnd().
		SetChangedFunc(func() {
			lm.App.Draw()
		})

	textView.SetBorder(true)
	textView.SetTitle(" " + title + " ")
	textView.SetTitleColor(tcell.ColorRed)
	return textView
}

func (lm *LayoutManager) Run() {
	var wg sync.WaitGroup
	for i, session := range lm.Sessions {
		wg.Add(1)
		tv := lm.createTextView(session.Name)

		row, col, colSpan := i/2, i%2, 1
		if i == len(lm.Sessions)-1 && len(lm.Sessions)%2 == 1 {
			colSpan = 2
		}

		lm.Grid.AddItem(tv, row, col, 1, colSpan, 0, 0, false)

		go session.ReadLogs(&wg, tview.ANSIWriter(tv))
	}

	if !isatty.IsTerminal(os.Stdout.Fd()) {
		log.Warnf("doclogs is only available in a tty environment")
		return
	}

	if err := lm.App.SetRoot(lm.Grid, true).EnableMouse(true).Run(); err != nil {
		log.Fatalf("could not init ui: %v", err)
		os.Exit(1)
	}
}
