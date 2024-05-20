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

	c.OnHTML("div-monsters-list-wrapper > ul-monsters-list", func(e *colly.HTMLElement) {
		counter := 0
		const limit = 10
		e.ForEach("li", func(i int, h *colly.HTMLElement) {
			if counter >= limit {
				return
			}

			Pokemon := Pokemon{}
			Pokemon.Name = h.ChildText("span")
			c.OnHTML("button[type=\"button\"]", func(e *colly.HTMLElement) {
				style := e.Attr("style")
				re := regexp.MustCompile(`background-image:\s*url\(([^)]+)\)`)
				match := re.FindStringSubmatch(style)
				if len(match) > 1 {
					Pokemon.Image = match[1]
				}
			})
			GeneralInfo := General{}
			c.OnHTML("div.detail-infobox > div.detail-types-and-num > div.detail-types", func(e *colly.HTMLElement) {
				e.ForEach("span", func(_ int, h *colly.HTMLElement) {
					GenType := h.Text
					GeneralInfo.Type = append(GeneralInfo.Type, GenType)
				})
			})
			c.OnHTML("div.detail-infobox > div.detail-types-and-num > div.detail-national-id", func(e *colly.HTMLElement) {
				GeneralInfo.ID = e.ChildText("span")
			})

			c.OnHTML("div.detail-infobox > div.detail-stats", func(e *colly.HTMLElement) {
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

			c.OnHTML("div.monster-minutia > strong", func(strong *colly.HTMLElement) {
				if strings.TrimSpace(strong.Text) == "Height:" {
					span := strong.DOM.Next()
					Profile.Height = span.Text()
				}
			})

			c.OnHTML("div.monster-minutia > strong", func(strong *colly.HTMLElement) {
				if strings.TrimSpace(strong.Text) == "Weight:" {
					span := strong.DOM.Next()
					Profile.Weight = span.Text()
				}
			})

			c.OnHTML("div.monster-minutia > strong", func(strong *colly.HTMLElement) {
				if strings.TrimSpace(strong.Text) == "Catch Rate:" {
					span := strong.DOM.Next()
					Profile.CatchRate = span.Text()
				}
			})

			c.OnHTML("div.monster-minutia > strong", func(strong *colly.HTMLElement) {
				if strings.TrimSpace(strong.Text) == "Gender Ratio:" {
					span := strong.DOM.Next()
					Profile.GenderRatio = span.Text()
				}
			})

			c.OnHTML("div.monster-minutia > strong", func(strong *colly.HTMLElement) {
				if strings.TrimSpace(strong.Text) == "Egg Groups:" {
					span := strong.DOM.Next()
					Profile.EggGroups = span.Text()
				}
			})

			c.OnHTML("div.monster-minutia > strong", func(strong *colly.HTMLElement) {
				if strings.TrimSpace(strong.Text) == "Hatch Steps:" {
					span := strong.DOM.Next()
					Profile.HatchSteps = span.Text()
				}
			})

			c.OnHTML("div.monster-minutia > strong", func(strong *colly.HTMLElement) {
				if strings.TrimSpace(strong.Text) == "Abilities:" {
					span := strong.DOM.Next()
					Profile.Abilities = span.Text()
				}
			})

			c.OnHTML("div.monster-minutia > strong", func(strong *colly.HTMLElement) {
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

			c.OnHTML("div.evolutions", func(e *colly.HTMLElement) {
				e.ForEach("div.evolution-row", func(_ int, l *colly.HTMLElement) {
					c.OnHTML("div.evolution-label", func(h *colly.HTMLElement) {
						Evoinfo := h.ChildText("span")
						Pokemon.Evolution = append(Pokemon.Evolution, Evoinfo)
					})
				})
			})

			//Moves
			c.OnHTML("div.monster-moves", func(e *colly.HTMLElement) {
				counter1 := 0
				const limit = 10
				e.ForEach("div.moves-row", func(_ int, l *colly.HTMLElement) {
					if counter1 >= limit {
						return
					}

					Moves := Moves{}
					c.OnHTML("div.moves-inner-row", func(h *colly.HTMLElement) {
						var MoveNum, MoveName, MoveType string
						e.ForEach("span", func(i int, l *colly.HTMLElement) {
							text := strings.TrimSpace(l.Text)
							switch i {
							case 0:
								MoveNum = text
							case 1:
								MoveName = text
							case 2:
								MoveType = text
							}
						})
						Moves.MoveNum = MoveNum
						Moves.MoveName = MoveName
						Moves.MoveType = MoveType
					})

					c.OnHTML("div.moves-row-detail > div.moves-row-stats", func(h *colly.HTMLElement) {
						c.OnHTML("strong", func(strong *colly.HTMLElement) {
							textList := []string{"Power:", "Acc:", "PP:"}
							strongText := strings.TrimSpace(strong.Text)
							for _, targetText := range textList {
								if strongText == targetText {
									span := strong.DOM.Next()
									spanText := strings.TrimSpace(span.Text())
									switch strongText {
									case "Power:":
										Moves.Power = spanText
									case "Acc:":
										Moves.Acc = spanText
									case "PP:":
										Moves.PP = spanText
									}
									break
								}
							}
						})
					})

					c.OnHTML("div.moves-row-detail", func(h *colly.HTMLElement) {
						c.OnHTML("move-description", func(l *colly.HTMLElement) {
							Moves.MoveDesc = strings.TrimSpace(l.Text)
						})
					})

					Pokemon.Move = append(Pokemon.Move, Moves)
				})

				counter1++
			})

			counter++
		})

	})

	d.OnHTML("div-mw-content-text > div.mw-parser-output > table.sortable roundy jquery-tablesorter > tbody", func(e *colly.HTMLElement) {
		counter := 0
		const limit = 10
		Pokemon := Pokemon{}
		e.ForEach("tr", func(i int, h *colly.HTMLElement) {
			if counter >= limit {
				return
			}

			tds := e.DOM.Find("td")
			if tds.Length() >= 4 {
				Pokemon.BaseXp = tds.Eq(3).Text()
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

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})

	go func() {
		c.Visit("https://pokedex.org/")
	}()

	go func() {
		d.Visit("https://bulbapedia.bulbagarden.net/wiki/List_of_Pok%C3%A9mon_by_effort_value_yield_(Generation_IX)")
	}()

	c.Wait()
	d.Wait()

	var fin1 = 0
	var fin2 = 0

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
		fin1 = 1
	})
	d.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
		fin2 = 1
	})

	if fin1 == 1 && fin2 == 1 {
		js, err := json.MarshalIndent(Pokedex, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Writing data to file")
		if err := os.WriteFile("pokedex.json", js, 0664); err == nil {
			fmt.Println("Data written to file successfully")
		}
	}
}
