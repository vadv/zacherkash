package transport

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

var bodyRewriters = make([]*bodyRewrite, 0)

type bodyRewrite struct {
	reg *regexp.Regexp // то что матчится
	sub []byte         // на значение
}

func BuildBodyRewrites(settings map[string]string) error {
	bodyRewriters = make([]*bodyRewrite, 0)
	for s1, s2 := range settings {
		reg, err := regexp.Compile(s1)
		if err != nil {
			return err
		}
		bodyRewriters = append(bodyRewriters, &bodyRewrite{reg: reg, sub: []byte(s2)})
	}
	return nil
}

type BodyRewriter struct {
	http.RoundTripper
}

func (t *BodyRewriter) RoundTrip(req *http.Request) (resp *http.Response, err error) {

	log.Printf(`[INFO] %s - %s %s "%s" "%s"`, req.RemoteAddr, req.Method, req.URL, req.UserAgent(), req.Referer())

	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	// rewrite
	for _, rew := range bodyRewriters {
		if rew.reg.Match(b) {
			log.Printf("[INFO] Found `%v`, replace it with: `%s`\n", rew.reg, rew.sub)
			b = rew.reg.ReplaceAll(b, rew.sub)
		}
	}

	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))

	return resp, nil
}
