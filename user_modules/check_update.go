package user_modules

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func CheckUpdate(appVersion string) {
	defer updateWG.Done()
	
	type GitResponse struct {
		TagName    string `json:"tag_name"`
		ReleaseUrl string `json:"html_url"`
	}
	
	var (
		apiResponse GitResponse
		err         error
		resp        *http.Response
	)
	
	httpClient := &http.Client{Timeout: 5 * time.Second}
	
	resp, err = httpClient.Get("https://api.github.com/repos/Underneach/Nicotine-String-Sorter/releases/latest")
	if err != nil {
		WaitLogo()
		PrintErr()
		fmt.Print("Failed to check updates: ", err, "\n\n")
		return
	}
	
	if resp.StatusCode != 200 {
		WaitLogo()
		PrintErr()
		fmt.Print("Failed to check updates: HTTP code ", resp.StatusCode, "\n\n")
		return
	}
	
	if err = json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		WaitLogo()
		PrintErr()
		fmt.Print("Failed to check updates: ", err, "\n\n")
		return
	}
	
	if apiResponse.TagName != appVersion {
		WaitLogo()
		PrintSuccess()
		fmt.Print("New version available: ")
		ColorBlue.Print(apiResponse.TagName, "\n")
		PrintSuccess()
		ColorBlue.Print(apiResponse.ReleaseUrl, "\n")
		PrintSuccess()
		fmt.Print("Download the new .exe and replace the current one\n\n")
	} else {
		WaitLogo()
		PrintSuccess()
		fmt.Print("You have the latest version of the sorter\n\n")
	}
	_ = resp.Body.Close()
}

func WaitLogo() {
	for !isLogoPrinted {
		time.Sleep(time.Millisecond * 100)
	}
}
