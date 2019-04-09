package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/gdamore/tcell"
)

/*func MyHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}*/

type Pair struct {
	D   float64
	Dot float64
}

type ByD []Pair

func (bd ByD) Len() int {
	return len(bd)
}

func (bd ByD) Swap(i, j int) {
	bd[i], bd[j] = bd[j], bd[i]
}

func (bd ByD) Less(i, j int) bool {
	return bd[i].D < bd[j].D
}

const nScreenWidth = 120
const nScreenHeight = 40
const nMapWidth = 16
const nMapHeight = 16
const styleDefault tcell.Style = 0

var fPlayerX = float64(13.7)
var fPlayerY = float64(4.09)
var fPlayerA float64
var fFOV = float64(3.14159 / 4.0)
var fDepth = float64(16.0)
var fSpeed = float64(5.0)

func reverse(in string) string {
	parts := strings.Split(in, " ")
	var out string
	for index := len(parts) - 1; index >= 0; index-- {
		out += parts[index] + " "
	}
	out = strings.Trim(out, " ")
	return out
}

func pollKeyEvent(screen tcell.Screen, out chan<- rune) {
	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			{
				switch ev.Rune() {
				case 'a', 'A':
					{
						out <- 'A'
					}
				case 'd', 'D':
					{
						out <- 'D'
					}
				case 'w', 'W':
					{
						out <- 'W'
					}
				case 's', 'S':
					{
						out <- 'S'
					}
				}
			}

		}
	}
}

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Printf("Error creating screen: %s", err.Error())
		return
	}
	out := make(chan rune, 1000)
	err = screen.Init()
	if err != nil {
		fmt.Printf("Error initializing screen: %s", err.Error())
		return
	}
	go pollKeyEvent(screen, out)
	//var bytesWritten int
	screen.HideCursor()
	var ourMap string
	ourMap += "################"
	ourMap += "#..............#"
	ourMap += "#.......########"
	ourMap += "#..............#"
	ourMap += "#......##......#"
	ourMap += "#......##......#"
	ourMap += "#..............#"
	ourMap += "###............#"
	ourMap += "##.............#"
	ourMap += "#......####..###"
	ourMap += "#......#.......#"
	ourMap += "#......#.......#"
	ourMap += "#..............#"
	ourMap += "#......#########"
	ourMap += "#..............#"
	ourMap += "################"
	var t1 = time.Now()
	var t2 = time.Now()
	for {
		t2 = time.Now()
		ft := float64(t2.Sub(t1))
		t1 = t2
		fElapsedTime := float64(ft / float64(time.Second))

		fmt.Printf("fElapsedTime: %.10f\n", float64(fElapsedTime))
		select {
		case r := <-out:
			{
				switch r {
				case 'a', 'A':
					{
						fPlayerA -= (fSpeed * 0.75) * fElapsedTime
					}
				case 'd', 'D':
					{
						fPlayerA += (fSpeed * 0.75) * fElapsedTime
					}
				case 'w', 'W':
					{
						fPlayerX += math.Sin(fPlayerA) * fSpeed * fElapsedTime
						fPlayerY += math.Cos(fPlayerA) * fSpeed * fElapsedTime
						currentRune, _, _, _ := screen.GetContent(int(fPlayerX), int(fPlayerY))
						if currentRune == rune('#') {
							fPlayerX -= math.Sin(fPlayerA) * fSpeed * fElapsedTime
							fPlayerY -= math.Cos(fPlayerA) * fSpeed * fElapsedTime
						}
					}
				case 's', 'S':
					{
						fPlayerX -= math.Sin(fPlayerA) * fSpeed * fElapsedTime
						fPlayerY -= math.Cos(fPlayerA) * fSpeed * fElapsedTime
						currentRune, _, _, _ := screen.GetContent(int(fPlayerX), int(fPlayerY))
						if currentRune == rune('#') {
							fPlayerX += math.Sin(fPlayerA) * fSpeed * fElapsedTime
							fPlayerY += math.Cos(fPlayerA) * fSpeed * fElapsedTime
						}
					}
				}
			}
		default:
			{

			}
		}
		// get key states (key pressed)
		/*ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			{
				switch ev.Rune() {
				case 'a', 'A':
					{
						fPlayerA -= (fSpeed * 0.75) * fElapsedTime
					}
				case 'd', 'D':
					{
						fPlayerA += (fSpeed * 0.75) * fElapsedTime
					}
				case 'w', 'W':
					{
						fPlayerX += math.Sin(fPlayerA) * fSpeed * fElapsedTime
						fPlayerY += math.Cos(fPlayerA) * fSpeed * fElapsedTime
						currentRune, _, _, _ := screen.GetContent(int(fPlayerX), int(fPlayerY))
						if currentRune == rune('#') {
							fPlayerX -= math.Sin(fPlayerA) * fSpeed * fElapsedTime
							fPlayerY -= math.Cos(fPlayerA) * fSpeed * fElapsedTime
						}
					}
				case 's', 'S':
					{
						fPlayerX -= math.Sin(fPlayerA) * fSpeed * fElapsedTime
						fPlayerY -= math.Cos(fPlayerA) * fSpeed * fElapsedTime
						currentRune, _, _, _ := screen.GetContent(int(fPlayerX), int(fPlayerY))
						if currentRune == rune('#') {
							fPlayerX += math.Sin(fPlayerA) * fSpeed * fElapsedTime
							fPlayerY += math.Cos(fPlayerA) * fSpeed * fElapsedTime
						}
					}
				}

			}

		}*/
		fmt.Printf("fPlayerX: %d, fPlayerY: %d, fPlayerA: %f, FPS: %f", int(fPlayerX), int(fPlayerY), fPlayerA, float64(1)/fElapsedTime)
		for x := 0; x < nScreenWidth; x++ {
			fRayAngle := float64((fPlayerA - fFOV/2.0) + (float64(x)/nScreenWidth)*fFOV)
			fStepSize := float64(1.0)
			fDistanceToWall := float64(0)
			bHitWall := false
			bBoundry := false

			fEyeX := math.Sin(fRayAngle)
			fEyeY := math.Cos(fRayAngle)

			for !bHitWall && fDistanceToWall < fDepth {
				fDistanceToWall += fStepSize
				nTestX := int(fPlayerX + fEyeX*fDistanceToWall)
				nTestY := int(fPlayerY + fEyeY*fDistanceToWall)

				if nTestX < 0 || nTestX >= nMapWidth || nTestY < 0 || nTestY >= nMapHeight {
					bHitWall = true
					fDistanceToWall = fDepth
				} else {
					if ourMap[(nTestX*nMapWidth+nTestY)] == byte('#') {
						bHitWall = true
						p := []Pair{}

						for tx := 0; tx < 2; tx++ {
							for ty := 0; ty < 2; ty++ {
								vy := float64(float64(nTestY) + float64(ty) - fPlayerY)
								vx := float64(float64(nTestX) + float64(tx) - fPlayerX)
								d := math.Sqrt((vx*vx + vy*vy))
								dot := float64((fEyeX * vx / d) + (fEyeY * vy / d))
								p = append(p, Pair{d, dot})
							}
						}

						sort.Sort(ByD(p))
						fBound := float64(0.01)
						if math.Acos(p[0].Dot) < fBound {
							bBoundry = true
						}
						if math.Acos(p[1].Dot) < fBound {
							bBoundry = true
						}

						if math.Acos(p[2].Dot) < fBound {
							bBoundry = true
						}
					}
				}
			}
			nCeiling := int((nScreenHeight / 2.0) - nScreenHeight/float64(fDistanceToWall))
			nFloor := int(nScreenHeight - nCeiling)

			nShade := rune(' ')
			if fDistanceToWall <= fDepth/4.0 {
				nShade = rune(0x2588)
			} else if fDistanceToWall < fDepth/3.0 {
				nShade = rune(0x2593)
			} else if fDistanceToWall < fDepth/2.0 {
				nShade = rune(0x2592)
			} else if fDistanceToWall < fDepth {
				nShade = rune(0x2591)
			} else {
				nShade = rune(' ')
			}

			if bBoundry {
				nShade = rune(' ')
			}

			for y := 0; y < nScreenHeight; y++ {
				if y <= nCeiling {
					screen.SetContent(x, y, rune(' '), nil, styleDefault)
				} else if y > nCeiling && y <= nFloor {
					screen.SetContent(x, y, nShade, nil, styleDefault)
				} else {
					b := float64(1.0 - (float64(y-nScreenHeight/2.0))/(float64(nScreenHeight/2.0)))
					if b < 0.25 {
						nShade = rune('#')
					} else if b < 0.5 {
						nShade = rune('x')
					} else if b < 0.75 {
						nShade = rune('.')
					} else if b < 0.9 {
						nShade = rune('-')
					} else {
						nShade = rune(' ')
					}
					screen.SetContent(x, y, nShade, nil, styleDefault)
				}
			}
		}
		for nx := 0; nx < nMapWidth; nx++ {
			for ny := 0; ny < nMapHeight; ny++ {
				screen.SetContent(nx, ny, rune(ourMap[ny*nMapWidth+nx]), nil, styleDefault)
			}
		}
		screen.SetContent(int(fPlayerX)+1, int(fPlayerY), rune('P'), nil, styleDefault)
		screen.Show()
		/*key := tcell.Key(rune('W'))
		evt := tcell.NewEventKey(key, rune('W'), tcell.ModNone)
		screen.PostEvent(evt)*/

	}
	defer screen.Fini()
	//fmt.Printf("Hello golang")
	//someStr := "This is a nice string to be reversed"
	/*http.HandleFunc("/", MyHandler)
	s := &http.Server{
		Addr:           ":8000",
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())*/
	//fmt.Printf("%s", reverse(someStr))
}
