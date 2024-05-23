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

func extractBackgroundImageURL(style string) string {
	re := regexp.MustCompile(`background-image:\s*url\(([^)]+)\)`)
	match := re.FindStringSubmatch(style)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func findNodeByAttr(n *html.Node, attrName, attrValue string) *html.Node {
	if n == nil {
		return nil
	}
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == attrName && strings.Contains(attr.Val, attrValue) {
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
	if n == nil {
		return nil
	}
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
	bulbapediaURL := "https://bulbapedia.bulbagarden.net/wiki/List_of_Pok%C3%A9mon_by_effort_value_yield_(Generation_IX)"
	const limit = 5
	var i int
	var Pokedex []Pokemon
	for i = 1; i <= 5; i++ {
		pokedexURL := fmt.Sprintf("https://pokedex.org/#/pokemon/%d", i)
		root, err := fetchHTML(pokedexURL)
		if err != nil {
			log.Fatal(err)
		}

		monstersList := findNodeByAttr(root, "class", "mui-panel")
		if monstersList != nil {
			pokemon := Pokemon{}
			h1 := findNodeByTag(monstersList, atom.H1)
			if h1 != nil {
				pokemon.Name = extractText(h1)
			}

			divImg := findNodeByAttr(monstersList, "class", "detail-sprite")
			if divImg != nil {
				style := extractAttribute(divImg, "style")
				pokemon.Image = extractBackgroundImageURL(style)
			}

			content := findNodeByAttr(root, "class", "detail-panel-content")
			if content != nil {
				detailInfobox := findNodeByAttr(content, "class", "detail-infobox")
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
			}
			var multFunc func(*html.Node)
			multFunc = func(n *html.Node) {
				if n.Type == html.ElementNode && n.Data == "div" && extractAttribute(n, "class") == "when-attacked-row" {
					var pokeType, pokeMult string
					for c := n.FirstChild; c != nil; c = c.NextSibling {
						if c.Type == html.ElementNode && c.Data == "span" {
							if class := extractAttribute(c, "class"); class == "monster-type" {
								pokeType = strings.TrimSpace(c.FirstChild.Data)
							} else if class == "monster-multiplier" {
								pokeMult = strings.TrimSpace(c.FirstChild.Data)

								if pokeType != "" && pokeMult != "" {
									Multiplier := Multiplier{
										DamageType: pokeType,
										DamageMult: pokeMult,
									}
									pokemon.DamageMultiplier = append(pokemon.DamageMultiplier, Multiplier)
								}
							}
						}
					}
				}

				for c := n.FirstChild; c != nil; c = c.NextSibling {
					multFunc(c)
				}
			}
			multFunc(root)

			var evoFunc func(*html.Node)
			evoFunc = func(n *html.Node) {
				if n.Type == html.ElementNode && n.Data == "div" && extractAttribute(n, "class") == "evolution-row" {
					for c := n.FirstChild; c != nil; c = c.NextSibling {
						if c.Type == html.ElementNode && c.Data == "div" && extractAttribute(c, "class") == "evolution-label" {
							evolutionInfo := strings.TrimSpace(c.FirstChild.FirstChild.Data)
							pokemon.Evolution = append(pokemon.Evolution, evolutionInfo)
						}
					}
				}

				for c := n.FirstChild; c != nil; c = c.NextSibling {
					evoFunc(c)
				}
			}
			evoFunc(root)

			var currentMove *Moves

			var parseNode func(*html.Node)
			parseNode = func(n *html.Node) {
				if n.Type == html.ElementNode {
					if n.Data == "div" && len(n.Attr) > 0 && n.Attr[0].Key == "class" {
						switch n.Attr[0].Val {
						case "moves-inner-row":
							currentMove = &Moves{}
							pokemon.Move = append(pokemon.Move, *currentMove)
						case "moves-row-stats":
							for c := n.FirstChild; c != nil; c = c.NextSibling {
								if c.Type == html.ElementNode && c.Data == "strong" {
									switch strings.TrimSpace(c.FirstChild.Data) {
									case "Power:":
										currentMove.Power = strings.TrimSpace(c.LastChild.Data)
									case "Acc:":
										currentMove.Acc = strings.TrimSpace(c.LastChild.Data)
									case "PP:":
										currentMove.PP = strings.TrimSpace(c.LastChild.Data)
									}
								}
							}
						case "move-description":
							currentMove.MoveDesc = strings.TrimSpace(n.FirstChild.Data)
						}
					}
				}
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					parseNode(c)
				}
			}

			parseNode(root)

			// Print parsed moves
			for _, move := range pokemon.Move {
				fmt.Printf("MoveName: %s, Power: %s, Acc: %s, PP: %s, Description: %s\n", move.MoveName, move.Power, move.Acc, move.PP, move.MoveDesc)
			}

			Pokedex = append(Pokedex, pokemon)
		}

	}

	// Fetch and parse additional data from Bulbapedia
	bulbapediaRoot, err := fetchHTML(bulbapediaURL)
	if err != nil {
		log.Fatal(err)
	}

	table := findNodeByAttr(bulbapediaRoot, "class", "jquery-tablesorter")

	tbody := findNodeByTag(table, atom.Tbody)
	if tbody != nil {
		rows := findNodesByTag(tbody, atom.Tr)
		for i, row := range rows {
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

	js, err := json.MarshalIndent(Pokedex, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Writing data to file")
	if err := os.WriteFile("pokedex.json", js, 0664); err == nil {
		fmt.Println("Data written to file successfully")
	}
}
