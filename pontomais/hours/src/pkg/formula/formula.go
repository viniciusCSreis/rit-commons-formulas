// This is the formula implementation class.
// Where you will code your methods and manipulate the inputs to perform the specific operation you wish to automate.

package formula

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"
)

const ApiVersion = "2"

type WorkDays struct {
	TimeCards []struct {
		Date string `json:"date"`
		Time string `json:"time"`
	} `json:"time_cards"`
}

type wordDaysResponse struct {
	Days []WorkDays `json:"work_days"`
}

type auth struct {
	clientId string
	token    string
	email    string
}

type Formula struct {
	Username string
	Password string
	Client   http.Client
}

func (f Formula) Run() {
	auth, err := f.login()
	if err != nil {
		panic(err)
	}
	days, err := f.wordDays(auth)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Work Hours:\n")
	for _, d := range days {
		fmt.Printf("---\n")
		workTime := int64(0)
		for i, t := range d.TimeCards {
			if i == 0 {
				fmt.Printf("Data: %s\n", t.Date)
			}
			hour, _ := time.Parse("15:04", t.Time)
			if i%2 == 0 {
				workTime -= hour.UTC().Unix()
			} else {
				workTime += hour.UTC().Unix()
			}
			fmt.Printf("- %s \n", t.Time)
		}
		fmt.Printf("WorkTime: %s\n", time.Unix(workTime, 0).UTC().Format("15:04"))
	}

}

func (f Formula) wordDays(auth auth) ([]WorkDays, error) {
	body, err := json.Marshal(map[string]string{
		"login":    f.Username,
		"password": f.Password,
	})
	if err != nil {
		return nil, err
	}
	today := time.Now()
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"https://api.pontomais.com.br/api/time_card_control/current/work_days?sort_direction=asc&sort_property=date&start_date=%s&end_date=%s&with_employee=true",
			today.Add(time.Hour*24*7*-1).Format("2006-01-02"),
			today.Format("2006-01-02"),
		),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	request.Header.Add("token-type", "Bearer")
	request.Header.Add("api-version", ApiVersion)
	request.Header.Add("access-token", auth.token)
	request.Header.Add("client", auth.clientId)
	request.Header.Add("uid", auth.email)

	resp, err := f.Client.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode <= 199 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("status of login is %d", resp.StatusCode)
	}
	var rBody wordDaysResponse
	err = json.NewDecoder(resp.Body).Decode(&rBody)
	if err != nil {
		return nil, err
	}

	return reverse(rBody.Days), err
}

func (f Formula) login() (auth, error) {
	body, err := json.Marshal(map[string]string{
		"login":    f.Username,
		"password": f.Password,
	})
	if err != nil {
		return auth{}, err
	}
	request, err := http.NewRequest(
		http.MethodPost,
		"https://api.pontomais.com.br/api/auth/sign_in",
		bytes.NewReader(body),
	)
	if err != nil {
		return auth{}, err
	}
	request.Header.Add("api-version", ApiVersion)
	request.Header.Add("authority", "api.pontomais.com.br")
	request.Header.Add("content-type", "application/json;charset=UTF-8")
	request.Header.Add("accept", "application/json, text/plain, */*")
	fmt.Printf("making login\n")
	resp, err := f.Client.Do(request)
	if err != nil {
		return auth{}, err
	}
	if resp.StatusCode <= 199 || resp.StatusCode >= 300 {
		return auth{}, fmt.Errorf("status of login is %d", resp.StatusCode)
	}
	var rBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&rBody)
	if err != nil {
		return auth{}, err
	}
	return auth{
		token:    fmt.Sprintf("%s", rBody["token"]),
		clientId: fmt.Sprintf("%s", rBody["client_id"]),
		email:    fmt.Sprintf("%s", reflect.ValueOf(rBody["data"]).MapIndex(reflect.ValueOf("email"))),
	}, nil
}

func reverse(wd []WorkDays) []WorkDays {
	if len(wd) == 0 {
		return wd
	}
	return append(reverse(wd[1:]), wd[0])
}