package dataTypes

const DeploymentScrapperTypes = "deployment"

type DeploymentScrapperChData struct {
	Deployments []DeploymentInfo
}

type DeploymentInfo struct {
	Name       string
	Namespace  string
	PodsLive   string
	PodsTotal  string
	CreateTime string
	Images     string
}
