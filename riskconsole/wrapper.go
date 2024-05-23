package riskconsole

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/volcengine/volc-sdk-golang/base"
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
	inputRequest := UploadFileInputRequest{
		AppId:    request.AppId,
		UploadId: request.UploadId,
		PartSize: request.PartSize,
		PartNum:  request.PartNum,
	}
	paramStr, err := json.Marshal(inputRequest)
	if err != nil {
		return nil, fmt.Errorf("UploadFileInputRequest: fail to marshal request, %v", err)
	}

	fieldItem := base.CreateMultiPartItemFormField("Input", string(paramStr))
	fileItem := base.CreateMultiPartItemFormFile("Data", request.FileName, bytes.NewReader(request.Content))
	resp, _, err := p.Client.CtxMultiPart(context.Background(), "UploadFile", nil, []*base.MultiPartItem{fieldItem, fileItem})
	if err != nil && p.Retry() {
		resp, _, err = p.Client.CtxMultiPart(context.Background(), "UploadFile", nil, []*base.MultiPartItem{fieldItem, fileItem})
	}

	if err != nil {
		return nil, fmt.Errorf("UploadFile: fail to do request, %v", err)
	}

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

	respBody, _, err := p.Client.Json(api, nil, string(reqData))
	if err != nil && p.Retry() {
		respBody, _, err = p.Client.Json(api, nil, string(reqData))
	}

	if err != nil {
		return respBody, fmt.Errorf("%s: fail to do request, %w", api, err)
	}

	return respBody, nil
}
