package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
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

func fetchHTML(url string) (*html.Node, error) {
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return html.Parse(strings.NewReader(string(body)))
}

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += extractText(c)
	}
	return strings.TrimSpace(text)
}

func extractAttribute(n *html.Node, attrName string) string {
	for _, attr := range n.Attr {
		if attr.Key == attrName {
			return attr.Val
		}
	}
	return ""
}

func findNodeByAttr(n *html.Node, attrName, attrValue string) *html.Node {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == attrName && attr.Val == attrValue {
				return n
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findNodeByAttr(c, attrName, attrValue); result != nil {
			return result
		}
	}
	return nil
}

func findNodeByTag(n *html.Node, tag atom.Atom) *html.Node {
	if n.DataAtom == tag {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findNodeByTag(c, tag); result != nil {
			return result
		}
	}
	return nil
}

func findNodesByTag(n *html.Node, tag atom.Atom) []*html.Node {
	var nodes []*html.Node
	if n.DataAtom == tag {
		nodes = append(nodes, n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, findNodesByTag(c, tag)...)
	}
	return nodes
}

func main() {
	pokedexURL := "https://pokedex.org/"
	bulbapediaURL := "https://bulbapedia.bulbagarden.net/wiki/List_of_Pok%C3%A9mon_by_effort_value_yield_(Generation_IX)"
	const limit = 5

	root, err := fetchHTML(pokedexURL)
	if err != nil {
		log.Fatal(err)
	}

	var Pokedex []Pokemon
	monstersList := findNodeByAttr(root, "id", "monsters-list")
	if monstersList != nil {
		liNodes := findNodesByTag(monstersList, atom.Li)
		for i, li := range liNodes {
			if i >= limit {
				break
			}
			pokemon := Pokemon{}
			span := findNodeByTag(li, atom.Span)
			if span != nil {
				pokemon.Name = extractText(span)
			}

			button := findNodeByTag(li, atom.Button)
			if button != nil {
				style := extractAttribute(button, "style")
				re := regexp.MustCompile(`background-image:\s*url\(([^)]+)\)`)
				match := re.FindStringSubmatch(style)
				if len(match) > 1 {
					pokemon.Image = match[1]
				}
			}

			detailInfobox := findNodeByAttr(root, "class", "detail-infobox")
			if detailInfobox != nil {
				generalInfo := General{}
				types := findNodeByAttr(detailInfobox, "class", "detail-types")
				if types != nil {
					for _, span := range findNodesByTag(types, atom.Span) {
						generalInfo.Type = append(generalInfo.Type, extractText(span))
					}
				}

				nationalID := findNodeByAttr(detailInfobox, "class", "detail-national-id")
				if nationalID != nil {
					generalInfo.ID = extractText(findNodeByTag(nationalID, atom.Span))
				}

				stats := findNodeByAttr(detailInfobox, "class", "detail-stats")
				if stats != nil {
					for _, statRow := range findNodesByTag(stats, atom.Div) {
						if extractAttribute(statRow, "class") == "detail-stats-row" {
							stat := Stat{}
							stat.StatName = extractText(statRow.FirstChild)
							stat.StatNum = extractText(statRow.LastChild)
							generalInfo.Stat = append(generalInfo.Stat, stat)
						}
					}
				}

				species := findNodeByAttr(root, "class", "monster-species")
				if species != nil {
					generalInfo.Species = extractText(species)
				}

				desc := findNodeByAttr(root, "class", "monster-description")
				if desc != nil {
					generalInfo.Desc = extractText(desc)
				}
				pokemon.GeneralInfo = generalInfo

				// Profile
				detail := findNodeByAttr(root, "class", "detail-below-header")
				if detail != nil {
					profile := Profile{}
					for _, strong := range findNodesByTag(detail, atom.Strong) {
						switch strings.TrimSpace(extractText(strong)) {
						case "Height:":
							profile.Height = extractText(strong.NextSibling)
						case "Weight:":
							profile.Weight = extractText(strong.NextSibling)
						case "Catch Rate:":
							profile.CatchRate = extractText(strong.NextSibling)
						case "Gender Ratio:":
							profile.GenderRatio = extractText(strong.NextSibling)
						case "Egg Groups:":
							profile.EggGroups = extractText(strong.NextSibling)
						case "Hatch Steps:":
							profile.HatchSteps = extractText(strong.NextSibling)
						case "Abilities:":
							profile.Abilities = extractText(strong.NextSibling)
						case "EVs:":
							profile.EVs = extractText(strong.NextSibling)
						}
					}
					pokemon.Profile = profile
				}

			}

			Pokedex = append(Pokedex, pokemon)
		}
	}

	// Fetch and parse additional data from Bulbapedia
	bulbapediaRoot, err := fetchHTML(bulbapediaURL)
	if err != nil {
		log.Fatal(err)
	}

	XP := findNodeByAttr(bulbapediaRoot, "class", "sortable.roundy.jquery-tablesorter")
	if XP != nil {
		tbody := findNodeByTag(XP, atom.Tbody)
		if tbody != nil {
			for i, row := range findNodesByTag(tbody, atom.Tr) {
				if i >= limit {
					break
				}
				cols := findNodesByTag(row, atom.Td)
				if len(cols) >= 4 {
					baseXp := extractText(cols[3])
					Pokedex[i].BaseXp = baseXp
				}
			}
		}
	}

	js, err := json.MarshalIndent(Pokedex, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Writing data to file")
	if err := os.WriteFile("pokedex.json", js, 0664); err == nil {
		fmt.Println("Data written to file successfully")
	}
}
