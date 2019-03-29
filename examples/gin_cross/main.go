package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/TingYunAPM/go/utils/zip"

	"github.com/TingYunAPM/go"
	"github.com/TingYunAPM/go/framework/gin"
	"github.com/gin-gonic/gin"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("use : " + os.Args[0] + " <tingyun.json> <listenaddress> <url>")
		return
	}
	err := tingyun.AppInit(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}
	defer tingyun.AppStop()

	r := tingyun_gin.Default()
	r.GET("/extern", func(c *gin.Context) {
		action := tingyun_gin.FindAction(c)

		if track_id := c.Request.Header.Get("X-Tingyun-Id"); len(track_id) > 0 {
			fmt.Println(track_id)
			action.SetTrackId(track_id)
		}
		externalComponent := action.CreateExternalComponent(os.Args[3], "main.Extern")
		result, response, err := HttpGet(os.Args[3], map[string]string{"X-Tingyun-Id": externalComponent.CreateTrackId()})

		if response != nil {
			if tx_data := response.Header.Get("X-Tingyun-Tx-Data"); len(tx_data) > 0 {
				externalComponent.SetTxData(tx_data)
			}
			if tx_data := action.GetTxData(); len(tx_data) > 0 {
				c.Header("X-Tingyun-Tx-Data", tx_data)
			}
			c.Header("Content-Type", response.Header.Get("Content-Type"))
			c.String(response.StatusCode, string(result))
		} else {
			c.String(http.StatusOK, "/extern error:"+err.Error())
		}
		externalComponent.Finish()
	})
	//	r.Run(":8020")
	r.Run(os.Args[2])
}

func HttpGet(url string, params map[string]string) ([]byte, *http.Response, error) {
	duration := time.Second * 10
	var err error = nil
	request, err := http.NewRequest("GET", url, nil)
	if nil != err {
		return nil, nil, err
	}
	//	defer request.Body.Close()
	useParams := make(map[string]string)
	useParams["Accept-Encoding"] = "identity, deflate"
	//	useParams["Content-Type"] = "Application/json;charset=UTF-8"
	useParams["User-Agent"] = "TingYun-Agent/GoLang"
	for k, v := range params {
		useParams[k] = v
	}
	for k, v := range useParams {
		request.Header.Add(k, v)
	}

	client := &http.Client{Timeout: duration}
	response, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		if b, err := ioutil.ReadAll(response.Body); err != nil { //server返回200，然后读数据失败....
			return nil, response, err
		} else {
			for k, v := range response.Header {
				fmt.Print("Header: " + k + "=")
				for i := range v {
					fmt.Print(" " + v[i])
				}
				fmt.Println("")
			}
			encoding := response.Header.Get("Content-Encoding")
			if encoding == "gzip" || encoding == "deflate" {
				d, err := zip.Inflate(b)
				if err == nil {
					return d, response, nil
				}
			}
			return b, response, nil
		}
	} else {
		return nil, response, nil
	}

}
