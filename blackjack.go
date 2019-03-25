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
            <button id = "start" onclick="external.invoke('start')">Start</button>
            <button id = "hit" onclick="external.invoke('hit')">Hit</button>
			<button id = "stand" onclick="external.invoke('stand')">Stand</button>
			<button id = "double" onclick="external.invoke('double')">Double</button>
			<button id = "split" onclick="external.invoke('split')">Split</button>
			<button id = "settle" onclick="external.invoke('settle')">Settle</button>
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
				<p class = "splithand"></p>
			</div>

			<div id = "player2">
				<h2>player2</h2>
				<p>blank</p>
				<p class = "wallet"></p>
				<p class = "bet"></p>
				<p class = "splithand"></p>
			</div>

			<div id = "player3">
				<h2>player3</h2>
				<p>blank</p>
				<p class = "wallet"></p>
				<p class = "bet"></p>
				<p class = "splithand"></p>
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
			left:500px;
			top:300px
		}
		#player3
		{
			position:absolute;
			left:900px;
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
		#hit
		{
			visibility:hidden
		}
		#stand
		{
			visibility:hidden
		}
		#double
		{
			visibility:hidden
		}
		#split
		{
			visibility:hidden
		}
		#settle
		{
			visibility:hidden
		}		 
	</style> 
	<script>
		
	</script>

</html>
`
var deck = cards.NewDeck(6)
var dealersHand = new(blackjack.Hand)
var player1Hand = new(blackjack.Hand)
var player2Hand = new(blackjack.Hand)
var player3Hand = new(blackjack.Hand)
var player1Split = new(blackjack.Hand)
var player2Split = new(blackjack.Hand)
var player3Split = new(blackjack.Hand)

var turn = "blank"

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
	case data == "split":
		split(w)
	}
}

func main() {
	url := startServer()

	w := webview.New(webview.Settings{
		Width:                  1200,
		Height:                 800,
		Title:                  "blackjack demo",
		Resizable:              true,
		URL:                    url,
		ExternalInvokeCallback: handleRPC,
	})

	defer w.Exit()
	w.Run()
}

func double(w webview.WebView) {
	if turn == "player1's turn" {
		if player1Hand.CanDouble() && player1Hand.CanBet() {
			player1Hand.Double()
			if player1Hand.Score() > 21 {
				w.Dialog(webview.DialogTypeAlert, 0, "Player1", "Player1 bust!\n")
			}
			if player2Hand.Bet != 0 {
				turn = "player2's turn"
			} else if player3Hand.Bet != 0 {
				turn = "player3's turn"
			} else {
				turn = "dealer's turn"
				w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML += ",` + dealersHand.Cards[1].ToStr() + `";`)
				w.Eval(`document.getElementById("settle").style.visibility = "visible";`)
				w.Eval(`document.getElementById("hit").style.visibility = "hidden";`)
				w.Eval(`document.getElementById("stand").style.visibility = "hidden";`)
				w.Eval(`document.getElementById("double").style.visibility = "hidden";`)
			}

			w.Dispatch(func() {
				w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
				w.Eval(`document.getElementById("player1").getElementsByTagName("p")[0].innerHTML = "` + player1Hand.ToStr() + `";`)
				w.Eval(`document.getElementById("player1").getElementsByClassName("wallet")[0].innerHTML = "` + strconv.Itoa(player1Hand.Wallet) + `";`)
				w.Eval(`document.getElementById("player1").getElementsByClassName("bet")[0].innerHTML = "` + strconv.Itoa(player1Hand.Bet) + `";`)
			})
		} else {
			w.Dialog(webview.DialogTypeAlert, 0, "Player1", "Player1 can not double!\n")
		}
	} else if turn == "player2's turn" {
		if player2Hand.CanDouble() && player2Hand.CanBet() {
			player2Hand.Double()
			if player2Hand.Score() > 21 {
				w.Dialog(webview.DialogTypeAlert, 0, "Player2", "Player2 bust!\n")
			}
			if player3Hand.Bet != 0 {
				turn = "player3's turn"
			} else {
				turn = "dealer's turn"
				w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML += ",` + dealersHand.Cards[1].ToStr() + `";`)
				w.Eval(`document.getElementById("settle").style.visibility = "visible";`)
				w.Eval(`document.getElementById("hit").style.visibility = "hidden";`)
				w.Eval(`document.getElementById("stand").style.visibility = "hidden";`)
				w.Eval(`document.getElementById("double").style.visibility = "hidden";`)
			}
			w.Dispatch(func() {
				w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
				w.Eval(`document.getElementById("player2").getElementsByTagName("p")[0].innerHTML = "` + player2Hand.ToStr() + `";`)
				w.Eval(`document.getElementById("player2").getElementsByClassName("wallet")[0].innerHTML = "` + strconv.Itoa(player2Hand.Wallet) + `";`)
				w.Eval(`document.getElementById("player2").getElementsByClassName("bet")[0].innerHTML = "` + strconv.Itoa(player2Hand.Bet) + `";`)
			})
		} else {
			w.Dialog(webview.DialogTypeAlert, 0, "Player2", "Player2 can not double!\n")
		}
	} else if turn == "player3's turn" {
		if player3Hand.CanDouble() && player3Hand.CanBet() {
			player3Hand.Double()
			if player3Hand.Score() > 21 {
				w.Dialog(webview.DialogTypeAlert, 0, "Player3", "Player3 bust!\n")
			}
			turn = "dealer's turn"
			w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML += ",` + dealersHand.Cards[1].ToStr() + `";`)
			w.Eval(`document.getElementById("settle").style.visibility = "visible";`)
			w.Eval(`document.getElementById("hit").style.visibility = "hidden";`)
			w.Eval(`document.getElementById("stand").style.visibility = "hidden";`)
			w.Eval(`document.getElementById("double").style.visibility = "hidden";`)
			w.Dispatch(func() {
				w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
				w.Eval(`document.getElementById("player3").getElementsByTagName("p")[0].innerHTML = "` + player3Hand.ToStr() + `";`)
				w.Eval(`document.getElementById("player3").getElementsByClassName("wallet")[0].innerHTML = "` + strconv.Itoa(player3Hand.Wallet) + `";`)
				w.Eval(`document.getElementById("player3").getElementsByClassName("bet")[0].innerHTML = "` + strconv.Itoa(player3Hand.Bet) + `";`)
			})
		} else {
			w.Dialog(webview.DialogTypeAlert, 0, "Player3", "Player3 can not double!\n")
		}
	} else if turn == "dealer's turn" {
		w.Dialog(webview.DialogTypeAlert, 0, "Dealer", "dealer can not double!\n")
	}

}

