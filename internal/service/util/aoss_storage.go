package util

import (
	"context"
	"fmt"
	"os/exec"

	cmlog "gitlab.senseauto.com/apcloud/library/common-go/log"
)

const (
	downloadCmdStub    = "aws s3 cp %s %s --endpoint-url=%s"
	listCmdStub        = "aws s3 ls %s --endpoint-url=%s"
	s3BasePathTemplate = "s3://%s/"
)

type AOSSStorageClient struct {
	EndPoint string
	Cxt      context.Context
}

func (c *AOSSStorageClient) Download(bucket string, remoteFileSrc string, localFilePath string) error {
	// example: "aws s3 cp s3://aoss-bi-data/infra/infra_dev_test/test1.txt . --endpoint-url=http://aoss-v2.st-sh-01.sensecoreapi-oss.cn"
	logger := cmlog.Extract(c.Cxt)
	s3FileBasePath := fmt.Sprintf(s3BasePathTemplate, bucket)
	downloadCmd := fmt.Sprintf(downloadCmdStub, s3FileBasePath+remoteFileSrc, localFilePath, c.EndPoint)
	logger.Infof("the download cmd would be: %s", downloadCmd)

	cmd := exec.Command("sh", "-c", downloadCmd)
	result, err := cmd.Output()
	logger.Infof("the download cmd result: %s", string(result))
	if err != nil {
		logger.Errorf("the download cmd err: %s", err)
	}
	return err
}

func (c *AOSSStorageClient) List(bucket string, remoteFileSrc string) ([]byte, error) {
	// example: "aws s3 ls s3://aoss-bi-data/infra/infra_dev_test/ --endpoint-url=http://aoss-v2.st-sh-01.sensecoreapi-oss.cn"
	logger := cmlog.Extract(c.Cxt)
	s3FileBasePath := fmt.Sprintf(s3BasePathTemplate, bucket)
	listCmd := fmt.Sprintf(listCmdStub, s3FileBasePath+remoteFileSrc, c.EndPoint)
	logger.Infof("the list cmd would be: %s", listCmd)

	cmd := exec.Command("sh", "-c", listCmd)
	output, err := cmd.Output()

	logger.Infof("the download cmd result: %s", string(output))
	if err != nil {
		logger.Errorf("the list cmd err: %s", err)
	}
	return output, err
}
