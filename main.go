package main

import (
	"flag"
	"fmt"
	stripe "github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/balance"
	"github.com/stripe/stripe-go/v72/balancetransaction"
	"github.com/stripe/stripe-go/v72/payout"
	"github.com/tomcraven/gotable"
	"os"
	"time"
)

func main(){
	var Key string
	flag.StringVar(&Key, "k", "", "Stripe 私钥!")
	flag.Parse()

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	stripe.Key = Key // "sk_test_51JQAN**********"

	b, err := balance.Get(nil)
	if err == nil {
		Available := float32(b.Available[0].Value)
		Pending := float32(b.Pending[0].Value)
		fmt.Printf("钱包可用余额:%.2f,冻结余额:%.2f\n",Available,Pending)
	}

	params := &stripe.BalanceTransactionListParams{}
	params.Filters.AddFilter("limit", "", "10")
	i := balancetransaction.List(params)

	t := gotable.NewTable([]gotable.Column{
		gotable.NewColumn("TXN ID", 20),
		gotable.NewColumn("CH ID", 20),
		gotable.NewColumn("Amount", 20),
		gotable.NewColumn("Date/Time", 20),
	})

	for i.Next() {
		bt := i.BalanceTransaction()
		t.Push(bt.ID,bt.Source.ID,float32(bt.Amount),time.Unix(bt.Created, 0).Format("2006-01-02 15:04:05"))
	}

	t.Print()

	var Amount float32
	fmt.Println("是否打算提现,如果需要,请输入提现金额(不能大于可用余额),否则输入0.")
	fmt.Scanln(&Amount)

	if Amount == 0 {
		os.Exit(0)
	}else{
		Amount = Amount * 100
		params := &stripe.PayoutParams{
			Amount: stripe.Int64(int64(Amount)),
			Currency: stripe.String(string(stripe.CurrencyJPY)),
		}
		payout.New(params)
		fmt.Println("已经申请提现,但是提现是异步操作,你需要30分钟后才知道结果.")
	}
}
