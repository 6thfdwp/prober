package housing

import "strings"

type SuburbProfileProber interface {
	DoMarketInsights()
}

// type suburbRodProber struct{}

type SuburbStreetProber interface {
}

type SuburbProfile struct {
	State    string
	Postcode string
	Name     string

	Insights     string
	Demographics string
}

const (
	PropertySite = "https://www.property.com.au"
	REASite      = "https://www.realestate.com.au"
	DomainSite   = "https://www.domain.com.au"
)

// n has format daisy-hill-qld-4127
func NewSuburb(n string) SuburbProfile {
	parts := strings.Split(n, "-")
	l := len(parts)
	state, postcode := parts[l-2], parts[l-1]

	return SuburbProfile{
		State: state, Postcode: postcode,
		Name: strings.Join(parts[:l-2], "-"),
	}
}

func (t *SuburbProfile) ToDmainFullUrl() string {
	// daisy-hill-qld-4127
	return DomainSite + "/" + t.Name + "-" + t.State + "-" + t.Postcode
}

func (t *SuburbProfile) ToPropertyStreetUrl(street string) string {
	return PropertySite + "/" + t.State + "/" + t.Name + "-" + t.Postcode + "/" + street
}
func (t *SuburbProfile) ToPropertyHouseUrl(streetlot string) string {
	return PropertySite + streetlot
}
func (t *SuburbProfile) ToREAFullUrl() string {
	//
	return REASite + "/" + t.State + "/" + t.Postcode
}

func (t *SuburbProfile) Stringify() {
	//
}

func (t *SuburbProfile) DoMktInsights(sub string)       {}
func (t *SuburbProfile) CollectSupplyDemand(sub string) {}
