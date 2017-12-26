package main


import (
	gc "github.com/Regela/goncurses"
	"github.com/famz/SetLocale"
	"log"
	"time"
	"math/rand"
	"os"
	"os/user"
)
var pause bool
var rnd *rand.Rand
var stdscr,gamescr,framescr *gc.Window
var rows, cols int

var MAX_X uint8
var MAX_Y uint8
var Usr *user.User
const(
	u = iota
	d
	l
	r

)
var dir,lastdir uint8
type point struct{
	X uint8
	Y uint8
}

type chpoint struct{
	P point
	Ch gc.Char
}

var	feed point
func FeedGenerate()  {
genfound:
	feed.X=uint8(rnd.Intn( int(MAX_X)))
	feed.Y=uint8(rnd.Intn( int(MAX_Y)))
	for i := range snake{
		if snake[i].P == feed || snakeHead.P == feed {goto genfound}//Перегениерация, если попало на змейку
	}
}
func ChangeDir(){
	for {
		c:=framescr.GetChar()
		if c == 'p'{
			pause = ! pause
			continue
		}
		switch c {
		case 'w':
			{
				if dir==u{
					move()
				}else if lastdir!=d {dir = u}
				break
			}
		case 's':
			{
				if dir==d{
					move()
				}else if lastdir!=u {dir = d}
				break
			}
		case 'd':
			{
				if dir==r{
					move()
				}else if lastdir!=l {dir = r}
				break
			}
		case 'a':
			{
				if dir==l{
					move()
				}else if lastdir!=r {dir = l}
				break
			}

		case 'q':
			{
				Exit()
			}

		}


	}
}
var snakeHead = chpoint{point{0,0},'|'}
var snake []chpoint
func moveTime(){
	for {
		move()

		time.Sleep(time.Second / 3)
	}
}

func MoveAddString(w *gc.Window,y int,x int, str string){
	w.Move(y,x)
	for i:=range str{
		w.AddChar(gc.Char(str[i]))
	}
}

func deadSnake(){
	my,mx:=framescr.MaxYX()
	MoveAddString(framescr,my/2,mx/2-6,"Game Over")
	MoveAddString(framescr,my/2+1,mx/2-6,"Goodbye")
	gamescr.NoutRefresh()
	framescr.NoutRefresh()
	gc.Update()
	time.Sleep(time.Second*2)
	Exit()
}

