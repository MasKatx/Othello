// mainパッケージ定義
package main

// パッケージ fmt, html/template, net/http をインポート
import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

/* テンプレート表示用データ構造体 Bord 定義
	フィールド
		Color（盤面表示）[][]string型
		Order（現在の順番）string型
		Index（配置リクエスト場所）string型
		WinFlg（勝利フラグ: n, b, w）string型
*/
type Bord struct {
	Color	[][]string
	Order	string
	Index	string
	WinFlg	string
}

/* テンプレートの設定
	ParseFiles関数（html/templateパッケージ）
		テンプレート template1.html を読み込む
		ParseFiles関数は、Template構造体のポインタを返す
	Must関数（html/templateパッケージ）
		テンプレートファイルの読み込みに失敗した場合パニック発生させて
		プログラムを終了する
		Must関数は Template構造体のポインタを返す
*/
var tmpl = template.Must(template.ParseFiles("othello.html"))

/* HTTPリクエストを処理する serverHTML構造体のメソッド定義
	メソッド名：ServeHTTP
	レシーバー：serverHTML構造体のポインタ
	引数：
		w（サーバーからクライアントへのレスポンス送信）http.ResponsWriter型
		r（クライアントから送信されるリクエストメッセージ）http.Resquest型へのポインタ
	処理：HTTPヘッダを設定して、Executeメソッドを実行する
		Executeメソッドは、テンプレートとserverHTML構造体を組み合わせた
		HTMLをクライアントに送信する
	戻り値：なし
*/
func (s Bord)checkPut(order string) [][]string {
	var data [][]string = [][]string{
		{"z", "z", "z", "z", "z", "z", "z", "z", "z", "z"}, 
		{"z", "", "", "", "", "", "", "", "", "z"},
		{"z", "", "", "", "", "", "", "", "", "z"}, 
		{"z", "", "", "", "", "", "", "", "", "z"}, 
		{"z", "", "", "", "w", "b", "", "", "", "z"}, 
		{"z", "", "", "", "b", "w", "", "", "", "z"}, 
		{"z", "", "", "", "", "", "", "", "", "z"}, 
		{"z", "", "", "", "", "", "", "", "", "z"}, 
		{"z", "", "", "", "", "", "", "", "", "z"}, 
		{"z", "z", "z", "z", "z", "z", "z", "z", "z", "z"}, 
	}
	var retData [][]string = [][]string{
		{"", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", ""}, 
		{"", "", "", "", "", "", "", ""}, 
		{"", "", "", "w", "b", "", "", ""}, 
		{"", "", "", "b", "w", "", "", ""}, 
		{"", "", "", "", "", "", "", ""}, 
		{"", "", "", "", "", "", "", ""}, 
		{"", "", "", "", "", "", "", ""}, 
	}
	//データコピーとフォーマット
	for i:=0; i<8; i++{
		for j:=0; j<8; j++{
			if s.Color[i][j] == "w" || s.Color[i][j] == "b" {
				data[i+1][j+1] = s.Color[i][j]
				retData[i][j] = s.Color[i][j]
			} else {
				data[i+1][j+1] = ""
				retData[i][j] = ""
			}
		}
	}
	var flg bool
	var opOrder string
	if order == "b" {
		opOrder = "w"
	} else {
		opOrder = "b"
	}

	for i:=0; i<8; i++{
		for j:=0; j<8; j++{
			flg = false
			if data[i+1][j+1] == "" {
				//おかれていない場合
				// 下
				if data[i+2][j+1] == opOrder {
					for k:=2; k<8; k++ {
						if data[i+k+1][j+1] == order {
							flg = true
						} else if data[i+k+1][j+1] != opOrder {
							break
						}
					}
				}
				// 上
				if data[i][j+1] == opOrder {
					for k:=2; k<8; k++ {
						if data[i-k+1][j+1] == order {
							flg = true
						} else if data[i-k+1][j+1] != opOrder {
							break
						}
					}
				}
				// 左
				if data[i+1][j] == opOrder {
					for k:=2; k<8; k++ {
						if data[i+1][j-k+1] == order {
							flg = true
						} else if data[i+1][j-k+1] != opOrder {
							break
						}
					}
				}
				// 右
				if data[i+1][j+2] == opOrder {
					for k:=2; k<8; k++ {
						if data[i+1][j+k+1] == order {
							flg = true
						} else if data[i+1][j+k+1] != opOrder {
							break
						}
					}
				}
				// 左上
				if data[i][j] == opOrder {
					for k:=2; k<8; k++ {
						if data[i-k+1][j-k+1] == order {
							flg = true
						} else if data[i-k+1][j-k+1] != opOrder {
							break
						}
					}
				}
				// 右上
				if data[i][j+2] == opOrder {
					for k:=2; k<8; k++ {
						if data[i-k+1][j+k+1] == order {
							flg = true
						} else if data[i-k+1][j+k+1] != opOrder {
							break
						}
					}
				}
				// 左下
				if data[i+2][j] == opOrder {
					for k:=2; k<8; k++ {
						if data[i+k+1][j-k+1] == order {
							flg = true
						} else if data[i+k+1][j-k+1] != opOrder {
							break
						}
					}
				}
				// 右下
				if data[i+2][j+2] == opOrder {
					for k:=2; k<8; k++ {
						if data[i+k+1][j+k+1] == order {
							flg = true
						} else if data[i+k+1][j+k+1] != opOrder {
							break
						}
					}
				}
			}
			if flg {
				retData[i][j] = "ok"
				flg = false
			}
		}
	}
	return retData
}

func (s Bord)putAndReverse(order string, y int, x int) [][]string {
	var data [][]string = [][]string{
		{"Z", "Z", "Z", "Z", "Z", "Z", "Z", "Z", "Z", "Z"},
		{"Z","", "", "", "", "", "", "", "", "Z"},
		{"Z","", "", "", "", "", "", "", "", "Z"}, 
		{"Z","", "", "", "", "", "", "", "", "Z"}, 
		{"Z","", "", "", "w", "b", "", "", "", "Z"}, 
		{"Z","", "", "", "b", "w", "", "", "", "Z"}, 
		{"Z","", "", "", "", "", "", "", "", "Z"}, 
		{"Z","", "", "", "", "", "", "", "", "Z"}, 
		{"Z","", "", "", "", "", "", "", "", "Z"}, 
		{"Z", "Z", "Z", "Z", "Z", "Z", "Z", "Z", "Z", "Z"},
	}
	//データコピーとフォーマット
	for i:=0; i<8; i++{
		for j:=0; j<8; j++{
			if s.Color[i][j] == "ok" {
				data[i+1][j+1] = ""
			} else {
				data[i+1][j+1] = s.Color[i][j]
			}
		}
	}

	x++
	y++

	var opOrder string
	if order == "b" {
		opOrder = "w"
	} else {
		opOrder = "b"
	}

	// 左
	for i:=1; i<8; i++ {
		if data[y][x-i] == order {
			for j:=i-1; j>0; j-- {
				data[y][x-j] = order
			}
			break;
		} else if data[y][x-i] != opOrder {
			break;
		}
	}
	// 上
	for i:=1; i<8; i++ {
		if data[y-i][x] == order {
			for j:=i-1; j>0; j-- {
				data[y-j][x] = order
			}
			break;
		} else if data[y-i][x] != opOrder {
			break;
		}
	}
	// 右
	for i:=1; i<8; i++ {
		if data[y][x+i] == order {
			for j:=i-1; j>0; j-- {
				data[y][x+j] = order
			}
			break;
		} else if data[y][x+i] != opOrder {
			break;
		}
	}
	// 下
	for i:=1; i<8; i++ {
		if data[y+i][x] == order {
			for j:=i-1; j>0; j-- {
				data[y+j][x] = order
			}
			break;
		} else if data[y+i][x] != opOrder {
			break;
		}
	}
	// 左上
	for i:=1; i<8; i++ {
		if data[y-i][x-i] == order {
			for j:=i-1; j>0; j-- {
				data[y-j][x-j] = order
			}
			break;
		} else if data[y-i][x-i] != opOrder {
			break;
		}
	}
	// 左下
	for i:=1; i<8; i++ {
		if data[y+i][x-i] == order {
			for j:=i-1; j>0; j-- {
				data[y+j][x-j] = order
			}
			break;
		} else if data[y+i][x-i] != opOrder {
			break;
		}
	}
	// 右上
	for i:=1; i<8; i++ {
		if data[y-i][x+i] == order {
			for j:=i-1; j>0; j-- {
				data[y-j][x+j] = order
			}
			break;
		} else if data[y-i][x+i] != opOrder {
			break;
		}
	}
	// 右下
	for i:=1; i<8; i++ {
		if data[y+i][x+i] == order {
			for j:=i-1; j>0; j-- {
				data[y+j][x+j] = order
			}
			break;
		} else if data[y+i][x+i] != opOrder {
			break;
		}
	}

	var retdata [][]string = [][]string{
		{"", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", ""}, 
		{"", "", "", "", "", "", "", ""}, 
		{"", "", "", "w", "b", "", "", ""}, 
		{"", "", "", "b", "w", "", "", ""}, 
		{"", "", "", "", "", "", "", ""}, 
		{"", "", "", "", "", "", "", ""}, 
		{"", "", "", "", "", "", "", ""}, 
	}

	for i:=0; i<8; i++{
		for j:=0; j<8; j++{
			retdata[i][j] = data[i+1][j+1]
		}
	}
	return retdata
}

func printBordLog(bord [][]string) (bool, bool, bool){
	var okFlg bool = false
	var bFlg bool = false
	var wFlg bool = false
	for _, v := range bord {
		fmt.Println("+--+--+--+--+--+--+--+--+")
		for _, w := range v {
			fmt.Print("|")
			if w == "" {
				fmt.Print("  ")
			} else {
				fmt.Printf("%2s", w)
				if w == "b" {
					bFlg = true
				} else if w == "w" {
					wFlg = true
				} else if w == "ok" {
					okFlg = true
				}
			}
		}
		fmt.Println("|")
	}
	fmt.Println("+--+--+--+--+--+--+--+--+")
	return okFlg, bFlg, wFlg
}

func (s *Bord) FormHandler(w http.ResponseWriter, r *http.Request) {
	// HTTPヘッダの設定
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//値取得 

	// 初回読み込み時以外はhtml側から盤面を取得する
	s.Index = r.FormValue("radio")
	fmt.Println(s.Index)
	var x, y int
	if s.Index != "" {
		var textFormName string;
		for i:=0; i<8; i++ {
			for j:=0; j<8; j++ {
				textFormName = "text"
				textFormName += strconv.Itoa(i) + strconv.Itoa(j)
				s.Color[i][j] = r.FormValue(textFormName)
			}
		}
		y, _ = strconv.Atoi(s.Index[0:1])
		x, _ = strconv.Atoi(s.Index[1:2])
		// // リバース
		s.Color = s.putAndReverse(s.Order, y, x)
		s.Color[y][x] = s.Order
	}

	// 順番交代
	if r.FormValue("order") == "b"{
		s.Order = "w"
	} else {
		s.Order = "b"
	}

	fmt.Println(s.Order)

	//おける場所のチェック
	s.Color = s.checkPut(s.Order)

	//ログ出力と存在の確認
	var okFlg bool = false
	var bFlg bool = false
	var wFlg bool = false
	var bCount int = 0
	var wCount int = 0
	okFlg, bFlg, wFlg = printBordLog(s.Color)

	//全て同じ色だったら終わる
	if !bFlg {
		// 全て黒
		s.WinFlg = "w"
	} else if !wFlg {
		// 全て白
		s.WinFlg = "b"
	} else if !okFlg {
		// おける場所なし
		s.Order = r.FormValue("order")			// 順序戻す
		s.Color = s.checkPut(s.Order)			// 再度ok確認
		okFlg, bFlg, wFlg = printBordLog(s.Color)	// ログ出力と存在の確認
		if !okFlg {
			for i:=0; i<8; i++ {
				for j:=0; j<8; j++ {
					if s.Color[i][j] == "b" {
						bCount ++
					} else if s.Color[i][j] == "w" {
						wCount ++
					}
				}
				if bCount > wCount {
					s.WinFlg = "b"
				} else if bCount < wCount {
					s.WinFlg = "w"
				} else {
					s.WinFlg = "e"
				}
			}
		}
	}

	// テンプレートとserverHTML構造体を組み合わせたHTMLをクライアントに送信
	result := tmpl.Execute(w, s)
	// クライアントへの送信にエラーがあるかを判断
	if result != nil {
		// エラーであれば panic 関数を使用して終了
		panic(result)
	}
}

// main関数定義
func main() {
	// ハンドラーごとにテンプレートに埋め込む値を設定
	// ハンドラー bord
	bord := &Bord{
		[][]string{
			{"", "", "", "", "", "", "", ""},
			{"", "", "", "", "", "", "", ""}, 
			{"", "", "", "", "", "", "", ""}, 
			{"", "", "", "w", "b", "", "", ""}, 
			{"", "", "", "b", "w", "", "", ""}, 
			{"", "", "", "", "", "", "", ""}, 
			{"", "", "", "", "", "", "", ""}, 
			{"", "", "", "", "", "", "", ""}, 
		},
		"b",
		"",
		"n",
	}
	// ハンドラー form_text の設定
	http.HandleFunc("/othello", bord.FormHandler)
	// Webサーバーを起動（ポート番号 8888）
	result := http.ListenAndServe(":8888", nil)
	if result != nil {
		fmt.Println(result)
	}
}