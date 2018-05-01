package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"io"
	"os/exec"
	"time"
	"os/signal"
	"syscall"
)

// Generated with some json to go thing
type Subreddit struct {
	Kind string `json:"kind"`
	Data struct {
		After           string `json:"after"`
		Dist            int    `json:"dist"`
		Modhash         string `json:"modhash"`
		WhitelistStatus string `json:"whitelist_status"`
		Children        []struct {
			Kind string `json:"kind"`
			Data struct {
				SubredditID         string        `json:"subreddit_id"`
				ApprovedAtUtc       interface{}   `json:"approved_at_utc"`
				SendReplies         bool          `json:"send_replies"`
				ModReasonBy         interface{}   `json:"mod_reason_by"`
				BannedBy            interface{}   `json:"banned_by"`
				NumReports          interface{}   `json:"num_reports"`
				RemovalReason       interface{}   `json:"removal_reason"`
				Subreddit           string        `json:"subreddit"`
				SelftextHTML        string        `json:"selftext_html"`
				Selftext            string        `json:"selftext"`
				Likes               interface{}   `json:"likes"`
				SuggestedSort       interface{}   `json:"suggested_sort"`
				UserReports         []interface{} `json:"user_reports"`
				SecureMedia         interface{}   `json:"secure_media"`
				IsRedditMediaDomain bool          `json:"is_reddit_media_domain"`
				Saved               bool          `json:"saved"`
				ID                  string        `json:"id"`
				BannedAtUtc         interface{}   `json:"banned_at_utc"`
				ModReasonTitle      interface{}   `json:"mod_reason_title"`
				ViewCount           interface{}   `json:"view_count"`
				Archived            bool          `json:"archived"`
				Clicked             bool          `json:"clicked"`
				NoFollow            bool          `json:"no_follow"`
				Author              string        `json:"author"`
				NumCrossposts       int           `json:"num_crossposts"`
				LinkFlairText       string        `json:"link_flair_text"`
				ModReports          []interface{} `json:"mod_reports"`
				CanModPost          bool          `json:"can_mod_post"`
				IsCrosspostable     bool          `json:"is_crosspostable"`
				Pinned              bool          `json:"pinned"`
				Score               int           `json:"score"`
				ApprovedBy          interface{}   `json:"approved_by"`
				Over18              bool          `json:"over_18"`
				ReportReasons       interface{}   `json:"report_reasons"`
				Domain              string        `json:"domain"`
				Hidden              bool          `json:"hidden"`
				Preview             struct {
					Images []struct {
						Source struct {
							URL    string `json:"url"`
							Width  int    `json:"width"`
							Height int    `json:"height"`
						} `json:"source"`
						Resolutions []struct {
							URL    string `json:"url"`
							Width  int    `json:"width"`
							Height int    `json:"height"`
						} `json:"resolutions"`
						Variants struct {
						} `json:"variants"`
						ID string `json:"id"`
					} `json:"images"`
					Enabled bool `json:"enabled"`
				} `json:"preview"`
				Thumbnail           string        `json:"thumbnail"`
				Edited              bool          `json:"edited"`
				LinkFlairCSSClass   string        `json:"link_flair_css_class"`
				AuthorFlairCSSClass interface{}   `json:"author_flair_css_class"`
				ContestMode         bool          `json:"contest_mode"`
				Gilded              int           `json:"gilded"`
				Downs               int           `json:"downs"`
				BrandSafe           bool          `json:"brand_safe"`
				SecureMediaEmbed    struct {
				} `json:"secure_media_embed"`
				MediaEmbed struct {
				} `json:"media_embed"`
				AuthorFlairText       interface{} `json:"author_flair_text"`
				Stickied              bool        `json:"stickied"`
				Visited               bool        `json:"visited"`
				CanGild               bool        `json:"can_gild"`
				IsSelf                bool        `json:"is_self"`
				ParentWhitelistStatus string      `json:"parent_whitelist_status"`
				Name                  string      `json:"name"`
				Spoiler               bool        `json:"spoiler"`
				Permalink             string      `json:"permalink"`
				SubredditType         string      `json:"subreddit_type"`
				Locked                bool        `json:"locked"`
				HideScore             bool        `json:"hide_score"`
				Created               float64     `json:"created"`
				URL                   string      `json:"url"`
				WhitelistStatus       string      `json:"whitelist_status"`
				Quarantine            bool        `json:"quarantine"`
				SubredditSubscribers  int         `json:"subreddit_subscribers"`
				CreatedUtc            float64     `json:"created_utc"`
				SubredditNamePrefixed string      `json:"subreddit_name_prefixed"`
				Ups                   int         `json:"ups"`
				Media                 interface{} `json:"media"`
				NumComments           int         `json:"num_comments"`
				Title                 string      `json:"title"`
				ModNote               interface{} `json:"mod_note"`
				IsVideo               bool        `json:"is_video"`
				Distinguished         interface{} `json:"distinguished"`
			} `json:"data"`
		} `json:"children"`
		Before interface{} `json:"before"`
	} `json:"data"`
}