func move(){
	//f:=false
	//oldSnakeHead:=snakeHead
	if pause {return}
	if dir == u {
		snakeHead.P.Y--
		snakeHead.Ch='|'
		if lastdir == l{
			snake[0].Ch='\\'
		} else if lastdir == r{
			snake[0].Ch='/'
		}
	} else if dir == d {
		snakeHead.P.Y++
		snakeHead.Ch='|'
		if lastdir == r{
			snake[0].Ch='\\'
		} else if lastdir == l{
			snake[0].Ch='/'
		}
	} else if dir == l {
		snakeHead.P.X--
		snakeHead.Ch='='
		if lastdir == u{
			snake[0].Ch='\\'
		} else if lastdir == d{
			snake[0].Ch='/'
		}
	} else if dir == r {
		snakeHead.P.X++
		snakeHead.Ch='='
		if lastdir == d{
			snake[0].Ch='\\'
		} else if lastdir == u{
			snake[0].Ch='/'
		}
	}

	lastdir=dir
	if snakeHead.P.Y>=MAX_Y || snakeHead.P.X>=MAX_X{//Если врезался в стенку
		//if snakeHead.X==255{
		//	framescr.MoveAddChar(int(snakeHead.Y+1),int(0),'X')
		//	framescr.NoutRefresh()
		//	gamescr.MoveAddChar(int(snake[0].Y),int(snake[0].X),'0')
		//}else if snakeHead.Y==255{
		//	framescr.MoveAddChar(int(0),int(snakeHead.X+1),'X')
		//	framescr.NoutRefresh()
		//	gamescr.MoveAddChar(int(snake[0].Y),int(snake[0].X),'0')
		//}else if snakeHead.X>=MAX_X{
		//	framescr.MoveAddChar(int(snakeHead.Y+1),int(MAX_X+1),'X')
		//	framescr.NoutRefresh()
		//	gamescr.MoveAddChar(int(snake[0].Y),int(snake[0].X),'0')
		//}else if snakeHead.Y>=MAX_Y{
		//	framescr.MoveAddChar(int(MAX_Y+1),int(snakeHead.X+1),'X')
		//	framescr.NoutRefresh()
		//	gamescr.MoveAddChar(int(snake[0].Y),int(snake[0].X),'0')
		//}
		deadSnake()
		//snakeHead=oldSnakeHead
		//return
	}
	if snakeHead.P==feed{
		FeedGenerate()
		CurScoreInc()
		//f=true
		snake=append(snake,snake[len(snake)-1])
	}
	//for i:=len(snake)-1; i >0; i--{
	//	if snakeHead==snake[i-1] {
	//		snakeHead=oldSnakeHead
	//		return
	//	}
	//}
	for i:=len(snake)-1; i >0; i--{
		if snakeHead.P==snake[i-1].P {
			deadSnake()
			//snakeHead=oldSnakeHead
			//return
		}
		snake[i]=snake[i-1]
	}
	snake[0]=snakeHead


	gamescr.Erase()
	for i:=len(snake)-1; i >0; i--{
		//gamescr.MoveAddChar(int(snake[i].Y),int(snake[i].X),'0')
		gamescr.MoveAddChar(int(snake[i].P.Y),int(snake[i].P.X*2),snake[i].Ch|gc.ColorPair(1))
		gamescr.MoveAddChar(int(snake[i].P.Y),int(snake[i].P.X*2+1),snake[i].Ch|gc.ColorPair(1))

	}
	gamescr.MoveAddChar(int(snakeHead.P.Y),int(snakeHead.P.X*2),gc.ACS_BLOCK|gc.ColorPair(1))
	gamescr.MoveAddChar(int(snakeHead.P.Y),int(snakeHead.P.X*2+1),gc.ACS_BLOCK|gc.ColorPair(1))

	gamescr.MoveAddChar(int(feed.Y),int(feed.X*2),'$')
	gamescr.MoveAddChar(int(feed.Y),int(feed.X*2+1),'$')
	gamescr.NoutRefresh()


	gc.Update()


}

func initPairs()  {
	gc.InitColor(1,255,0,0)
	gc.InitColor(2,0,255,0)

	gc.InitPair(1,2,1)


}

func main(){
	MAX_X=15
	MAX_Y=15
	snake=make([]chpoint,3)
	for i:=len(snake)-1; i >=0; i-- {
		snake[i]=snakeHead
	}
	dir=r
	lastdir=dir
	pause=false
	SetLocale.SetLocale(SetLocale.LC_ALL, "")
	rnd=rand.New(rand.NewSource(time.Now().UnixNano()))
	var err error
	stdscr, err = gc.Init()
	gc.StartColor()
	initPairs()
	if err != nil {
		gc.End()
		log.Fatal(err)
	}
	gc.Echo(false)
	gc.CBreak(true)
	gc.Cursor( 0)
	rows, cols = stdscr.MaxYX()
	if rows < int(MAX_Y+2) || cols < int(MAX_X+2)  {
		gc.End()
		log.Fatal("Too small")
	}
	framescr, err = gc.NewWindow(int(MAX_Y+2),int(MAX_X*2+2), 0, 0)
	if err != nil {

		gc.End()
		log.Fatal(err)
	}
	framescr.Box(gc.ACS_BOARD, gc.ACS_BOARD)
	framescr.MovePrint(1, 0, "")
	framescr.NoutRefresh()
	gc.Update()
	Usr, err = user.Current()
	if err != nil {
		log.Fatal( err )
	}


	InitScores()
	//time.Sleep(time.Second)
	gamescr, err = gc.NewWindow(int(MAX_Y),int(MAX_X*2), 1, 1)
	if err != nil {

		gc.End()
		log.Fatal(err)
	}
	gamescr.NoutRefresh()

	FeedGenerate()
	gc.Update()

	//
	go ChangeDir()

	time.Sleep(time.Second)
	moveTime()
	gc.End()

}

func Exit()  {
	gc.End()
	AddCurScoreAndSave()
	os.Exit(0)
}