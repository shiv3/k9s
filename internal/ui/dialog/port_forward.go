package dialog

import (
	"strings"

	"github.com/derailed/k9s/internal/ui"
	"github.com/derailed/tview"
	"github.com/gdamore/tcell"
)

const portForwardKey = "portforward"

// ShowPortForward pops a port forwarding configuration dialog.
func ShowPortForward(p *ui.Pages, port string, okFn func(address, lport, cport string)) {
	f := tview.NewForm()
	f.SetItemPadding(0)
	f.SetButtonsAlign(tview.AlignCenter).
		SetButtonBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetButtonTextColor(tview.Styles.PrimaryTextColor).
		SetLabelColor(tcell.ColorAqua).
		SetFieldTextColor(tcell.ColorOrange)

	p1, p2, address := port, port, "localhost"
	f.AddInputField("Pod Port:", p1, 20, nil, func(p string) {
		p1 = p
	})
	f.AddInputField("Local Port:", p2, 20, nil, func(p string) {
		p2 = p
	})
	f.AddInputField("Address:", address, 20, nil, func(h string) {
		address = h
	})

	f.AddButton("OK", func() {
		okFn(address, stripPort(p2), stripPort(p1))
	})
	f.AddButton("Cancel", func() {
		DismissPortForward(p)
	})

	modal := tview.NewModalForm("<PortForward>", f)
	modal.SetDoneFunc(func(_ int, b string) {
		DismissPortForward(p)
	})
	p.AddPage(portForwardKey, modal, false, false)
	p.ShowPage(portForwardKey)
}

// DismissPortForward dismiss the port forward dialog.
func DismissPortForward(p *ui.Pages) {
	p.RemovePage(portForwardKey)
}

// ----------------------------------------------------------------------------
// Helpers...

// StripPort removes the named port id if present.
func stripPort(p string) string {
	tokens := strings.Split(p, ":")
	if len(tokens) == 2 {
		return strings.Replace(tokens[1], "╱UDP", "", 1)
	}

	return p
}
