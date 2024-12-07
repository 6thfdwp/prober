package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/6thfdwp/prober/internal/housing"
	"github.com/go-rod/rod"
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
func collectSupplyDemand(browser *rod.Browser, subProfile housing.SuburbProfile) string {
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

func collectExtra(suburb housing.SuburbProfile) string {
	browser := rod.New().MustConnect()
	defer browser.Close()

	url := suburb.ToYourInvestFullUrl()
	page := browser.MustPage(url).MustWaitLoad()
	log.Printf("## visiting %s", url)

	content := page.MustElement(".key-demographics").MustText()
	log.Printf("## done %s", suburb.Name)
	return content
}
func collectMktInsights(suburb housing.SuburbProfile) string {
	browser := rod.New().MustConnect()
	defer browser.Close()
	// use Domain.com as data source
	url := suburb.ToDmainFullUrl()
	page := browser.MustPage(url).MustWaitLoad()
	log.Printf("## visiting page %s loaded", url)

	page.MustScreenshotFullPage("./outputs/screenshots/sub.png")
	page.MustElement("[name='4 Bedroom House']").MustClick().MustWaitLoad()
	time.Sleep(1 * time.Second)
	mktInsights := page.MustElement(".suburb-insights").MustText()
	pops := page.MustElement("[data-testid='demographics']").MustText()

	log.Printf("## done for %s", suburb.Name)
	return mktInsights + "demographics: " + pops
}

func collectHouseProfile(browser *rod.Browser, url string) string {
	page := browser.MustPage(url).MustWaitLoad()
	log.Printf("## collecting house profile loaded %s", url)
	// page.MustElement(`[aria-label^="HeaderPropertyFeatures"]`).MustElements("a")
	est := page.MustElement("[aria-label='Property value']").MustText()
	// page.MustElement("[aria-label='Property timeline']")
	feat := page.MustElement("[aria-label='Property features']").MustText()
	zones := page.MustElement("[aria-label='Government planning overlays & zones']").MustText()

	// return fmt.Sprintf("Estimated: %s, Features: %s, Gov overlays: %s", est, feat, zones)
	// log.Printf("## house key info %s", fmt.Sprintf("Estimated: %s, Features: %s", est, feat))
	return fmt.Sprintf("Estimated: %s, Features: %s, Zones: %s", est, feat, zones)
}

func onExecSubStreet(subName string, street string, lots []string) string {
	if len(lots) == 0 {
		log.Fatalln("need lots for the street: " + street)
		return ""
	}

	path, _ := launcher.LookPath()

	// u := "ws://127.0.0.1:9222/devtools/browser/fdcdabc6-0c90-48da-ab98-86824158bb4d"
	log.Printf("## browser path %s", path)
	// browser := rod.New().DefaultDevice(MyDevice).ControlURL(u).MustConnect()
	wsurl := launcher.NewUserMode().Bin(path).MustLaunch()
	browser := rod.New().ControlURL(wsurl).MustConnect()
	// time.Sleep(3 * time.Second)
	defer browser.Close()

	subProfile := housing.NewSuburb(subName)
	url := subProfile.ToPropertyStreetUrl(street)
	log.Printf("## visiting street page %s", url)

	time.Sleep(2 * time.Second)
	page := browser.MustPage(url).MustWaitLoad()

	log.Printf("## street page %s loaded", url)
	// page.MustScreenshotFullPage("./screenshots/street.png")
	log.Printf("## street page loaded, waiting 3 secs to continue ")
	time.Sleep(3 * time.Second)

	links := page.MustElements("a")
	// links := page.MustWaitElementsMoreThan("a", 3).MustElements("a")
	log.Printf("links found %d", len(links))

	res := map[string]string{}
	// lotLinks := []string{}
	for _, link := range links {
		href := link.MustAttribute("href")

		if href == nil {
			log.Printf("## found invalid link ")
			continue
		}
		// only visiting house number we're interested in
		matched := slices.ContainsFunc(lots, func(num string) bool {
			return strings.Contains(*href, num+"-pid")
		})
		if !matched {
			continue
		}

		houseKeyInfo := collectHouseProfile(browser, subProfile.ToPropertyHouseUrl(*href))
		res[subName+"/"+street] = houseKeyInfo
		res[*href] = houseKeyInfo

		// sleepDuration := time.Duration(rand.Int()) * time.Second
		// time.Sleep(sleepDuration)
		time.Sleep(1 * time.Second)
	}

	// ==
	// utils.Pause()
	d, _ := json.Marshal(res)
	return string(d)
}

func onExec(input []string) string {
	res := make(map[string]string)
	for _, subname := range input {
		suburb := housing.NewSuburb(subname)
		content := collectMktInsights(suburb)
		extra := collectExtra(suburb)

		// spd := collectSupplyDemand(browser, suburb)
		res[subname] = content + "extra demographics:" + extra

	}
	d, _ := json.Marshal(res)
	return string(d)
}

func NewSuburbCmd() *cobra.Command {
	var suburbs, street, lotsNum string

	// Usage: ./prober suburb -n daisy-hill-qld-4127 -s gladewood-dr -l (--lots) 45,48
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
			lots := strings.Split(lotsNum, ",")
			if street != "" && lotsNum != "" {
				log.Printf("## start exec suburb street probing %+v, lots: %s", street, lots)
				output := onExecSubStreet(suburbList[0], street, lots)

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
	suburbCmd.Flags().StringVarP(&street, "street", "s", "", "Street name of suburb")
	suburbCmd.Flags().StringVarP(&lotsNum, "lots", "l", "", "House lot number on the Street of a suburb")

	return suburbCmd
}
