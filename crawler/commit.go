package crawler

import (
	"auto-report/model"
	"encoding/json"
	"github.com/devfeel/mapper"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

/**
自动上报
*/
func Commit(username, password string) error {
	// 获取每日上报 Id
	client, err := Login(username, password)
	if err != nil {
		log.Println(err)
		return err
	}
	// 获取每日上报 Id
	id, err := getId(client)
	if err != nil {
		log.Println(err)
		return err
	}
	// 点击 “新增“，获取默认信息
	commitData, err := pressNewButton(client, id)
	if err != nil {
		log.Println(err)
		return err
	}
	// 提交
	commitInfo(client, commitData)
	return nil
}

// 获取每日上报 Id
func getId(client http.Client) (string, error) {
	req, err := http.NewRequest("POST", POSTID, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var commitResult model.CommitResult
	err = json.Unmarshal(body, &commitResult)
	if err != nil {
		log.Println(err)
		return "", err
	}
	id := commitResult.Module
	return id, nil
}

/**
点击 “新增“，获取默认信息
*/
func pressNewButton(client http.Client, id string) (model.CommitData, error) {
	reportId, err := json.Marshal(model.ID{id})
	if err != nil {
		log.Println(err)
		return model.CommitData{}, err
	}
	params := url.Values{
		"info": {string(reportId)},
	}
	resp, err := client.Post(REPORT, "application/x-www-form-urlencoded; charset=UTF-8", strings.NewReader(params.Encode()))
	if err != nil {
		log.Println(err)
		return model.CommitData{}, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var resultData model.ResultData
	err = json.Unmarshal(body, &resultData)
	if err != nil {
		log.Println(err)
		return model.CommitData{}, err
	}
	content := resultData.Module.Data[0]
	commitData := commit(&content)
	log.Println(commitData)
	return commitData, nil
}

/**
提交
*/
func commitInfo(client http.Client, commitData model.CommitData) error {
	cd, err := json.Marshal(commitData)
	if err != nil {
		log.Println(err)
		return err
	}
	params := url.Values{
		"info": {string(cd)},
	}
	resp, err := client.Post(COMMITURL, "application/x-www-form-urlencoded; charset=UTF-8", strings.NewReader(params.Encode()))
	if err != nil {
		log.Println(err)
		return err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var commitResult model.CommitResult
	err = json.Unmarshal(body, &commitResult)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(commitResult)
	return nil
}

/**
利用默认的信息构造提交数据
*/
func commit(reportData *model.ReportData) model.CommitData {
	modelData := &model.ModelData{}
	mapper.Mapper(reportData, modelData)
	commitData := model.CommitData{
		Model: *modelData,
	}
	return commitData
}
