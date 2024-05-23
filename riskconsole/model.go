package riskconsole

import (
	"encoding/json"
	"fmt"
	"github.com/volcengine/volc-sdk-golang/base"
)

type PushTrafficRiskDataRequest struct {
	Scene        string
	BusinessType string
	DataType     string
	AppId        int64
	FilePath     string
}

type PushTrafficRiskDataResponse struct {
	Success   bool
	PackageId string
}

type FileInfo struct {
	FilePath string
	FileName string
	PartNum  int
	FileHash string
	FileSize int64
}

type GetUploadIdRequest struct {
	AppId    int64  `json:"app_id"`
	FileName string `json:"file_name"`
	FileHash string `json:"file_hash"`
	PartSize int    `json:"part_size"`
}

type CommonResponse struct {
	LogId     string `json:"log_id"`
	ErrCode   string `json:"err_code"`
	ErrMsg    string `json:"err_msg"`
	Timestamp int64  `json:"timestamp"`
}

type GetUploadIdResponse struct {
	CommonResponse
	Data GetUploadId `json:"data"`
}

type GetUploadId struct {
	UploadId string `json:"upload_id"`
	Status   int    `json:"status"`
}

type GetUploadedPartListRequest struct {
	AppId    int64  `json:"app_id"`
	UploadId string `json:"upload_id"`
}

type GetUploadedPartListResponse struct {
	CommonResponse
	Data []int `json:"data"`
}

type UploadFileRequest struct {
	AppId    int64  `json:"app_id"`
	UploadId string `json:"upload_id"`
	PartSize int    `json:"part_size"`
	Content  []byte `json:"content"`
	PartNum  int    `json:"part_num"`
	FileName string `json:"file_name"`
}

type UploadFileInputRequest struct {
	AppId    int64  `json:"app_id"`
	UploadId string `json:"upload_id"`
	PartSize int    `json:"part_size"`
	PartNum  int    `json:"part_num"`
}

type UploadFileResponse struct {
	CommonResponse
	Data string `json:"data"`
}

type CompleteUploadFileRequest struct {
	AppId        int64  `json:"app_id"`
	UploadId     string `json:"upload_id"`
	Scene        string `json:"scene"`
	BusinessType string `json:"business_type"`
	DataType     string `json:"data_type"`
}

type CompleteUploadFileResponse struct {
	CommonResponse
	Data string `json:"data"`
}

func UnmarshalResultInto(data []byte, result interface{}) error {
	resp := new(base.CommonResponse)
	if err := json.Unmarshal(data, resp); err != nil {
		return fmt.Errorf("fail to unmarshal response, %v", err)
	}
	errObj := resp.ResponseMetadata.Error
	if errObj != nil && errObj.CodeN != 0 {
		return fmt.Errorf("request %s error %s", resp.ResponseMetadata.RequestId, errObj.Message)
	}

	data, err := json.Marshal(resp.Result)
	if err != nil {
		return fmt.Errorf("fail to marshal result, %v", err)
	}

	if err = json.Unmarshal(data, result); err != nil {
		return fmt.Errorf("fail to unmarshal result, %v", err)
	}

	return nil
}
