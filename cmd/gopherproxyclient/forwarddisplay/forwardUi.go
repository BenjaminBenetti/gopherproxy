package forwarddisplay

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/CanadianCommander/gopherproxy/cmd/gopherproxyclient/proxy"
	"github.com/rivo/tview"
)

type ForwardUi struct {
	Running         bool
	RefreshInterval time.Duration
	clientMannager  *proxy.ClientManager

	uiApp         *tview.Application
	gridLayout    *tview.Grid
	forwardsTable *tview.Table
	clientList    *tview.List
	metrics       *tview.TextView
}

// ============================================
// Constructors
// ============================================

func NewForwardUi(clientManager *proxy.ClientManager) *ForwardUi {
	return &ForwardUi{
		Running:         false,
		RefreshInterval: 1 * time.Second,
		clientMannager:  clientManager,
	}
}

// ============================================
// Public Methods
// ============================================

func (ui *ForwardUi) Build() {
	ui.gridLayout = tview.NewGrid().SetRows(0, 1).SetColumns(-2, 0)

	// Forwards Table
	ui.forwardsTable = tview.NewTable()
	ui.forwardsTable.SetTitle("Forwarding Rules")
	ui.forwardsTable.SetBorder(true)

	// Clients Table
	ui.clientList = tview.NewList()
	ui.clientList.SetTitle("Channel Clients")
	ui.clientList.SetBorder(true)

	// Metrics bar
	ui.metrics = tview.NewTextView()
	ui.metrics.SetTitle("Metrics")

	// layout
	ui.gridLayout.AddItem(ui.forwardsTable, 0, 0, 1, 1, 0, 0, false)
	ui.gridLayout.AddItem(ui.clientList, 0, 1, 1, 1, 0, 0, false)
	ui.gridLayout.AddItem(ui.metrics, 1, 0, 1, 2, 0, 0, false)
}

func (ui *ForwardUi) StartDrawing() {
	ui.Running = true
	ui.uiApp = tview.NewApplication()
	go ui.drawLoop()
	if err := ui.uiApp.SetRoot(ui.gridLayout, true).SetFocus(ui.clientList).Run(); err != nil {
		panic(err)
	}
}

func (ui *ForwardUi) StopDrawing() {
	ui.Running = false
}

// ============================================
// Go Routines
// ============================================

func (ui *ForwardUi) drawLoop() {
	for ui.Running {
		// update UI values
		ui.updateFowardRulesTable()
		ui.updateClientsList()
		ui.updateMetrics()

		ui.uiApp.Draw()
		time.Sleep(ui.RefreshInterval)
	}
}

// ============================================
// Private Methods
// ============================================

func (ui *ForwardUi) updateFowardRulesTable() {
	ui.forwardsTable.Clear()
	for idx, rule := range ui.clientMannager.ForwardingRules {
		builder := strings.Builder{}
		fmt.Fprintf(&builder, "  %d -> %s -> %s:%d", rule.LocalPort, rule.RemoteClient, rule.RemoteHost, rule.RemotePort)

		str := builder.String()
		if !rule.Valid {
			str = "[red]" + str + " (offline)[-]"
		} else {
			str = "[green]" + str + " (online)[-]"
		}

		ui.forwardsTable.SetCell(idx, 0, tview.NewTableCell(str))
	}
}

func (ui *ForwardUi) updateClientsList() {
	for idx, client := range ui.clientMannager.StateManager.ChannelMembers {

		secondaryText := "Remote"
		if client.Id == ui.clientMannager.Client.Id {
			secondaryText = "You"
		}

		if idx < ui.clientList.GetItemCount() {
			ui.clientList.SetItemText(idx, client.Name, secondaryText)
		} else {
			ui.clientList.AddItem(client.Name, secondaryText, rune(client.Name[0]), nil)
		}
	}

	// remove any extra items
	itemDiff := ui.clientList.GetItemCount() - len(ui.clientMannager.StateManager.ChannelMembers)
	if itemDiff > 0 {
		for range itemDiff {
			ui.clientList.RemoveItem(ui.clientList.GetItemCount() - 1)
		}
	}

}

func (ui *ForwardUi) updateMetrics() {
	ui.metrics.Clear()

	builder := strings.Builder{}
	fmt.Fprintf(&builder, "Tx: %f Rx: %f", rand.Float64()*100, rand.Float64()*100)

	ui.metrics.SetText(builder.String())
}
