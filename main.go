package main

import "github.com/huytd/mss/source"

func main() {
	print("Processing... ")
	url := source.GetURL("http://www.nhaccuatui.com/bai-hat/yeu-khac-viet.0vxsibzy25Tx.html")
	print(url)
	print("\n\n")
	print("Processing... ")
	curl := source.GetURL("http://chiasenhac.vn/mp3/vietnam/v-pop/gui-anh-xa-nho~bich-phuong~tsvt3v3dqfw2wm.html")
	print(curl)
}