// Awesome Background Changing Desktop
func main() {
	fmt.Print(time.Now())
	fmt.Println(" - Starting application.")

	executePeriodically(time.Minute, run)
}

func run(index int) int {
	t := time.Now()
	fmt.Println(t)

	// Check here for more https://www.reddit.com/r/sfwpornnetwork/wiki/network
	subreddits := []string{"EarthPorn", "wallpapers", "wallpaper", "notitle", "itookapicture", "SkyPorn", "astrophotography", "spaceporn", "pic", "Cinemagraphs"}
	limit := 1
	pathToPic := "/Users/Tim/go/src/ABCD/output/background.jpg"

	if index >= len(subreddits) {
		index = 0
	}

	fmt.Println("Scraping subreddit ", subreddits[index])
	apiurl := fmt.Sprintf("https://www.reddit.com/r/%s/top.json?limit=%v", subreddits[index], limit)
	index++

	request, err := http.NewRequest("GET", apiurl, nil)
	request.Header.Set("User-Agent", "osx:scrapedbackground:v0.1 (by /u/zertykx)")
	if err != nil {
		fmt.Println("Error when creating new request: ", err)
		return index
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error when executing request: ", err)
		return index
	}

	var record Subreddit

	if err := json.NewDecoder(response.Body).Decode(&record); err != nil {
		fmt.Println(err)
	}

	if len(record.Data.Children) > 0 {
		for _, child := range record.Data.Children {
			data := child.Data
			fmt.Println("Title: " + data.Title)

			if len(data.Preview.Images) > 0 {
				for _, image := range data.Preview.Images {
					download(image.Source.URL, pathToPic)
					args := fmt.Sprintf("'tell application \"Finder\" to set desktop picture to \"%s\" as POSIX file'", pathToPic)
					exec.Command("osascript", "-e", args).Run()
					exec.Command("killall", "Dock").Run()
				}
			} else {
				fmt.Println("No images.")
			}
		}
	} else {
		fmt.Println("No records.")
	}

	return index
}

// No idea how this works.
// Copied from https://github.com/STAR-ZERO/tumbling/blob/master/tumbling/tumbling.go
func download(url string, out string) error {
	file, err := os.Create(out)
	if err != nil {
		return err
	}
	defer file.Close()

	resp, err := http.Get(url)
	if err != nil {
		os.Remove(file.Name())
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		os.Remove(file.Name())
		return err
	}

	return nil
}

func executePeriodically(d time.Duration, f func(int) int) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP)
	index := 0

	ticker := time.NewTicker(d)
	for {
		select {
		case <-ticker.C:
			index = f(index)
		case <-quit:
			fmt.Print(time.Now())
			fmt.Println(" - Quitting application.")
			ticker.Stop()
			return
		}
	}
}