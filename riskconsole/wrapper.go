package riskconsole

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/volcengine/volc-sdk-golang/base"
	"strconv"
)

func (p *RiskConsole) GetUploadId(request GetUploadIdRequest) (*GetUploadIdResponse, error) {
	result := new(GetUploadIdResponse)
	respBody, err := p.performRequest("GetUploadId", request)
	if err != nil {
		return nil, err
	}

	if err = UnmarshalResultInto(respBody, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *RiskConsole) GetUploadedPartList(request GetUploadedPartListRequest) (*GetUploadedPartListResponse, error) {
	result := new(GetUploadedPartListResponse)
	respBody, err := p.performRequest("GetUploadedPartList", request)
	if err != nil {
		return nil, err
	}

	if err = UnmarshalResultInto(respBody, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *RiskConsole) UploadFile(request UploadFileRequest) (*UploadFileResponse, error) {
	appIdField := base.CreateMultiPartItemFormField("app_id", strconv.FormatInt(request.AppId, 10))
	uploadIdField := base.CreateMultiPartItemFormField("upload_id", request.UploadId)
	partSizeField := base.CreateMultiPartItemFormField("part_size", strconv.Itoa(request.PartSize))
	partNumField := base.CreateMultiPartItemFormField("part_num", strconv.Itoa(request.PartNum))
	fileNameField := base.CreateMultiPartItemFormField("file_name", request.FileName)
	fileContentField := base.CreateMultiPartItemFormFile("content", request.FileName, bytes.NewReader(request.Content))
	form := []*base.MultiPartItem{appIdField, uploadIdField, partSizeField, partNumField, fileNameField, fileContentField}
	fmt.Printf("UploadFileRequest request: app_id=%d,upload_id=%s,part_size=%d,part_num=%d,file_name=%s\n",
		request.AppId, request.UploadId, request.PartSize, request.PartNum, request.FileName)
	for _, item := range form {
		fmt.Printf("UploadFileRequest form: %+v\n", *item)
	}
	resp, _, err := p.Client.CtxMultiPart(context.Background(), "UploadFile", nil, form)
	if err != nil && p.Retry() {
		resp, _, err = p.Client.CtxMultiPart(context.Background(), "UploadFile", nil, form)
	}

	if err != nil {
		return nil, fmt.Errorf("UploadFile: fail to do request, %v", err)
	}

	fmt.Printf("UploadFileResponse response: %s\n", string(resp))

	result := new(UploadFileResponse)
	if err = UnmarshalResultInto(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *RiskConsole) CompleteUploadFile(request CompleteUploadFileRequest) (*CompleteUploadFileResponse, error) {
	result := new(CompleteUploadFileResponse)
	respBody, err := p.performRequest("CompleteUploadFile", request)
	if err != nil {
		return nil, err
	}

	if err = UnmarshalResultInto(respBody, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *RiskConsole) performRequest(api string, request interface{}) ([]byte, error) {
	reqData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("%sRequest: fail to marshal request, %w", api, err)
	}

	fmt.Printf("%sRequest request: %s\n", api, string(reqData))

	respBody, _, err := p.Client.Json(api, nil, string(reqData))
	if err != nil && p.Retry() {
		respBody, _, err = p.Client.Json(api, nil, string(reqData))
	}

	if err != nil {
		return respBody, fmt.Errorf("%s: fail to do request, %w", api, err)
	}

	fmt.Printf("%sRequest response: %s\n", api, string(respBody))

	return respBody, nil
}
