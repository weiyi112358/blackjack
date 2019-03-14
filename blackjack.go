package main

import (
	"log"
	"net"
	"net/http" 
    //"fmt"
    "./blackjack"
    "./cards"
	"github.com/zserge/webview"
)



var indexHTML = `
<!doctype html>
<html>
	<head>
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
	</head>
	<body>
            <button onclick="external.invoke('start')">Start</button>
            <button onclick="external.invoke('hit')">Hit</button>
            <button onclick="external.invoke('stand')">Stand</button>
			<div id="dealer">
				<h2> dealer</h2>
				<p>blank</p>
				
				
			</div>

			<div id = "player1">
				<h2>player1</h2>
				<p>blank</p>
				<p class = "wallet"></p>
				<p class = "bet"></p>
			</div>

			<div id = "player2">
				<h2>player2</h2>
				<p>blank</p>
				<p class = "wallet"></p>
				<p class = "bet"></p>
			</div>

			<div id = "player3">
				<h2>player3</h2>
				<p>blank</p>
				<p class = "wallet"></p>
				<p class = "bet"></p>
			</div>
	</body>
	<style> 
		#dealer
		{
			position:absolute;
			left:400px;
			top:50px
		}
		#player1
		{
			position:absolute;
			left:100px;
			top:300px
		}
		#player2
		{
			position:absolute;
			left:350px;
			top:300px
		}
		#player3
		{
			position:absolute;
			left:600px;
			top:300px
		}		 
	</style> 
	<script>
		
	</script>

</html>
`
var deck = cards.NewDeck(6)
var dealersHand = blackjack.NewHand(deck)
var playersHand = blackjack.NewHand(deck)

func startServer() string {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer ln.Close()
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(indexHTML))
		})
		log.Fatal(http.Serve(ln, nil))
	}()
	return "http://" + ln.Addr().String()
}

func handleRPC(w webview.WebView, data string) {
	switch {
	case data == "stand":
		stand(w)
	case data == "start":
		start(w)
	case data == "hit":
		hit(w)
	}
}

func main() {
    url := startServer()
    
	w := webview.New(webview.Settings{
		Width:     800,
		Height:    800,
		Title:     "blackjack demo",
		Resizable: true,
		URL:       url,
		ExternalInvokeCallback: handleRPC,
    })
    

	defer w.Exit()
	w.Run()
}


func start(w webview.WebView){

    
    playersHand.bet = playersHand.putBet(5)

    w.Dispatch(func() {
        w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML = "`+dealersHand.Cards[0].ToStr()+`";`)
		w.Eval(`document.getElementById("player1").getElementsByTagName("p")[0].innerHTML = "`+playersHand.ToStr()+`";`)
		w.Eval(`document.getElementById("player1").getElementsByClassName("wallet")[0].innerHTML = "`+ playersHand.wallet+`";`)
		w.Eval(`document.getElementById("player1").getElementsByClassName("bet")[0].innerHTML = "`+ playersHand.bet+`";`)
    })
    
    

    if playersHand.IsBlackjack() {
        w.Dialog(webview.DialogTypeAlert, 0, "Blackjack", "Blackjack!")

        if dealersHand.ScoreCard(0) >= 10 {
            w.Dispatch(func() {
                w.Eval(`document.getElementById("dealer").innerHTML = `+dealersHand.ToStr()+`;`)
            })
        }

        if dealersHand.IsBlackjack() {
            w.Dialog(webview.DialogTypeAlert, 0, "Blackjack", "Blackjack!The game is a push\n")
        } else {
            w.Dialog(webview.DialogTypeAlert, 0, "You win!", "You win!\n")
        }
    } else {
       // playHand(dealersHand, playersHand)
    }
        
}
func hit(w webview.WebView){
    card := playersHand.Hit()
    w.Dispatch(func() {
        w.Eval(`document.getElementById("player1").getElementsByTagName("p")[0].innerHTML += ",`+card.ToStr()+`";`)
    })
    if playersHand.Score() > 21 {
        w.Dialog(webview.DialogTypeAlert, 0, "You bust!", "You bust!\n")    
    }

    w.Dispatch(func() {
        w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML += ",`+dealersHand.Cards[1].ToStr()+`";`)        
    })
    w.Terminate()

    
}


func stand(w webview.WebView){
    for dealersHand.Score() < 17 {
        card := dealersHand.Hit()
        w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML += ",`+card.ToStr()+`";`)
    }

    if dealersHand.Score() > 21 {
        w.Dialog(webview.DialogTypeAlert, 0, "dealer bust!", "dealer bust!\n")
        w.Terminate()
    }

    if dealersHand.Score() > playersHand.Score() {
        w.Dialog(webview.DialogTypeAlert, 0, "dealer win!", "dealer win!\n")
	w.Terminate()
    } else if dealersHand.Score() < playersHand.Score() {
        w.Dialog(webview.DialogTypeAlert, 0, "you win!", "you win!\n")
	w.Terminate()
    } else {
        w.Dialog(webview.DialogTypeAlert, 0, "push", "push\n")
	w.Terminate()
    }

}



func getBet(wallet *float64) int {
	bet := 0
	valid := false
	for valid == false {
		valid = true
		bet = 5
		
		if float64(bet) > *wallet {
			valid = false
		}
		
	}
	*wallet = *wallet - float64(bet)
	return bet
}
