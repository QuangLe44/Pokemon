package main

import (
	"Pokemon/BaseInfo"
	"Pokemon/MonsterMoves"
	"Pokemon/MonsterType"
	"Pokemon/Stats"
	"Pokemon/description"
	"Pokemon/evolution"
	"Pokemon/exp"
	"Pokemon/move"
)

func main() {
	BaseInfo.Crawl()
	move.Crawl()
	MonsterMoves.Crawl()
	Stats.Crawl()
	evolution.Crawl()
	description.Crawl()
	MonsterType.Crawl()
	exp.Crawl()
	exp.Remove()
}
