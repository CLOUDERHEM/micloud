package authorizer

import (
	"io.github.clouderhem.micloud/authorizer/cookie"
	"io.github.clouderhem.micloud/cloud/setting/renewal"
	"io.github.clouderhem.micloud/consts"
	"io.github.clouderhem.micloud/utility/request"
	"net/http"
)

func DoRequest(req *http.Request) ([]byte, *http.Response, error) {
	var body []byte
	var resp *http.Response
	var err error
	for i := 0; i < consts.DefaultRetryTimes; i++ {
		postProcessReq(req)
		body, resp, err = request.DoRequest(req)
		if err != nil {
			continue
		}
		if resp.StatusCode == http.StatusUnauthorized {
			c, err := renewal.Renewal()
			if err != nil {
				return body, resp, err
			} else {
				cookie.SetCookie(c)
			}
		}
		break
	}

	return body, resp, err
}

func postProcessReq(req *http.Request) {
	req.Header.Set("Cookie", cookie.GetCookie())
	req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"128\", \"Not;A=Brand\";v=\"24\", \"Google Chrome\";v=\"128\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
}
