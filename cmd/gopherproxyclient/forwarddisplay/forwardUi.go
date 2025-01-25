package forwarddisplay

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/CanadianCommander/gopherproxy/cmd/gopherproxyclient/proxy"
	"github.com/CanadianCommander/gopherproxy/internal/proxcom"
	"github.com/rivo/tview"
)

type ForwardUi struct {
	Running         bool
	RefreshInterval time.Duration
	clientManager   *proxy.ClientManager

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
		RefreshInterval: 250 * time.Millisecond,
		clientManager:   clientManager,
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
	var selectedChannelMember *proxcom.ChannelMember = nil
	forwardingRules := ui.clientManager.AllForwardingRules()

	// if possible identify the selected client and show only their forwarding rules
	if ui.clientList.GetCurrentItem() < len(ui.clientManager.StateManager.ChannelMembers) {
		selectedChannelMember = ui.clientManager.StateManager.ChannelMembers[ui.clientList.GetCurrentItem()]
		forwardingRules = selectedChannelMember.ForwardingRules

		ui.forwardsTable.SetTitle("Forwarding Rules - " + selectedChannelMember.Name)
	}

	// outgoing routes section
	ui.forwardsTable.SetCell(0, 0, tview.NewTableCell("[yellow]======== Outgoing Routes ========").
		SetAlign(tview.AlignCenter).
		SetSelectable(false))

	for idx, rule := range forwardingRules {
		builder := strings.Builder{}
		fmt.Fprintf(&builder, "  %d -> %s -> %s:%d", rule.LocalPort, rule.RemoteClient, rule.RemoteHost, rule.RemotePort)

		str := builder.String()
		if !rule.Valid {
			str = "[red]" + str + " (offline)[-]"
		} else {
			str = "[green]" + str + " (online)[-]"
		}

		ui.forwardsTable.SetCell(idx+1, 0, tview.NewTableCell(str))
	}

	// incoming routes section
	if selectedChannelMember != nil {
		ui.forwardsTable.SetCell(len(forwardingRules)+1, 0, tview.NewTableCell("[yellow] ======== Incoming Routes ========").
			SetAlign(tview.AlignCenter).
			SetSelectable(false))

		incomingRules := ui.clientManager.AllForwardingRulesTargetingClient(selectedChannelMember.Name)

		for idx, rule := range incomingRules {
			builder := strings.Builder{}
			fmt.Fprintf(&builder, "  %s:%d <- %s <- %d", rule.RemoteHost, rule.RemotePort, rule.RemoteClient, rule.LocalPort)
			str := builder.String()
			if !rule.Valid {
				str = "[red]" + str + " (offline)[-]"
			} else {
				str = "[green]" + str + " (online)[-]"
			}

			ui.forwardsTable.SetCell(idx+len(forwardingRules)+2, 0, tview.NewTableCell(str))
		}
	}
}

func (ui *ForwardUi) updateClientsList() {
	for idx, client := range ui.clientManager.StateManager.ChannelMembers {

		secondaryText := "Remote"
		if client.Id == ui.clientManager.Client.Id {
			secondaryText = "You"
		}

		if idx < ui.clientList.GetItemCount() {
			ui.clientList.SetItemText(idx, client.Name, secondaryText)
		} else {
			ui.clientList.AddItem(client.Name, secondaryText, rune(client.Name[0]), nil)
		}
	}

	// remove any extra items
	itemDiff := ui.clientList.GetItemCount() - len(ui.clientManager.StateManager.ChannelMembers)
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
