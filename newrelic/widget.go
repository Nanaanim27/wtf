package newrelic

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell"
	"github.com/olebedev/config"
	"github.com/rivo/tview"
	"github.com/senorprogrammer/wtf/wtf"
	nr "github.com/yfronto/newrelic"
)

var Config *config.Config

type Widget struct {
	wtf.TextWidget
}

func NewWidget() *Widget {
	widget := Widget{
		TextWidget: wtf.NewTextWidget("New Relic", "newrelic"),
	}

	widget.addView()

	return &widget
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) Refresh() {
	if widget.Disabled() {
		return
	}

	app, _ := Application()
	deploys, _ := Deployments()

	widget.View.SetTitle(fmt.Sprintf(" New Relic: [green]%s[white] ", app.Name))
	widget.RefreshedAt = time.Now()

	widget.View.Clear()
	fmt.Fprintf(widget.View, "%s", widget.contentFrom(deploys))
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) addView() {
	view := tview.NewTextView()

	view.SetBorder(true)
	view.SetBorderColor(tcell.ColorGray)
	view.SetDynamicColors(true)
	view.SetTitle(widget.Name)
	view.SetWrap(false)

	widget.View = view
}

func (widget *Widget) contentFrom(deploys []nr.ApplicationDeployment) string {
	str := fmt.Sprintf(
		" %s\n",
		"[red]Latest Deploys[white]",
	)

	revisions := []string{}

	for _, deploy := range deploys {
		if (deploy.Revision != "") && wtf.Exclude(revisions, deploy.Revision) {
			lineColor := "white"
			if wtf.IsToday(deploy.Timestamp) {
				lineColor = "cornflowerblue"
			}

			str = str + fmt.Sprintf(
				" [green]%s[%s] %s %-16s[white]\n",
				deploy.Revision[0:8],
				lineColor,
				deploy.Timestamp.Format("Jan 02, 15:04 MST"),
				wtf.NameFromEmail(deploy.User),
			)

			revisions = append(revisions, deploy.Revision)

			if len(revisions) == Config.UInt("wtf.newrelic.deployCount", 5) {
				break
			}
		}
	}

	return str
}
