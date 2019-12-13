package terminal

import (
	"context"

	"github.com/Masterlu1998/kube-viewer/kScrapper"
	"github.com/Masterlu1998/kube-viewer/kScrapper/workload"
	ui "github.com/gizak/termui/v3"
)

func DeploymentGraphAction(ctx context.Context, tdb *TerminalDashBoard, sm *kScrapper.ScrapperManagement, namespace string) {
	s := sm.ScrapperMap[workload.ResourceScrapperTypes]
	s.StartScrapper(ctx, namespace)
	t := tdb.ResourceTable
	for d := range s.Watch() {
		workloadSData, ok := d.([]workload.WorkloadInfo)
		if !ok {
			continue
		}

		t.Rows = [][]string{{"name", "namespace", "pods", "create time", "images"}}
		for _, wd := range workloadSData {
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

}