func settle(w webview.WebView) {
	if turn == "dealer's turn" {
		for dealersHand.Score() < 17 {
			card := dealersHand.Hit()
			w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML += ",` + card.ToStr() + `";`)
		}
		w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML = "` + dealersHand.ToStr() + `";`)
		if player1Hand.HasSplit {
			if player1Hand.Score() == 21 || player1Split.Score() == 21 {
				if dealersHand.Score() != 21 {
					player1Hand.Wallet += player1Hand.Bet * 2
				}
			} else if player1Hand.IsBust == true && player1Split.IsBust == true {
				if dealersHand.Score() > 21 {
					player1Hand.Wallet += player1Hand.Bet
				}
			} else {
				if dealersHand.Score() > 21 {
					player1Hand.Wallet += player1Hand.Bet * 2
				} else if dealersHand.Score() > player1Hand.Score() && dealersHand.Score() > player1Split.Score() {

				} else if dealersHand.Score() < player1Hand.Score() || dealersHand.Score() < player1Split.Score() {
					player1Hand.Wallet += player1Hand.Bet * 2
				} else {
					player1Hand.Wallet += player1Hand.Bet
				}
			}
		} else {
			if player1Hand.Score() == 21 {
				if dealersHand.Score() != 21 {
					player1Hand.Wallet += player1Hand.Bet * 2
				}
			} else if player1Hand.IsBust == true {
				if dealersHand.Score() > 21 {
					player1Hand.Wallet += player1Hand.Bet
				}
			} else {
				if dealersHand.Score() > 21 {
					player1Hand.Wallet += player1Hand.Bet * 2
				} else if dealersHand.Score() > player1Hand.Score() {

				} else if dealersHand.Score() < player1Hand.Score() {
					player1Hand.Wallet += player1Hand.Bet * 2
				} else {
					player1Hand.Wallet += player1Hand.Bet
				}
			}
		}
		if player2Hand.HasSplit {
			if player2Hand.Score() == 21 || player2Split.Score() == 21 {
				if dealersHand.Score() != 21 {
					player2Hand.Wallet += player2Hand.Bet * 2
				}
			} else if player2Hand.IsBust == true && player2Split.IsBust == true {
				if dealersHand.Score() > 21 {
					player2Hand.Wallet += player2Hand.Bet
				}
			} else {
				if dealersHand.Score() > 21 {
					player2Hand.Wallet += player2Hand.Bet * 2
				} else if dealersHand.Score() > player2Hand.Score() && dealersHand.Score() > player2Split.Score() {

				} else if dealersHand.Score() < player2Hand.Score() || dealersHand.Score() < player2Split.Score() {
					player2Hand.Wallet += player2Hand.Bet * 2
				} else {
					player2Hand.Wallet += player2Hand.Bet
				}
			}
		} else {
			if player2Hand.Score() == 21 {
				if dealersHand.Score() != 21 {
					player2Hand.Wallet += player2Hand.Bet * 2
				}
			} else if player2Hand.IsBust == true {
				if dealersHand.Score() > 21 {
					player2Hand.Wallet += player2Hand.Bet
				}
			} else {
				if dealersHand.Score() > 21 {
					player2Hand.Wallet += player2Hand.Bet * 2
				} else if dealersHand.Score() > player2Hand.Score() {

				} else if dealersHand.Score() < player2Hand.Score() {
					player2Hand.Wallet += player2Hand.Bet * 2
				} else {
					player2Hand.Wallet += player2Hand.Bet
				}
			}
		}
		if player3Hand.HasSplit {
			if player3Hand.Score() == 21 || player3Split.Score() == 21 {
				if dealersHand.Score() != 21 {
					player3Hand.Wallet += player3Hand.Bet * 2
				}
			} else if player3Hand.IsBust == true && player3Split.IsBust == true {
				if dealersHand.Score() > 21 {
					player3Hand.Wallet += player3Hand.Bet
				}
			} else {
				if dealersHand.Score() > 21 {
					player3Hand.Wallet += player3Hand.Bet * 2
				} else if dealersHand.Score() > player3Hand.Score() && dealersHand.Score() > player3Split.Score() {

				} else if dealersHand.Score() < player3Hand.Score() || dealersHand.Score() < player3Split.Score() {
					player3Hand.Wallet += player3Hand.Bet * 2
				} else {
					player3Hand.Wallet += player3Hand.Bet
				}
			}
		} else {
			if player3Hand.Score() == 21 {
				if dealersHand.Score() != 21 {
					player3Hand.Wallet += player3Hand.Bet * 2
				}
			} else if player3Hand.IsBust == true {
				if dealersHand.Score() > 21 {
					player3Hand.Wallet += player3Hand.Bet
				}
			} else {
				if dealersHand.Score() > 21 {
					player3Hand.Wallet += player3Hand.Bet * 2
				} else if dealersHand.Score() > player3Hand.Score() {

				} else if dealersHand.Score() < player3Hand.Score() {
					player3Hand.Wallet += player3Hand.Bet * 2
				} else {
					player3Hand.Wallet += player3Hand.Bet
				}
			}
		}

		turn = "isSettled"

		w.Dispatch(func() {
			w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
			w.Eval(`document.getElementById("player1").getElementsByClassName("wallet")[0].innerHTML = "` + strconv.Itoa(player1Hand.Wallet) + `";`)
			w.Eval(`document.getElementById("player2").getElementsByClassName("wallet")[0].innerHTML = "` + strconv.Itoa(player2Hand.Wallet) + `";`)
			w.Eval(`document.getElementById("player3").getElementsByClassName("wallet")[0].innerHTML = "` + strconv.Itoa(player3Hand.Wallet) + `";`)
		})
		w.Eval(`document.getElementById("settle").style.visibility = "hidden";`)
		w.Eval(`document.getElementById("start").style.visibility = "visible";`)

		if player1Hand.HasSplit {
			player1Split.Flush()
			w.Eval(`document.getElementById("player1").getElementsByClassName("splithand")[0].style.visibility = "hidden";`)
		}
		if player2Hand.HasSplit {
			player2Split.Flush()
			w.Eval(`document.getElementById("player2").getElementsByClassName("splithand")[0].style.visibility = "hidden";`)
		}
		if player3Hand.HasSplit {
			player3Split.Flush()
			w.Eval(`document.getElementById("player3").getElementsByClassName("splithand")[0].style.visibility = "hidden";`)
		}

		player1Hand.Flush()
		player2Hand.Flush()
		player3Hand.Flush()

	}

}

