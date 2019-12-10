package terminal

import (
	"errors"

	"github.com/Masterlu1998/kube-viewer/dataTypes"
	ui "github.com/gizak/termui/v3"
	"github.com/sirupsen/logrus"
)

func deploymentGraphAction(tdb *TerminalDashBoard, s dataTypes.Scrapper) error {
	t := tdb.Table
	for d := range s.GetDataCh() {
		workloadSData, ok := d.(dataTypes.DeploymentScrapperChData)
		if !ok {
			logrus.Error("invalid chData type \"chData\"")
			return errors.New("invalid chData type \"chData\"")
		}
		t.Rows = make([][]string, 0)
		workloadTableHeader := []string{"name", "namespace", "pods", "create time", "images"}
		t.Rows = append(t.Rows, workloadTableHeader)

		for _, wd := range workloadSData.Deployments {
			var deploymentContent []string
			deploymentContent = append(deploymentContent,
				wd.Name,
				wd.Namespace,
				wd.PodsLive+"/"+wd.PodsTotal,
				wd.CreateTime,
				wd.Images)
			t.Rows = append(t.Rows, deploymentContent)
		}
		ui.Render(tdb.Grid)
	}

	return nil
}
