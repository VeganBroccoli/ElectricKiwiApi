package ElectricKiwiApi

import (
	"fmt"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const baseUrl = "https://my.electrickiwi.co.nz"

type HourOfPower int

const (
	Time0000 HourOfPower = iota
	Time0030
	Time0100
	Time0130
	Time0200
	Time0230
	Time0300
	Time0330
	Time0400
	Time0430
	Time0500
	Time0530
	Time0600
	Time0900
	Time0930
	Time1000
	Time1030
	Time1100
	Time1130
	Time1200
	Time1230
	Time1300
	Time1330
	Time1400
	Time1430
	Time1500
	Time1530
	Time1600
	Time2100
	Time2130
	Time2200
	Time2230
	Time2300
)

var hoursOfPower = [...]int{
	1,2,3,4,5,6,7,8,9,10,11,12,13,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,43,44,45,46,47,
}

var hoursOfPowerStr = [...]string{
	"0000", "0030", "0100", "0130", "0200", "0230", "0300", "0330", "0400", "0430", "0500", "0530", "0600",
	"0900", "0930", "1000", "1030", "1100", "1130", "1200", "1230", "1300", "1330", "1400", "1430", "1500", "1530", "1600",
	"2100", "2130", "2200", "2230", "2300",
}

func NewHourOfPower(str string) (HourOfPower, error) {
	for i := 0; i < len(hoursOfPowerStr); i++ {
		if hoursOfPowerStr[i] == str {
			return HourOfPower(i), nil
		}
	}

	return 0, fmt.Errorf("'%s' is not a valid time for the conversion", str)
}

func (h HourOfPower) ApiString() string {
	return strconv.Itoa(hoursOfPower[h])
}

type Session struct {
	client *http.Client
	billingStatusID string
}

var billingStatusRegex *regexp.Regexp

func init() {
	billingStatusRegex = regexp.MustCompile(`(name="active_billing_status_id" )(value="([^"]+)")`)
}

func NewSession(email, password string) (Session, error) {
	sess := Session{}

	opts := cookiejar.Options{PublicSuffixList: publicsuffix.List}
	jar, err := cookiejar.New(&opts)
	if err != nil {
		return sess, err
	}

	sess.client = &http.Client{
		Jar:           jar,
		Timeout:       0,
	}

	err = sess.login(email, password)
	if err != nil {
		return sess, err
	}

	return sess, nil
}

func (s Session) login(email, password string) error {
	form := url.Values{}
	form.Add("LoginForm[username]", email)
	form.Add("LoginForm[password]", password)

	req, err := postFormRequest("/Site/login", form)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	results := billingStatusRegex.FindStringSubmatch(string(body))
	s.billingStatusID = results[3] // gross :)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected 200 on login, received %d", resp.StatusCode)
	}

	return nil
}

func (s Session) UpdateHourOfPower(hour HourOfPower) error {
	form := url.Values{}
	form.Add("KiwikPayment[free_hour_consumption]", hour.ApiString())
	form.Add("active_billing_status_id", s.billingStatusID)

	req, err := postFormRequest("/account/update-hour-of-power", form)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected 200 on update to hour of power, got: %d", resp.StatusCode)
	}

	return nil
}

func postFormRequest(path string, form url.Values) (*http.Request, error){
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s%s", baseUrl, path),
		strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	return req, nil
}