func start(w webview.WebView) {

	if turn == "blank" {
		dealersHand = blackjack.NewHand(deck)
		player1Hand = blackjack.NewHand(deck)
		player2Hand = blackjack.NewHand(deck)
		player3Hand = blackjack.NewHand(deck)
		player1Hand.PutBet()
		player2Hand.PutBet()
		player3Hand.PutBet()

		turn = "player1's turn"

		w.Dispatch(func() {
			w.Eval(`document.getElementById("hit").style.visibility = "visible";`)
			w.Eval(`document.getElementById("stand").style.visibility = "visible";`)
			w.Eval(`document.getElementById("double").style.visibility = "visible";`)
			w.Eval(`document.getElementById("split").style.visibility = "visible";`)

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

		w.Dispatch(func() {
			w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
			w.Eval(`document.getElementById("start").style.visibility = "hidden";`)
			w.Eval(`document.getElementById("settle").style.visibility = "hidden";`)
		})

	} else if turn == "isSettled" {
		dealersHand = blackjack.NewHand(deck)
		if player1Hand.CanBet() {
			player1Hand.Refresh()
			player1Hand.PutBet()
			turn = "player1's turn"
			if player2Hand.CanBet() {
				player2Hand.Refresh()
				player2Hand.PutBet()
			}
			if player3Hand.CanBet() {
				player3Hand.Refresh()
				player3Hand.PutBet()
			}

		} else {
			if player2Hand.CanBet() {
				player2Hand.Refresh()
				player2Hand.PutBet()
				turn = "player2's turn"
				if player3Hand.CanBet() {
					player3Hand.Refresh()
					player3Hand.PutBet()
				}
			} else {
				if player3Hand.CanBet() {
					player3Hand.Refresh()
					player3Hand.PutBet()
					turn = "player3's turn"
				} else {
					w.Dialog(webview.DialogTypeAlert, 0, "Players", "Your wallet is empty! Game Over!\n")
					w.Terminate()
				}
			}

		}

		w.Dispatch(func() {
			w.Eval(`document.getElementById("hit").style.visibility = "visible";`)
			w.Eval(`document.getElementById("stand").style.visibility = "visible";`)
			w.Eval(`document.getElementById("double").style.visibility = "visible";`)
			w.Eval(`document.getElementById("split").style.visibility = "visible";`)

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

		w.Dispatch(func() {
			w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
			w.Eval(`document.getElementById("start").style.visibility = "hidden";`)
			w.Eval(`document.getElementById("settle").style.visibility = "hidden";`)
		})

	} else {
		w.Dialog(webview.DialogTypeAlert, 0, "Game is not finished", "Game should be settled\n")
	}

}
func hit(w webview.WebView) {
	if turn == "player1's turn" {
		card := player1Hand.Hit()
		w.Dispatch(func() {
			w.Eval(`document.getElementById("player1").getElementsByTagName("p")[0].innerHTML += ",` + card.ToStr() + `";`)
		})
		if player1Hand.Score() > 21 {
			w.Dialog(webview.DialogTypeAlert, 0, "Player1", "Player1 bust!\n")
			player1Hand.IsBust = true

			if player1Hand.HasSplit {
				turn = "player1's Split turn"
			} else if player2Hand.Bet != 0 {
				turn = "player2's turn"
			} else if player3Hand.Bet != 0 {
				turn = "player3's turn"
			} else {
				turn = "dealer's turn"
				w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML += ",` + dealersHand.Cards[1].ToStr() + `";`)
				w.Eval(`document.getElementById("settle").style.visibility = "visible";`)
				w.Eval(`document.getElementById("hit").style.visibility = "hidden";`)
				w.Eval(`document.getElementById("stand").style.visibility = "hidden";`)
				w.Eval(`document.getElementById("double").style.visibility = "hidden";`)
				w.Eval(`document.getElementById("split").style.visibility = "hidden";`)
			}
		}
		w.Dispatch(func() {
			w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
		})
	} else if turn == "player1's Split turn" {
		player1Split.Hit()

		w.Eval(`document.getElementById("player1").getElementsByClassName("splithand")[0].innerHTML = "` + player1Split.ToStr() + `";`)
		if player1Split.Score() > 21 {
			w.Dialog(webview.DialogTypeAlert, 0, "Player1", "Player1 Split bust!\n")
			player1Split.IsBust = true

			if player2Hand.Bet != 0 {
				turn = "player2's turn"
			} else if player3Hand.Bet != 0 {
				turn = "player3's turn"
			} else {
				turn = "dealer's turn"
				w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML += ",` + dealersHand.Cards[1].ToStr() + `";`)
				w.Eval(`document.getElementById("settle").style.visibility = "visible";`)
				w.Eval(`document.getElementById("hit").style.visibility = "hidden";`)
				w.Eval(`document.getElementById("stand").style.visibility = "hidden";`)
				w.Eval(`document.getElementById("double").style.visibility = "hidden";`)
				w.Eval(`document.getElementById("split").style.visibility = "hidden";`)
			}
		}
		w.Dispatch(func() {
			w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
		})
	} else if turn == "player2's turn" {
		card := player2Hand.Hit()
		w.Dispatch(func() {
			w.Eval(`document.getElementById("player2").getElementsByTagName("p")[0].innerHTML += ",` + card.ToStr() + `";`)
		})
		if player2Hand.Score() > 21 {
			w.Dialog(webview.DialogTypeAlert, 0, "Player2", "Player2 bust!\n")
			player2Hand.IsBust = true
			if player1Hand.HasSplit {
				turn = "player2's Split turn"
			} else if player3Hand.Bet != 0 {
				turn = "player3's turn"
			} else {
				turn = "dealer's turn"
				w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML += ",` + dealersHand.Cards[1].ToStr() + `";`)
				w.Eval(`document.getElementById("settle").style.visibility = "visible";`)
				w.Eval(`document.getElementById("hit").style.visibility = "hidden";`)
				w.Eval(`document.getElementById("stand").style.visibility = "hidden";`)
				w.Eval(`document.getElementById("double").style.visibility = "hidden";`)
				w.Eval(`document.getElementById("split").style.visibility = "hidden";`)
			}
		}
		w.Dispatch(func() {
			w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
		})

	} else if turn == "player2's Split turn" {
		player2Split.Hit()

		w.Eval(`document.getElementById("player2").getElementsByClassName("splithand")[0].innerHTML = "` + player2Split.ToStr() + `";`)
		if player2Split.Score() > 21 {
			w.Dialog(webview.DialogTypeAlert, 0, "Player2", "Player2 Split bust!\n")
			player2Split.IsBust = true

			if player3Hand.Bet != 0 {
				turn = "player3's turn"
			} else {
				turn = "dealer's turn"
				w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML += ",` + dealersHand.Cards[1].ToStr() + `";`)
				w.Eval(`document.getElementById("settle").style.visibility = "visible";`)
				w.Eval(`document.getElementById("hit").style.visibility = "hidden";`)
				w.Eval(`document.getElementById("stand").style.visibility = "hidden";`)
				w.Eval(`document.getElementById("double").style.visibility = "hidden";`)
				w.Eval(`document.getElementById("split").style.visibility = "hidden";`)
			}
		}
		w.Dispatch(func() {
			w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
		})
	} else if turn == "player3's turn" {
		card := player3Hand.Hit()
		w.Dispatch(func() {
			w.Eval(`document.getElementById("player3").getElementsByTagName("p")[0].innerHTML += ",` + card.ToStr() + `";`)
		})
		if player3Hand.Score() > 21 {
			w.Dialog(webview.DialogTypeAlert, 0, "Player3", "Player3 bust!\n")
			w.Dispatch(func() {
				w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML += ",` + dealersHand.Cards[1].ToStr() + `";`)
			})
			player3Hand.IsBust = true
			if player3Hand.HasSplit {
				turn = "player3's Split turn"
			} else {
				turn = "dealer's turn"

				w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML += ",` + dealersHand.Cards[1].ToStr() + `";`)
				w.Eval(`document.getElementById("settle").style.visibility = "visible";`)
				w.Eval(`document.getElementById("hit").style.visibility = "hidden";`)
				w.Eval(`document.getElementById("stand").style.visibility = "hidden";`)
				w.Eval(`document.getElementById("double").style.visibility = "hidden";`)
				w.Eval(`document.getElementById("split").style.visibility = "hidden";`)

			}
		}
		w.Dispatch(func() {
			w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
		})
	} else if turn == "player3's Split turn" {
		player3Split.Hit()

		w.Eval(`document.getElementById("player3").getElementsByClassName("splithand")[0].innerHTML = "` + player3Split.ToStr() + `";`)
		if player3Split.Score() > 21 {
			w.Dialog(webview.DialogTypeAlert, 0, "Player3", "Player3 Split bust!\n")
			player3Split.IsBust = true

			turn = "dealer's turn"
			w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML += ",` + dealersHand.Cards[1].ToStr() + `";`)
			w.Eval(`document.getElementById("settle").style.visibility = "visible";`)
			w.Eval(`document.getElementById("hit").style.visibility = "hidden";`)
			w.Eval(`document.getElementById("stand").style.visibility = "hidden";`)
			w.Eval(`document.getElementById("double").style.visibility = "hidden";`)
			w.Eval(`document.getElementById("split").style.visibility = "hidden";`)

		}
		w.Dispatch(func() {
			w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
		})
	}
}

func stand(w webview.WebView) {
	if turn == "player1's turn" {
		if player1Hand.HasSplit {
			turn = "player1's Split turn"
		} else if player2Hand.Bet != 0 {
			turn = "player2's turn"
		} else if player3Hand.Bet != 0 {
			turn = "player3's turn"
		} else {
			turn = "dealer's turn"
			w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML += ",` + dealersHand.Cards[1].ToStr() + `";`)
			w.Eval(`document.getElementById("settle").style.visibility = "visible";`)
			w.Eval(`document.getElementById("hit").style.visibility = "hidden";`)
			w.Eval(`document.getElementById("stand").style.visibility = "hidden";`)
			w.Eval(`document.getElementById("double").style.visibility = "hidden";`)
			w.Eval(`document.getElementById("split").style.visibility = "hidden";`)
		}
		w.Dispatch(func() {
			w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
		})
	} else if turn == "player1's Split turn" {
		if player2Hand.Bet != 0 {
			turn = "player2's turn"
		} else if player3Hand.Bet != 0 {
			turn = "player3's turn"
		} else {
			turn = "dealer's turn"
			w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML += ",` + dealersHand.Cards[1].ToStr() + `";`)
			w.Eval(`document.getElementById("settle").style.visibility = "visible";`)
			w.Eval(`document.getElementById("hit").style.visibility = "hidden";`)
			w.Eval(`document.getElementById("stand").style.visibility = "hidden";`)
			w.Eval(`document.getElementById("double").style.visibility = "hidden";`)
			w.Eval(`document.getElementById("split").style.visibility = "hidden";`)
		}
		w.Dispatch(func() {
			w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
		})
	} else if turn == "player2's turn" {
		if player2Hand.HasSplit {
			turn = "player2's Split turn"
		} else if player3Hand.Bet != 0 {
			turn = "player3's turn"
		} else {
			turn = "dealer's turn"
			w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML += ",` + dealersHand.Cards[1].ToStr() + `";`)
			w.Eval(`document.getElementById("settle").style.visibility = "visible";`)
			w.Eval(`document.getElementById("hit").style.visibility = "hidden";`)
			w.Eval(`document.getElementById("stand").style.visibility = "hidden";`)
			w.Eval(`document.getElementById("double").style.visibility = "hidden";`)

		}
		w.Dispatch(func() {
			w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
		})
	} else if turn == "player2's Split turn" {
		if player3Hand.Bet != 0 {
			turn = "player3's turn"
		} else {
			turn = "dealer's turn"
			w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML += ",` + dealersHand.Cards[1].ToStr() + `";`)
			w.Eval(`document.getElementById("settle").style.visibility = "visible";`)
			w.Eval(`document.getElementById("hit").style.visibility = "hidden";`)
			w.Eval(`document.getElementById("stand").style.visibility = "hidden";`)
			w.Eval(`document.getElementById("double").style.visibility = "hidden";`)

		}
		w.Dispatch(func() {
			w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
		})
	} else if turn == "player3's turn" {
		if player1Hand.HasSplit {
			turn = "player3's Split turn"
		} else {
			turn = "dealer's turn"
			w.Eval(`document.getElementById("settle").style.visibility = "visible";`)
			w.Eval(`document.getElementById("hit").style.visibility = "hidden";`)
			w.Eval(`document.getElementById("stand").style.visibility = "hidden";`)
			w.Eval(`document.getElementById("double").style.visibility = "hidden";`)
			w.Eval(`document.getElementById("split").style.visibility = "hidden";`)
			w.Dispatch(func() {
				w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML += ",` + dealersHand.Cards[1].ToStr() + `";`)
				w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
			})
		}
	} else if turn == "player3's Split turn" {
		turn = "dealer's turn"
		w.Eval(`document.getElementById("settle").style.visibility = "visible";`)
		w.Eval(`document.getElementById("hit").style.visibility = "hidden";`)
		w.Eval(`document.getElementById("stand").style.visibility = "hidden";`)
		w.Eval(`document.getElementById("double").style.visibility = "hidden";`)
		w.Eval(`document.getElementById("split").style.visibility = "hidden";`)
		w.Dispatch(func() {
			w.Eval(`document.getElementById("dealer").getElementsByTagName("p")[0].innerHTML += ",` + dealersHand.Cards[1].ToStr() + `";`)
			w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
		})
	}
}

