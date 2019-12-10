package terminal

import (
	"context"
	"fmt"

	"github.com/Masterlu1998/kube-viewer/kScrapper/resource"
	ui "github.com/gizak/termui/v3"
)

func workloadGraphAction(ctx context.Context, tdb *TerminalDashBoard, s *resource.ResourceScrapper, workloadTypes resource.ResourceTypes) {
	s.StartResourceScrapper(ctx, workloadTypes)
	t := tdb.Table
	for d := range s.GetDataCh() {
		workloadSData, ok := d.(resource.WorkloadData)
		if !ok {
			continue
		}
		t.Rows = [][]string{resourceTableHeader}
		fmt.Println(workloadSData.Infos)
		for _, wd := range workloadSData.Infos {
			var deploymentContent []string
			deploymentContent = append(deploymentContent,
				wd.Kind,
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
