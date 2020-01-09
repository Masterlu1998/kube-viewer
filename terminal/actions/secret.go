package actions

import (
	"errors"

	"github.com/Masterlu1998/kube-viewer/kScrapper/common"
	"github.com/Masterlu1998/kube-viewer/kScrapper/secret"
)

var (
	secretTableHeader   = []string{"name", "namespace", "create time"}
	secretTableColWidth = []int{40, 25, 40}
)

func secretListDataGetter(c common.KubernetesData) ([]string, [][]string, []int, error) {
	secretInfos, ok := c.([]secret.Info)
	if !ok {
		return nil, nil, nil, errors.New("convert to secret.Info failed")
	}

	newSecretTableData := make([][]string, 0)
	for _, secretInfo := range secretInfos {
		newSecretTableData = append(newSecretTableData, []string{
			secretInfo.Name,
			secretInfo.Namespace,
			secretInfo.CreateTime,
		})
	}
	return secretTableHeader, newSecretTableData, secretTableColWidth, nil
}

func BuildSecretListAction() ActionHandler {
	return listResourceAction(secretListDataGetter, secret.SecretScrapperTypes)
}
