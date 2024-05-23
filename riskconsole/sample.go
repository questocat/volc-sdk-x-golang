package riskconsole

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func (p *RiskConsole) PushTrafficRiskData(req *PushTrafficRiskDataRequest) (*PushTrafficRiskDataResponse, error) {
	result := new(PushTrafficRiskDataResponse)
	fileInfo, err := getFileInfo(req.FilePath)
	if err != nil {
		return nil, err
	}

	getUploadIdRequest := GetUploadIdRequest{
		AppId:    req.AppId,
		FileName: fileInfo.FileName,
		FileHash: fileInfo.FileHash,
		PartSize: PartFileSize,
	}

	uploadIdResp, err := p.GetUploadId(getUploadIdRequest)
	if err != nil {
		return nil, err
	}

	getUploadedPartListRequest := GetUploadedPartListRequest{
		AppId:    req.AppId,
		UploadId: uploadIdResp.Data.UploadId,
	}

	uploadedPartListResp, err := p.GetUploadedPartList(getUploadedPartListRequest)
	if err != nil {
		return nil, err
	}

	partListMap := make(map[int]struct{})
	for _, partNum := range uploadedPartListResp.Data {
		partListMap[partNum] = struct{}{}
	}

	f, err := os.Open(req.FilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var bytesRead int
	buffer := make([]byte, PartFileSize)
	partNum := 1

	for {
		bytesRead, err = f.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if _, ok := partListMap[partNum]; ok {
			partNum++
			continue
		}

		uploadFileRequest := UploadFileRequest{
			AppId:    req.AppId,
			FileName: fileInfo.FileName,
			UploadId: uploadIdResp.Data.UploadId,
			PartNum:  partNum,
			PartSize: bytesRead,
		}

		partNum++

		if bytesRead == PartFileSize {
			uploadFileRequest.Content = buffer
		} else {
			newBuffer := make([]byte, bytesRead)
			copy(newBuffer, buffer[:bytesRead])
			uploadFileRequest.Content = newBuffer
		}

		var uploadFileResponse *UploadFileResponse
		uploadFileResponse, err = p.UploadFile(uploadFileRequest)
		if err != nil {
			return nil, err
		}

		if uploadFileResponse.ErrCode != "0" {
			result.Success = false
			return result, err
		}
	}

	completeUploadFileRequest := CompleteUploadFileRequest{
		AppId:        req.AppId,
		UploadId:     uploadIdResp.Data.UploadId,
		BusinessType: req.BusinessType,
		DataType:     req.DataType,
		Scene:        req.Scene,
	}

	completeUploadFileResult, err := p.CompleteUploadFile(completeUploadFileRequest)
	if err != nil {
		return nil, err
	}

	if completeUploadFileResult.ErrCode != "0" {
		result.Success = false
		return result, err
	}

	result.Success = true
	result.PackageId = completeUploadFileResult.Data

	return result, nil
}

func getFileInfo(filePath string) (*FileInfo, error) {
	if filePath == "" {
		return nil, fmt.Errorf("filePath is empty")
	}

	var file *os.File
	var err error
	file, err = os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	hash := md5.New()
	partNum := 0
	buffer := make([]byte, PartFileSize)
	var fileSize int64
	var bytesRead int
	for {
		bytesRead, err = file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		hash.Write(buffer[:bytesRead])
		fileSize += int64(bytesRead)
		partNum++
	}

	fileInfo := &FileInfo{
		FilePath: filePath,
		FileName: filepath.Base(filePath),
		PartNum:  partNum,
		FileHash: hex.EncodeToString(hash.Sum(nil)),
		FileSize: fileSize,
	}

	return fileInfo, nil
}
