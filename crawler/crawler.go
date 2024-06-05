package main

import (
	"Pokemon/BaseInfo"
	"Pokemon/MonsterMoves"
	"Pokemon/MonsterType"
	"Pokemon/Stats"
	"Pokemon/description"
	"Pokemon/evolution"
	"Pokemon/move"
)

func main() {
	BaseInfo.Crawl()
	move.Crawl()
	MonsterMoves.Crawl()
	Stats.Crawl()
	evolution.Crawl()
	MonsterType.Crawl()
	description.Crawl()
}
