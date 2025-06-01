package main

import (
	"fmt"
	"io"
	"net/http"
)

const (
	searchUrl              = "https://www.trackwrestling.com/AjaxFunctions.jsp"
	profilesSearchFunction = "getTWMembersJSP"
)

func TrackSearchProfiles(w http.ResponseWriter, r *http.Request) {
	sessionData, err := GetTrackSessionData()

	url := fmt.Sprintf("%s?TIM=%s&twSessionId=%s&function=%s&twId=%s&firstName=%s&lastName=%s&teamName=%s&hometown=%s",
		searchUrl,
		sessionData.Tim,
		sessionData.Id,
		profilesSearchFunction,
		r.URL.Query().Get("twId"),
		r.URL.Query().Get("firstName"),
		r.URL.Query().Get("lastName"),
		r.URL.Query().Get("teamName"),
		r.URL.Query().Get("hometown"),
	)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	Assert(err == nil, "error creating request")

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:138.0) Gecko/20100101 Firefox/138.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Sec-GPC", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", fmt.Sprintf("https://www.trackwrestling.com/membership/TWMemberList.jsp?TIM=%s&twSessionId=%s&fromDomain=0.0", sessionData.Tim, sessionData.Id))
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Priority", "u=0")

	resp, err := client.Do(req)
	Assert(err == nil, "error making request")

	body, err := io.ReadAll(resp.Body)
	Assert(err == nil, "error reading response body")

	fmt.Fprintln(w, string(body))

}
