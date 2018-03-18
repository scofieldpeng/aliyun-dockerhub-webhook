package main

import (
	"gopkg.in/gin-gonic/gin.v1"
)

type webhook struct {
	// 镜像名称
	Image string
	// 标签名称
	Tag string
	// 镜像的所在可用区
	Region string
}

// 阿里云webhook推送的原始内容
/*
{
    "push_data": {
        "digest": "sha256:457f4aa83fc9a6663ab9d1b0a6e2dce25a12a943ed5bf2c1747c58d48bbb4917",
        "pushed_at": "2016-11-29 12:25:46",
        "tag": "latest"
    },
    "repository": {
        "date_created": "2016-10-28 21:31:42",
        "name": "repoTest",
        "namespace": "namespace",
        "region": "cn-hangzhou",
        "repo_authentication_type": "NO_CERTIFIED",
        "repo_full_name": "namespace/repoTest",
        "repo_origin_type": "NO_CERTIFIED",
        "repo_type": "PUBLIC"
    }
}
*/
type webhookAliyunData struct {
	PushData struct {
		Digest   string `json:"digest"`
		PushTime string `json:"pushed_at"`
		Tag      string `json:"tag"`
	} `json:"push_data"`
	Repo struct {
		CreateTime     string `json:"date_created"`
		Name           string `json:"name"`
		NameSpace      string `json:"namespace"`
		Region         string `json:"region"`
		RepoAuthType   string `json:"repo_authentication_type"`
		RepoFullName   string `json:"repo_full_name"`
		RepoOriginType string `json:"repo_origin_type"`
		RepoType       string `json:"repo_type"`
	} `json:"repository"`
}

// 解析内容
// 会将阿里云webhook推送的内容进行解析，如果解析失败将会返回error
func (w *webhook) Parse(ctx *gin.Context) (err error) {
	var (
		originData webhookAliyunData
	)

	if err = ctx.BindJSON(&originData); err != nil {
		return
	}

	w.Image = originData.Repo.RepoFullName
	w.Tag = originData.PushData.Tag
	w.Region = originData.Repo.Region

	return
}
