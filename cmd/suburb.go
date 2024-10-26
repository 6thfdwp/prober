package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/6thfdwp/prober/internal/housing"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/launcher"

	"github.com/spf13/cobra"
)

func checkErr(err error) {
	var evalErr *rod.EvalError
	if errors.Is(err, context.DeadlineExceeded) { // timeout error
		fmt.Println("timeout err")
	} else if errors.As(err, &evalErr) { // eval error
		fmt.Println(evalErr.LineNumber)
	} else if err != nil {
		fmt.Println("can't handle", err)
	}
}
func collectSupplyDemand(browser *rod.Browser, input string) string {

	subProfile := housing.NewSuburb(input)

	url := subProfile.ToREAFullUrl()
	log.Printf("## visiting %s", url)
	page := browser.MustPage(url).MustWaitLoad()
	log.Printf("## page loaded for %s", url)

	content := ""
	el := page.MustElement("[id='house-price-data-buy-all-bedrooms']")
	content = el.MustText()
	log.Printf("## page content got for %v", el)

	// err := rod.Try(func() {
	// 	// REA has id for each house type house-price-data-buy-4-bedrooms
	// 	el := page.MustElement("#house-price-data-buy-4-bedrooms")
	// 	content = el.MustText()
	// 	log.Printf("## page content got for %v", el)
	// })
	// checkErr(err)
	return content
}

func collectExtra(browser *rod.Browser, input string) string {
	parts := strings.Split(input, "-")
	l := len(parts)
	state, postcode := parts[l-2], parts[l-1]
	sub := postcode + "-" + strings.Join(parts[:l-2], "-")
	url := "https://www.yourinvestmentpropertymag.com.au/top-suburbs/" + state + "/" + sub
	log.Printf("## visiting %s", url)
	page := browser.MustPage(url).MustWaitLoad()
	log.Printf("## page loaded for %s", url)

	content := page.MustElement(".key-demographics").MustText()
	log.Printf("## done %s", url)
	return content
}
func collectMktInsights(browser *rod.Browser, suburb string) string {
	// use Domain.com as data source
	url := "https://www.domain.com.au/suburb-profile/" + suburb
	page := browser.MustPage(url).MustWaitLoad()
	log.Printf("## page %s loaded", url)

	page.MustScreenshotFullPage("./screenshots/sub.png")
	page.MustElement("[name='4 Bedroom House']").MustClick()
	mktInsights := page.MustElement(".suburb-insights").MustText()
	pops := page.MustElement("[data-testid='demographics']").MustText()

	log.Printf("## done %s", url)
	return mktInsights + "demographics: " + pops
}

var MyDevice = devices.Device{
	Title:          "Chrome computer",
	Capabilities:   []string{"touch", "mobile"},
	UserAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
	AcceptLanguage: "en",
	Screen: devices.Screen{
		DevicePixelRatio: 2,
		Horizontal: devices.ScreenSize{
			Width:  1500,
			Height: 900,
		},
		Vertical: devices.ScreenSize{
			Width:  1500,
			Height: 900,
		},
	},
}

func collectHouseProfile(browser *rod.Browser, url string) string {
	page := browser.MustPage(url).MustWaitLoad()
	log.Printf("## collecting house profile loaded %s", url)
	// page.MustElement(`[aria-label^="HeaderPropertyFeatures"]`).MustElements("a")
	est := page.MustElement("[aria-label='Property value']").MustText()
	// page.MustElement("[aria-label='Property timeline']")
	feat := page.MustElement("[aria-label='Property features']").MustText()
	// zones := page.MustElement("[aria-label='Government planning overlays & zones']").MustText()

	// return fmt.Sprintf("Estimated: %s, Features: %s, Gov overlays: %s", est, feat, zones)
	log.Printf("## house key info %s", fmt.Sprintf("Estimated: %s, Features: %s", est, feat))
	return fmt.Sprintf("Estimated: %s, Features: %s", est, feat)
}

func onExecSubStreet(fullSub string, street string) string {
	path, _ := launcher.LookPath()

	// u := "ws://127.0.0.1:9222/devtools/browser/fdcdabc6-0c90-48da-ab98-86824158bb4d"
	log.Printf("## browser path %s", path)
	// browser := rod.New().DefaultDevice(MyDevice).ControlURL(u).MustConnect()
	wsurl := launcher.NewUserMode().Bin(path).MustLaunch()
	browser := rod.New().ControlURL(wsurl).MustConnect()
	time.Sleep(3 * time.Second)
	defer browser.Close()

	subProfile := housing.NewSuburb(fullSub)
	url := subProfile.ToPropertyStreetUrl(street)
	log.Printf("## visiting street page %s", url)

	time.Sleep(2 * time.Second)
	page := browser.MustPage(url).MustWaitLoad()

	log.Printf("## street page %s loaded", url)
	// page.MustScreenshotFullPage("./screenshots/street.png")

	links := page.MustElements("a")
	// links := page.MustWaitElementsMoreThan("a", 3).MustElements("a")
	log.Printf("## links found %d, waiting 3 secs to continue ", len(links))
	time.Sleep(3 * time.Second)

	// res := make([]string, 0)
	res := map[string]string{}
	for _, link := range links {
		href := link.MustAttribute("href")

		if href == nil {
			log.Printf("## found invalid link ")
			continue
		}
		// only visiting house number we're interested in
		if !strings.Contains(*href, "45-pid") {
			continue
		}
		// !strings.Contains(*href, "10-pid") ||
		// !strings.Contains(*href, "20-pid") ||
		// !strings.Contains(*href, "30-pid") ||
		// !strings.Contains(*href, "45-pid") {
		// !strings.Contains(*href, "42-pid") || !strings.Contains(*href, "50-pid") || !strings.Contains(*href, "60-pid") {

		houseKeyInfo := collectHouseProfile(browser, subProfile.ToPropertyHouseUrl(*href))
		res[fullSub+"/"+street] = houseKeyInfo
		time.Sleep(1 * time.Second)
	}

	// ==
	// utils.Pause()
	d, _ := json.Marshal(res)
	return string(d)
}

func onExec(input []string) string {
	browser := rod.New().MustConnect()
	defer browser.Close()

	res := make(map[string]string)
	for _, suburb := range input {
		content := collectMktInsights(browser, suburb)
		extra := collectExtra(browser, suburb)

		// spd := collectSupplyDemand(browser, suburb)
		res[suburb] = content + "extra demographics:" + extra

	}
	d, _ := json.Marshal(res)
	return string(d)
}

func NewSuburbCmd() *cobra.Command {
	var suburbs, street string

	var suburbCmd = &cobra.Command{
		Use:   "suburb",
		Short: "Suburb command to probe suburb details",
		Long:  `Use this command to probe suburb details by passing suburb names.`,
		Run: func(cmd *cobra.Command, args []string) {
			if suburbs == "" {
				fmt.Println("No suburb names provided. Use -n flag to pass suburb names.")
				return
			}

			suburbList := strings.Split(suburbs, ",")
			if street != "" {
				log.Printf("## start exec suburb street probing %+v", street)
				output := onExecSubStreet(suburbList[0], street)

				fmt.Println(output)
				return
			}

			log.Printf("## start exec suburb profile probing %+v", suburbList)
			// Split the suburbs by commas and print each one
			output := onExec(suburbList)
			// send to standard output which can be piped
			fmt.Println(output)
		},
	}

	suburbCmd.Flags().StringVarP(&suburbs, "names", "n", "", "Comma-separated list of suburb names")
	suburbCmd.Flags().StringVarP(&street, "street", "s", "", "Street of suburb names")
	return suburbCmd
}
