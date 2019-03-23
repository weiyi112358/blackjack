package main

import (
	"log"
	"net"
	"net/http"
	"strconv"

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
			<button onclick="external.invoke('double')">Double</button>
			<button onclick="external.invoke('settle')">Settle</button>
			<label id = "turn" >Blank</label>
			<label id = "walletTitle"> Wallet:</label>
			<label id = "betTitle"> Bet:</label>
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
		#walletTitle
		{
			position:absolute;
			left:10px;
			top:400px
		}
		#betTitle
		{
			position:absolute;
			left:10px;
			top:440px
		}		 
	</style> 
	<script>
		
	</script>

</html>
`
var deck = cards.NewDeck(6)
var dealersHand = blackjack.NewHand(deck)
var player1Hand = blackjack.NewHand(deck)
var player2Hand = blackjack.NewHand(deck)
var player3Hand = blackjack.NewHand(deck)

var turn = "blank";

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
	case data == "settle":
		settle(w)
	
	case data == "double":
		double(w)
	}
}

func main() {
	url := startServer()

	w := webview.New(webview.Settings{
		Width:                  800,
		Height:                 800,
		Title:                  "blackjack demo",
		Resizable:              true,
		URL:                    url,
		ExternalInvokeCallback: handleRPC,
	})

	defer w.Exit()
	w.Run()
}

func double(w webview.WebView){
	if turn == "player1's turn"{
		player1Hand.Wallet = player1Hand.Wallet - player1Hand.Bet
		player1Hand.Bet = player1Hand.Bet*2		
		w.Dispatch(func() {
			w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
			w.Eval(`document.getElementById("player1").getElementsByClassName("wallet")[0].innerHTML = "` + strconv.Itoa(player1Hand.Wallet) + `";`)
			w.Eval(`document.getElementById("player1").getElementsByClassName("bet")[0].innerHTML = "` + strconv.Itoa(player1Hand.Bet) + `";`)
		})
	} else if turn == "player2's turn" {
		player2Hand.Wallet = player2Hand.Wallet - player2Hand.Bet
		player2Hand.Bet = player2Hand.Bet*2	
		w.Dispatch(func() {
			w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
			w.Eval(`document.getElementById("player2").getElementsByClassName("wallet")[0].innerHTML = "` + strconv.Itoa(player2Hand.Wallet) + `";`)
			w.Eval(`document.getElementById("player2").getElementsByClassName("bet")[0].innerHTML = "` + strconv.Itoa(player2Hand.Bet) + `";`)
		})
	} else if turn == "player3's turn" {
		player3Hand.Wallet = player3Hand.Wallet - player3Hand.Bet
		player3Hand.Bet = player3Hand.Bet*2	
		w.Dispatch(func() {
			w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
			w.Eval(`document.getElementById("player3").getElementsByClassName("wallet")[0].innerHTML = "` + strconv.Itoa(player3Hand.Wallet) + `";`)
			w.Eval(`document.getElementById("player3").getElementsByClassName("bet")[0].innerHTML = "` + strconv.Itoa(player3Hand.Bet) + `";`)
		})
	} 

}

func settle(w webview.WebView){
	if turn == "end"{
		for dealersHand.Score() < 17 {
			card := dealersHand.Hit()
			w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML += ",` + card.ToStr() + `";`)
		}
	
		if player1Hand.Score() == 21 {
			if dealersHand.Score() != 21 {
				player1Hand.Wallet += player1Hand.Bet*2
			}
		} else if player1Hand.IsBust ==true{
			if dealersHand.Score() > 21{
				player1Hand.Wallet += player1Hand.Bet
			}
		} else{
			if dealersHand.Score()>21{
				player1Hand.Wallet += player1Hand.Bet*2
			}else if dealersHand.Score() > player1Hand.Score() {
				
			} else if dealersHand.Score() < player1Hand.Score() {
				player1Hand.Wallet += player1Hand.Bet*2
			} else {
				player1Hand.Wallet += player1Hand.Bet
			}
		}
	
		if player2Hand.Score() == 21 {
			if dealersHand.Score() != 21 {
				player2Hand.Wallet += player2Hand.Bet*2
			}
		} else if player2Hand.IsBust ==true{
			if dealersHand.Score() > 21{
				player2Hand.Wallet += player2Hand.Bet
			}
		} else{
			if dealersHand.Score()>21{
				player2Hand.Wallet += player2Hand.Bet*2
			}else if dealersHand.Score() > player2Hand.Score() {
				
			} else if dealersHand.Score() < player2Hand.Score() {
				player2Hand.Wallet += player2Hand.Bet*2
			} else {
				player2Hand.Wallet += player2Hand.Bet
			}
		}
	
		if player3Hand.Score() == 21 {
			if dealersHand.Score() != 21 {
				player3Hand.Wallet += player3Hand.Bet*2
			}
		} else if player3Hand.IsBust ==true{
			if dealersHand.Score() > 21{
				player3Hand.Wallet += player3Hand.Bet
			}
		} else{
			if dealersHand.Score()>21{
				player3Hand.Wallet += player3Hand.Bet*2
			}else if dealersHand.Score() > player3Hand.Score() {
				
			} else if dealersHand.Score() < player3Hand.Score() {
				player3Hand.Wallet += player3Hand.Bet*2
			} else {
				player3Hand.Wallet += player3Hand.Bet
			}
		}

		w.Dispatch(func() {
			w.Eval(`document.getElementById("player1").getElementsByClassName("wallet")[0].innerHTML = "` + strconv.Itoa(player1Hand.Wallet) + `";`)
			w.Eval(`document.getElementById("player2").getElementsByClassName("wallet")[0].innerHTML = "` + strconv.Itoa(player2Hand.Wallet) + `";`)
			w.Eval(`document.getElementById("player3").getElementsByClassName("wallet")[0].innerHTML = "` + strconv.Itoa(player3Hand.Wallet) + `";`)
		})

		player1Hand.Refresh()
		player2Hand.Refresh()
		player3Hand.Refresh()
		dealersHand = blackjack.NewHand(deck)

		turn = "isSettled"

	}

	


}

