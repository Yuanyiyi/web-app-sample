package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"

	"github.com/sirupsen/logrus"

	"net"
	"time"
)

var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        10,               // 最大空闲连接数
		MaxIdleConnsPerHost: 3,                // 每个主机的最大空闲连接数
		IdleConnTimeout:     10 * time.Second, // 空闲连接的超时时间
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
	Timeout: 10 * time.Second,
}

func PostStruct(url string, bodyStruct interface{}) (res []byte, err error) {
	res, err = json.Marshal(bodyStruct)
	if err != nil {
		logrus.Errorf("[PostStruct]json.Marshal, url: %s, error: %s", url, err.Error())
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(res))
	if err != nil {
		logrus.Errorf("[PostStruct]NewRequest url: %s, error: %s", url, err.Error())
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("[PostStruct]Do url: %s, body: %v, error: %s", url, bodyStruct, err.Error())
		return
	}
	res, err = io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		logrus.Errorf("[PostStruct]response url: %s, error: %s", url, err.Error())
		return
	}

	return
}

func PostJson(url string, bodyJson string) (res []byte, err error) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Error(err)
		}
	}()
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(bodyJson)))
	if err != nil {
		logrus.Errorf("[PostJson]NewRequest url: %s, error: %s", url, err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("[PostJson]Do url: %s, body: %v, error: %s", url, bodyJson, err.Error())
		return
	}

	res, err = io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		logrus.Errorf("[PostJson]response url: %s, error: %s", url, err.Error())
		return
	}
	return
}

/*
*
发送的GET请求 需要设置header
Testner 20210123
*/
func GetJson(url string, authorization string) (res []byte, err error) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Error(err)
		}
	}()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Errorf("[GetJson]NewRequest url: %s, error: %s", url, err.Error())
		return
	}
	if authorization != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authorization))
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("[GetJson]Do url: %s, error: %s", url, err.Error())
		return
	}
	res, err = io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		logrus.Errorf("[GetJson]response url: %s, error: %s", url, err.Error())
		return
	}
	return
}

func CreateGetReqCtx(req interface{}, handlerFunc gin.HandlerFunc) (isSuccess bool, resp string) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	encode := structToURLValues(req).Encode()
	c.Request, _ = http.NewRequest("GET", "/?"+encode, nil)
	handlerFunc(c)
	return w.Code == http.StatusOK, w.Body.String()
}

func CreatePostReqCtx(req interface{}, handlerFunc gin.HandlerFunc) (isSuccess bool, resp string) {
	responseRecorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseRecorder)
	body, _ := json.Marshal(req)
	ctx.Request, _ = http.NewRequest("POST", "/", bytes.NewBuffer(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handlerFunc(ctx)
	return responseRecorder.Code == http.StatusOK, responseRecorder.Body.String()
}

// 将结构体转换为 URL 参数
func structToURLValues(s interface{}) url.Values {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	values := url.Values{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("form")
		if tag == "" {
			continue
		}

		value := v.Field(i).Interface()
		values.Set(tag, valueToString(value))
	}

	return values
}

// 由于 get 请求常常参数并不会特别复杂，通常的几种类型就应该可以包括，有需要可以继续添加
func valueToString(v interface{}) string {
	switch v := v.(type) {
	case int:
		return strconv.Itoa(v)
	case string:
		return v
	default:
		return ""
	}
}
