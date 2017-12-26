package main

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"sort"
	gc "github.com/Regela/goncurses"
	"log"
	"strconv"
)

var scores []int
var scorescr *gc.Window

var CurScore int

func AddScore(score int) {
	scores=append(scores,score)
	sort.Ints(scores)
}

func CurScoreInc(){
	CurScore++
	MoveAddString(scorescr,2,1,strconv.FormatInt(int64(CurScore), 16))
	scorescr.NoutRefresh()
}

func InitScores(){
	LoadScores()
	CurScore=0
	var err error
	scorescr, err = gc.NewWindow(int(MAX_Y)+2,10, 0, int(MAX_X)*2+3)
	if err != nil {
		gc.End()
		log.Fatal(err)
	}

	scorescr.Box('|','~' )
	MoveAddString(scorescr,1,1,"CurScr")
	MoveAddString(scorescr,3,1,"HiScrs")
	for i:=0 ; i <  10 && i < len(scores); i ++{
		MoveAddString(scorescr,i+4,1,strconv.FormatInt(int64(scores[len(scores)-1-i]),16))
	}
	scorescr.NoutRefresh()
	gc.Update()
}

func AddCurScoreAndSave(){
	AddScore(CurScore)
	SaveScores()
}

func LoadScores(){
	if _, err := os.Stat(Usr.HomeDir+"/.config/GoSnakeScores.json"); os.IsNotExist(err) {
		scores=make([]int,1)
		SaveScores()
	}
	data,_ := ioutil.ReadFile(Usr.HomeDir+"/.config/GoSnakeScores.json")
	json.Unmarshal(data, &scores)
}

func SaveScores(){
	data,_ := json.Marshal(scores)
	err := ioutil.WriteFile(Usr.HomeDir+"/.config/GoSnakeScores.json",data,0755)
	if err != nil {log.Fatal(err)}
}