func start(w webview.WebView) {

	if turn == "blank" || turn == "isSettled" {
	player1Hand.Bet = player1Hand.PutBet()
	player2Hand.Bet = player2Hand.PutBet()
	player3Hand.Bet = player3Hand.PutBet()

	w.Dispatch(func() {
		w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML = "` + dealersHand.Cards[0].ToStr() + `";`)
		w.Eval(`document.getElementById("player1").getElementsByTagName("p")[0].innerHTML = "` + player1Hand.ToStr() + `";`)
		w.Eval(`document.getElementById("player1").getElementsByClassName("wallet")[0].innerHTML = "` + strconv.Itoa(player1Hand.Wallet) + `";`)
		w.Eval(`document.getElementById("player1").getElementsByClassName("bet")[0].innerHTML = "` + strconv.Itoa(player1Hand.Bet) + `";`)

		
		w.Eval(`document.getElementById("player2").getElementsByTagName("p")[0].innerHTML = "` + player2Hand.ToStr() + `";`)
		w.Eval(`document.getElementById("player2").getElementsByClassName("wallet")[0].innerHTML = "` + strconv.Itoa(player2Hand.Wallet) + `";`)
		w.Eval(`document.getElementById("player2").getElementsByClassName("bet")[0].innerHTML = "` + strconv.Itoa(player2Hand.Bet) + `";`)

		
		w.Eval(`document.getElementById("player3").getElementsByTagName("p")[0].innerHTML = "` + player3Hand.ToStr() + `";`)
		w.Eval(`document.getElementById("player3").getElementsByClassName("wallet")[0].innerHTML = "` + strconv.Itoa(player3Hand.Wallet) + `";`)
		w.Eval(`document.getElementById("player3").getElementsByClassName("bet")[0].innerHTML = "` + strconv.Itoa(player3Hand.Bet) + `";`)
	})

	

	if player1Hand.IsBlackjack() {
		//w.Dialog(webview.DialogTypeAlert, 0, "Blackjack", "Blackjack!")

		if dealersHand.ScoreCard(0) >= 10 {
			w.Dispatch(func() {
				w.Eval(`document.getElementById("dealer").innerHTML = ` + dealersHand.ToStr() + `;`)
			})
		}

		if dealersHand.IsBlackjack() {
			//w.Dialog(webview.DialogTypeAlert, 0, "Blackjack", "Blackjack!The game is a push\n")
		} else {
			//player1Hand.Wallet+=player1Hand.Bet
			w.Eval(`document.getElementById("player1").getElementsByClassName("wallet")[0].innerHTML = "` + strconv.Itoa(player1Hand.Wallet) + `";`)
		}
	}
	if player2Hand.IsBlackjack() {
		//w.Dialog(webview.DialogTypeAlert, 0, "Blackjack", "Blackjack!")

		if dealersHand.ScoreCard(0) >= 10 {
			w.Dispatch(func() {
				w.Eval(`document.getElementById("dealer").innerHTML = ` + dealersHand.ToStr() + `;`)
			})
		}

		if dealersHand.IsBlackjack() {
			//w.Dialog(webview.DialogTypeAlert, 0, "Blackjack", "Blackjack!The game is a push\n")
		} else {
			//player2Hand.Wallet+=player2Hand.Bet
			w.Eval(`document.getElementById("player2").getElementsByClassName("wallet")[0].innerHTML = "` + strconv.Itoa(player1Hand.Wallet) + `";`)
		}
	} 
	if player3Hand.IsBlackjack() {
		//w.Dialog(webview.DialogTypeAlert, 0, "Blackjack", "Blackjack!")

		if dealersHand.ScoreCard(0) >= 10 {
			w.Dispatch(func() {
				w.Eval(`document.getElementById("dealer").innerHTML = ` + dealersHand.ToStr() + `;`)
			})
		}

		if dealersHand.IsBlackjack() {
			//w.Dialog(webview.DialogTypeAlert, 0, "Blackjack", "Blackjack!The game is a push\n")
		} else {
			//player3Hand.Wallet+=player3Hand.Bet
			w.Eval(`document.getElementById("player3").getElementsByClassName("wallet")[0].innerHTML = "` + strconv.Itoa(player1Hand.Wallet) + `";`)
		}
	}

	turn = "player1's turn"
	w.Dispatch(func() {
		w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
	})

	} else
	{
		w.Dialog(webview.DialogTypeAlert, 0, "Game is not finished,it should be settled", "Game is not finished,it should be settled\n")
	}

	
	
	

}
func hit(w webview.WebView) {
	if turn == "player1's turn" {
		card := player1Hand.Hit()
		w.Dispatch(func() {
			w.Eval(`document.getElementById("player1").getElementsByTagName("p")[0].innerHTML += ",` + card.ToStr() + `";`)
		})
		if player1Hand.Score() > 21 {
			w.Dialog(webview.DialogTypeAlert, 0, "Player1 bust!", "Player1 bust!\n")
			player1Hand.IsBust = true
			turn  = "player2's turn"

			w.Dispatch(func() {
				w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
			})			
		}
			
	} else if turn == "player2's turn" {
		card := player2Hand.Hit()
		w.Dispatch(func() {
			w.Eval(`document.getElementById("player2").getElementsByTagName("p")[0].innerHTML += ",` + card.ToStr() + `";`)
		})
		if player2Hand.Score() > 21 {
			w.Dialog(webview.DialogTypeAlert, 0, "Player2 bust!", "Player2 bust!\n")
			player2Hand.IsBust = true
			turn = "player3's turn"
			w.Dispatch(func() {
				w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
			})
		}
		
		
	} else if turn == "player3's turn" {
		card := player3Hand.Hit()
		w.Dispatch(func() {
			w.Eval(`document.getElementById("player3").getElementsByTagName("p")[0].innerHTML += ",` + card.ToStr() + `";`)
		})
		if player3Hand.Score() > 21 {
			w.Dialog(webview.DialogTypeAlert, 0, "Player3 bust!", "Player3 bust!\n")
			w.Dispatch(func() {
				w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML += ",` + dealersHand.Cards[1].ToStr() + `";`)
			})
			player3Hand.IsBust = true
			turn = "end"
			w.Dispatch(func() {
				w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
			})
		}				
	}
	

}

func stand(w webview.WebView) {
	if turn == "player1's turn"{
		turn = "player2's turn"
		w.Dispatch(func() {
			w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
		})
	} else if turn == "player2's turn" {
		turn = "player3's turn"
		w.Dispatch(func() {
			w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
		})
	} else if turn == "player3's turn" {
		turn = "end"
		w.Dispatch(func() {
			w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML += ",` + dealersHand.Cards[1].ToStr() + `";`)
			w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
		})
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