func split(w webview.WebView) {
	if turn == "player1's turn" {
		if player1Hand.CanSplit() {
			player1Split = player1Hand.Split()
			player1Hand.HasSplit = true
			w.Eval(`document.getElementById("player1").getElementsByClassName("splithand")[0].style.visibility = "visible";`)

			w.Dispatch(func() {
				w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
				w.Eval(`document.getElementById("player1").getElementsByTagName("p")[0].innerHTML = "` + player1Hand.ToStr() + `";`)
				w.Eval(`document.getElementById("player1").getElementsByClassName("wallet")[0].innerHTML = "` + strconv.Itoa(player1Hand.Wallet) + `";`)
				w.Eval(`document.getElementById("player1").getElementsByClassName("bet")[0].innerHTML = "` + strconv.Itoa(player1Hand.Bet) + `";`)
				w.Eval(`document.getElementById("player1").getElementsByClassName("splithand")[0].innerHTML = "` + player1Split.ToStr() + `";`)
			})
		} else {
			w.Dialog(webview.DialogTypeAlert, 0, "Player1", "Player1 can not split!\n")
		}
	} else if turn == "player2's turn" {
		if player2Hand.CanSplit() {
			player2Split = player2Hand.Split()
			player2Hand.HasSplit = true
			w.Eval(`document.getElementById("player2").getElementsByClassName("splithand")[0].style.visibility = "visible";`)

			w.Dispatch(func() {
				w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
				w.Eval(`document.getElementById("player2").getElementsByTagName("p")[0].innerHTML = "` + player2Hand.ToStr() + `";`)
				w.Eval(`document.getElementById("player2").getElementsByClassName("wallet")[0].innerHTML = "` + strconv.Itoa(player2Hand.Wallet) + `";`)
				w.Eval(`document.getElementById("player2").getElementsByClassName("bet")[0].innerHTML = "` + strconv.Itoa(player2Hand.Bet) + `";`)
				w.Eval(`document.getElementById("player2").getElementsByClassName("splithand")[0].innerHTML = "` + player2Split.ToStr() + `";`)
			})
		} else {
			w.Dialog(webview.DialogTypeAlert, 0, "Player2", "Player2 can not split!\n")
		}
	} else if turn == "player3's turn" {
		if player3Hand.CanSplit() {
			player3Split = player3Hand.Split()
			player3Hand.HasSplit = true
			w.Eval(`document.getElementById("player3").getElementsByClassName("splithand")[0].style.visibility = "visible";`)

			w.Dispatch(func() {
				w.Eval(`document.getElementById("turn").innerHTML = "` + turn + `";`)
				w.Eval(`document.getElementById("player3").getElementsByTagName("p")[0].innerHTML = "` + player3Hand.ToStr() + `";`)
				w.Eval(`document.getElementById("player3").getElementsByClassName("wallet")[0].innerHTML = "` + strconv.Itoa(player3Hand.Wallet) + `";`)
				w.Eval(`document.getElementById("player3").getElementsByClassName("bet")[0].innerHTML = "` + strconv.Itoa(player3Hand.Bet) + `";`)
				w.Eval(`document.getElementById("player3").getElementsByClassName("splithand")[0].innerHTML = "` + player3Split.ToStr() + `";`)
			})
		} else {
			w.Dialog(webview.DialogTypeAlert, 0, "Player3", "Player3 can not split!\n")
		}
	} else {
		w.Dialog(webview.DialogTypeAlert, 0, "Players", "Can not split twice!\n")
	}
}
