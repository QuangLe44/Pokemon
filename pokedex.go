package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Stat struct {
	StatName string
	StatNum  string
}

type General struct {
	ID      string
	Type    []string
	Stat    []Stat
	Species string
	Desc    string
}

type Profile struct {
	Height      string
	Weight      string
	CatchRate   string
	GenderRatio string
	EggGroups   string
	HatchSteps  string
	Abilities   string
	EVs         string
}

type Multiplier struct {
	DamageType string
	DamageMult string
}

type Evolution struct {
	EvoImage string
	EvoInfo  string
}

type Moves struct {
	MoveNum  string
	MoveName string
	MoveType string
	MoveStat []string
	MoveDesc string
}

type Pokemon struct {
	Name             string
	Image            string
	GeneralInfo      General
	Profile          Profile
	DamageMultiplier []Multiplier
	Evolution        []Evolution
	Move             []Moves
}

func main() {
	c := colly.NewCollector()
	c.SetRequestTimeout(120 * time.Second)
	Pokedex := make([]Pokemon, 0)

	c.OnHTML("div-monsters-list-wrapper ul-monsters-list", func(e *colly.HTMLElement) {
		e.ForEach("li", func(i int, h *colly.HTMLElement) {
			Pokemon := Pokemon{}
			Pokemon.Name = h.ChildText("span")
			c.OnHTML("button[type=\"button\"]", func(e *colly.HTMLElement) {
				style := e.Attr("style")
				re := regexp.MustCompile(`background-image:\s*url\(["']?([^"')]+)["']?\)`)
				match := re.FindStringSubmatch(style)
				if len(match) > 1 {
					Pokemon.Image = match[1]
				}
			})
			GeneralInfo := General{}
			c.OnHTML("div.detail-infobox div.detail-types-and-num div.detail-types", func(e *colly.HTMLElement) {
				e.ForEach("span", func(_ int, h *colly.HTMLElement) {
					GenType := h.Text
					GeneralInfo.Type = append(GeneralInfo.Type, GenType)
				})
			})
			c.OnHTML("div.detail-infobox div.detail-types-and-num div.detail-national-id", func(e *colly.HTMLElement) {
				GeneralInfo.ID = e.ChildText("span")
			})

			c.OnHTML("div.detail-infobox div.detail-stats", func(e *colly.HTMLElement) {
				e.ForEach("div.detail-stats-row", func(_ int, h *colly.HTMLElement) {
					Status := Stat{}
					Status.StatName = h.ChildText("span:not([class])")
					Status.StatNum = h.ChildText("span.stat-bar")
					GeneralInfo.Stat = append(GeneralInfo.Stat, Status)
				})
			})

			c.OnHTML("div.monster-species", func(e *colly.HTMLElement) {
				GeneralInfo.Species = e.Text
			})

			c.OnHTML("div.monster-description", func(e *colly.HTMLElement) {
				GeneralInfo.Desc = e.Text
			})

			Pokemon.GeneralInfo = GeneralInfo

			Profile := Profile{}

			c.OnHTML("div.monster-minutia strong", func(strong *colly.HTMLElement) {
				if strings.TrimSpace(strong.Text) == "Height:" {
					span := strong.DOM.Next()
					Profile.Height = span.Text()
				}
			})

			c.OnHTML("div.monster-minutia strong", func(strong *colly.HTMLElement) {
				if strings.TrimSpace(strong.Text) == "Weight:" {
					span := strong.DOM.Next()
					Profile.Weight = span.Text()
				}
			})

			c.OnHTML("div.monster-minutia strong", func(strong *colly.HTMLElement) {
				if strings.TrimSpace(strong.Text) == "Catch Rate:" {
					span := strong.DOM.Next()
					Profile.CatchRate = span.Text()
				}
			})

			c.OnHTML("div.monster-minutia strong", func(strong *colly.HTMLElement) {
				if strings.TrimSpace(strong.Text) == "Gender Ratio:" {
					span := strong.DOM.Next()
					Profile.GenderRatio = span.Text()
				}
			})

			c.OnHTML("div.monster-minutia strong", func(strong *colly.HTMLElement) {
				if strings.TrimSpace(strong.Text) == "Egg Groups:" {
					span := strong.DOM.Next()
					Profile.EggGroups = span.Text()
				}
			})

			c.OnHTML("div.monster-minutia strong", func(strong *colly.HTMLElement) {
				if strings.TrimSpace(strong.Text) == "Hatch Steps:" {
					span := strong.DOM.Next()
					Profile.HatchSteps = span.Text()
				}
			})

			c.OnHTML("div.monster-minutia strong", func(strong *colly.HTMLElement) {
				if strings.TrimSpace(strong.Text) == "Abilities:" {
					span := strong.DOM.Next()
					Profile.Abilities = span.Text()
				}
			})

			c.OnHTML("div.monster-minutia strong", func(strong *colly.HTMLElement) {
				if strings.TrimSpace(strong.Text) == "EVs:" {
					span := strong.DOM.Next()
					Profile.EVs = span.Text()
				}
			})

			Pokemon.Profile = Profile

			c.OnHTML("div.when-attacked", func(e *colly.HTMLElement) {
				e.ForEach("div.when-attacked-row", func(_ int, h *colly.HTMLElement) {
					var pokeType, pokeMult string
					h.ForEach("span", func(_ int, l *colly.HTMLElement) {
						if l.Attr("class") == "monster-type" {
							CurrType := l.Text
							if CurrType != "" {
								pokeType = CurrType
							}
						} else if l.Attr("class") == "monster-multiplier" {
							CurrMult := l.Text
							if CurrMult != "" {
								pokeMult = CurrMult
							}

							if pokeType != "" && pokeMult != "" {
								Multiplier := Multiplier{}
								Multiplier.DamageType = pokeType
								Multiplier.DamageMult = pokeMult
								Pokemon.DamageMultiplier = append(Pokemon.DamageMultiplier, Multiplier)
							}
						}
					})
				})
			})

			//evolution
		})

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", r.Request.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
		js, err := json.MarshalIndent(Pokedex, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Writing data to file")
		if err := os.WriteFile("pokedex.json", js, 0664); err == nil {
			fmt.Println("Data written to file successfully")
		}

	})

	c.Visit("")
}
