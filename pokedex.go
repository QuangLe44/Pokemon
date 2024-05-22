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

type Moves struct {
	MoveNum  string
	MoveName string
	MoveType string
	Power    string
	Acc      string
	PP       string
	MoveDesc string
}

type Pokemon struct {
	Name             string
	Image            string
	GeneralInfo      General
	Profile          Profile
	DamageMultiplier []Multiplier
	Evolution        []string
	Move             []Moves
	BaseXp           string
}

func main() {
	c := colly.NewCollector()
	c.SetRequestTimeout(120 * time.Second)
	d := colly.NewCollector()
	d.SetRequestTimeout(120 * time.Second)

	Pokedex := make([]Pokemon, 0)

	c.OnHTML("div.monsters-list-wrapper > ul.monsters-list", func(e *colly.HTMLElement) {
		counter := 0
		const limit = 5
		e.ForEach("li", func(i int, h *colly.HTMLElement) {
			if counter >= limit {
				return
			}
			Pokemon := Pokemon{}
			Pokemon.Name = h.ChildText("span")

			// Nested collector to get details of each Pokemon
			detailCollector := c.Clone()
			detailCollector.OnHTML("div.monster-details", func(e *colly.HTMLElement) {
				// Image
				style := e.ChildAttr("button[type=\"button\"]", "style")
				re := regexp.MustCompile(`background-image:\s*url\(([^)]+)\)`)
				match := re.FindStringSubmatch(style)
				if len(match) > 1 {
					Pokemon.Image = match[1]
				}

				// General Info
				GeneralInfo := General{}
				e.ForEach("div.detail-types-and-num > div.detail-types > span", func(_ int, h *colly.HTMLElement) {
					GenType := h.Text
					GeneralInfo.Type = append(GeneralInfo.Type, GenType)
				})
				GeneralInfo.ID = e.ChildText("div.detail-types-and-num > div.detail-national-id > span")
				e.ForEach("div.detail-stats > div.detail-stats-row", func(_ int, h *colly.HTMLElement) {
					Status := Stat{}
					Status.StatName = h.ChildText("span:not([class])")
					Status.StatNum = h.ChildText("span.stat-bar")
					GeneralInfo.Stat = append(GeneralInfo.Stat, Status)
				})
				GeneralInfo.Species = e.ChildText("div.monster-species")
				GeneralInfo.Desc = e.ChildText("div.monster-description")
				Pokemon.GeneralInfo = GeneralInfo

				// Profile
				Profile := Profile{}
				e.ForEach("div.monster-minutia > strong", func(_ int, strong *colly.HTMLElement) {
					span := strong.DOM.Next()
					switch strings.TrimSpace(strong.Text) {
					case "Height:":
						Profile.Height = span.Text()
					case "Weight:":
						Profile.Weight = span.Text()
					case "Catch Rate:":
						Profile.CatchRate = span.Text()
					case "Gender Ratio:":
						Profile.GenderRatio = span.Text()
					case "Egg Groups:":
						Profile.EggGroups = span.Text()
					case "Hatch Steps:":
						Profile.HatchSteps = span.Text()
					case "Abilities:":
						Profile.Abilities = span.Text()
					case "EVs:":
						Profile.EVs = span.Text()
					}
				})
				Pokemon.Profile = Profile

				// Damage Multiplier
				e.ForEach("div.when-attacked > div.when-attacked-row", func(_ int, h *colly.HTMLElement) {
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

				// Evolution
				e.ForEach("div.evolutions > div.evolution-row", func(_ int, l *colly.HTMLElement) {
					Evoinfo := l.ChildText("div.evolution-label > span")
					Pokemon.Evolution = append(Pokemon.Evolution, Evoinfo)
				})

				// Moves
				e.ForEach("div.monster-moves > div.moves-row", func(_ int, l *colly.HTMLElement) {
					Moves := Moves{}
					l.ForEach("div.moves-inner-row > span", func(i int, h *colly.HTMLElement) {
						text := strings.TrimSpace(h.Text)
						switch i {
						case 0:
							Moves.MoveNum = text
						case 1:
							Moves.MoveName = text
						case 2:
							Moves.MoveType = text
						}
					})
					l.ForEach("div.moves-row-detail > div.moves-row-stats > strong", func(_ int, strong *colly.HTMLElement) {
						span := strong.DOM.Next()
						switch strings.TrimSpace(strong.Text) {
						case "Power:":
							Moves.Power = span.Text()
						case "Acc:":
							Moves.Acc = span.Text()
						case "PP:":
							Moves.PP = span.Text()
						}
					})
					Moves.MoveDesc = l.ChildText("div.moves-row-detail > move-description")
					Pokemon.Move = append(Pokemon.Move, Moves)
				})

				Pokedex = append(Pokedex, Pokemon)
			})

			detailCollector.Visit(h.ChildAttr("a", "href"))

			counter++
		})
	})

	d.OnHTML("div.mw-content-text > div.mw-parser-output > table.sortable.roundy.jquery-tablesorter > tbody", func(e *colly.HTMLElement) {
		counter := 0
		const limit = 5
		e.ForEach("tr", func(i int, h *colly.HTMLElement) {
			if counter >= limit {
				return
			}
			tds := h.DOM.Find("td")
			if tds.Length() >= 4 {
				Pokedex[i].BaseXp = tds.Eq(3).Text()
			}
			counter++
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", r.Request.URL)
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Got this error:", err)
	})

	d.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	d.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", r.Request.URL)
	})
	d.OnError(func(r *colly.Response, err error) {
		fmt.Println("Got this error:", err)
	})

	c.Visit("https://pokedex.org/")
	d.Visit("https://bulbapedia.bulbagarden.net/wiki/List_of_Pok%C3%A9mon_by_effort_value_yield_(Generation_IX)")

	c.Wait()
	d.Wait()

	js, err := json.MarshalIndent(Pokedex, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Writing data to file")
	if err := os.WriteFile("pokedex.json", js, 0664); err == nil {
		fmt.Println("Data written to file successfully")
	}
}
