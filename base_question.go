package leetcodegraphql

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type BaseQuestion struct {
	QuestionId        string `json:"questionId"`
	QuestionTitle     string `json:"questionTitle"`
	Content           string `json:"content"`
	Difficulty        string `json:"difficulty"`
	DiscussUrl        string `json:"discussUrl"`
	CategoryTitle     string `json:"categoryTitle"`
	SubmitUrl         string `json:"submitUrl"`
	InterpretUrl      string `json:"interpretUrl"`
	CodeDefinition    string `json:"codeDefinition"`
	MetaData          string `json:"metaData"`
	EnvInfo           string `json:"envInfo"`
	Article           string `json:"article"`
	QuestionDetailUrl string `json:"questionDetailUrl"`
	DiscussCategoryId string `json:"discussCategoryId"`
}

func (q BaseQuestion) Valid() bool {
	return q.QuestionId != "" && q.QuestionTitle != ""
}

func (q BaseQuestion) GetCodeDefinition(lang string) (code string, err error) {
	s, err := strconv.Unquote(strconv.Quote(q.CodeDefinition))
	if err != nil {
		return
	}

	type Code struct {
		Text        string `json:"text"`
		Value       string `json:"value"`
		DefaultCode string `json:"defaultCode"`
	}
	var codes []Code
	if err = json.Unmarshal([]byte(s), &codes); err != nil {
		return
	}

	for _, c := range codes {
		if strings.ToLower(c.Text) == lang || c.Value == lang {
			code = c.DefaultCode
			return
		}
	}
	return
}

func (q BaseQuestion) GetEnvInfo(lang string) (info []string, err error) {
	s, err := strconv.Unquote(strconv.Quote(q.EnvInfo))
	if err != nil {
		return
	}

	m := make(map[string][]string)
	if err = json.Unmarshal([]byte(s), &m); err != nil {
		return
	}

	if temp, ok := m[lang]; ok && len(temp) != 0 {
		info = temp
	}
	return
}

func (q *BaseQuestion) Do(titleSlug string) error {
	body := strings.NewReader(`{"query":"query getQuestionDetail($titleSlug: String!) {\n  isCurrentUserAuthenticated\n  question(titleSlug: $titleSlug) {\n    questionId\n    questionTitle\n    content\n    difficulty\n    discussUrl\n    libraryUrl\n    mysqlSchemas\n    randomQuestionUrl\n    sessionId\n    categoryTitle\n    submitUrl\n    interpretUrl\n    codeDefinition\n    sampleTestCase\n    enableTestMode\n    metaData\n    enableRunCode\n    enableSubmit\n    judgerAvailable\n    emailVerified\n    envInfo\n    urlManager\n    likesDislikes {\n      likes\n      dislikes\n      __typename\n    }\n    article\n    questionDetailUrl\n    isLiked\n    discussCategoryId\n    nextChallengePairs\n    __typename\n  }\n}\n","variables":{"titleSlug":"` +
		titleSlug + `"},"operationName":"getQuestionDetail"}`)
	req, err := http.NewRequest("POST", "https://leetcode.com/graphql", body)
	if err != nil {
		return err
	}
	req.Header.Set("x-csrftoken", "uvORacsFvMydVNFluzue7hUMzM1F77MnYRbl4VBKTLBTQmxte9SWIYcM0mMJUovA")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("referer", fmt.Sprintf(
		"https://leetcode.com/problems/%s/description/",
		titleSlug,
	))
	client := &http.Client{
		Timeout: time.Second * 15,
	}
	req.AddCookie(&http.Cookie{
		Name:    "csrftoken",
		Value:   "uvORacsFvMydVNFluzue7hUMzM1F77MnYRbl4VBKTLBTQmxte9SWIYcM0mMJUovA",
		Path:    "/",
		Domain:  ".leetcode.com",
		Secure:  true,
		Expires: time.Now(),
	})
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resp := &Response{
		Data: ResponseData{
			Question: q,
		},
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		return err
	}

	//println(fmt.Sprintf("resp: %#v", resp))

	return nil
}