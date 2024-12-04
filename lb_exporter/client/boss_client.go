package bossclient

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"exporter-demo/logger"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Boss2Client struct {
	url           string
	access_key    string
	secret_key    string
	ignore_verify bool
}

// create new client object
func NewClient(host string, port int, protocol string, access_key string, secret_key string, ignore_verify bool) (*Boss2Client, error) {
	// check protocol
	if protocol != "http" && protocol != "https" {
		return nil, errors.New("protocol should be http or https")
	}

	// format url
	url := fmt.Sprintf("%s://%s:%d/boss2/", protocol, host, port)
	return &Boss2Client{
		url:           url,
		access_key:    access_key,
		secret_key:    secret_key,
		ignore_verify: ignore_verify,
	}, nil
}

func (b2c *Boss2Client) GenSignature(access_key string, action string, paramsStr string, timestamp string) string {
	// format msg
	reqStr := fmt.Sprintf("access_key_id=%s&action=%s&params=%s&time_stamp=%s", access_key, action, url.QueryEscape(paramsStr), url.QueryEscape(timestamp))
	msg := fmt.Sprintf("POST\n/boss2/\n%s", reqStr)

	// sign with secret key
	digest := hmac.New(sha256.New, []byte(b2c.secret_key))
	digest.Write([]byte(msg))
	signature := base64.StdEncoding.EncodeToString(digest.Sum(nil))

	return signature
}

func (b2c *Boss2Client) Request(action string, params map[string]interface{}) (string, error) {
	// gen params string
	params["action"] = action
	paramsStr, _ := json.Marshal(params)

	// gen signature
	timestamp := time.Now().UTC().Format(time.RFC3339)
	signature := b2c.GenSignature(b2c.access_key, action, string(paramsStr), timestamp)

	// add signature to request body
	data := url.Values{}
	data.Set("action", action)
	data.Set("params", string(paramsStr))
	data.Set("access_key_id", b2c.access_key)
	data.Set("time_stamp", timestamp)
	data.Set("signature", signature)

	// send post request
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: b2c.ignore_verify},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Post(b2c.url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	if err != nil {
		logger.Error("an error occured when requesting to boss", err)
		fmt.Println("an error occured when requesting to boss", err)
	}
	defer resp.Body.Close()

	//Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("an error occured when getting response from boss, %v", err)
	}

	return string(responseBody), nil
}

// 负载均衡
func (b2c *Boss2Client) Boss2DescribeLoadBalancers(zone string, status []string, limit int, offset int) (string, error) {
	action := "Boss2DescribeLoadBalancers"

	// gen params
	params := map[string]interface{}{
		"zone":   zone,
		"status": status,
		"limit":  limit,
		"offset": offset,
		//"resources": []string{"lb-x1beyj7r"},
	}
	if len(status) == 0 || status == nil {
		params["status"] = []string{"active"}
	}

	// send requests
	resp, err := b2c.Request(action, params)
	if err != nil {
		logger.Error("Boss2DescribeLoadBalancers", err)
		fmt.Println("Boss2DescribeLoadBalancers", err)
	}

	return resp, nil

}

// 监听器
func (b2c *Boss2Client) DescribeLoadBalancerListeners(zone string, limit int, offset int) (string, error) {
	action := "DescribeLoadBalancerListeners"

	// gen params
	params := map[string]interface{}{
		"zone":   zone,
		"limit":  limit,
		"offset": offset,
		//"resources": []string{"lb-x1beyj7r"},
	}
	//if len(status) == 0 || status == nil {
	//	params["status"] = []string{"active"}
	//}

	// send requests
	resp, err := b2c.Request(action, params)
	if err != nil {
		logger.Error("DescribeLoadBalancerListeners", err)
		fmt.Println("DescribeLoadBalancerListeners", err)
	}

	return resp, nil

}

// 负载均衡
func (b2c *Boss2Client) GetLoadBalancerMonitor(zone, resource, step, protocol, subtype, startTime, endTime string, meters []string) (string, error) {
	action := "GetLoadBalancerMonitor"

	// gen params
	params := map[string]interface{}{
		"zone":       zone,
		"resource":   resource,
		"meters":     meters,
		"step":       step,
		"protocol":   protocol,
		"subtype":    subtype,
		"start_time": startTime,
		"end_time":   endTime,
		//"resources": []string{"lb-x1beyj7r"},
	}

	// send requests
	resp, err := b2c.Request(action, params)
	if err != nil {
		logger.Error("GetLoadBalancerMonitor", err)
		fmt.Println("GetLoadBalancerMonitor", err)
	}

	return resp, nil

}
