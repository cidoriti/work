//У etherscan.io есть API, позволяющий получить информацию о транзакциях в блоке в сети ethereum, по номеру блока:
//https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag=0x10d4f&boolean=true (tag - номер блока в 16 системе)
//В этом методе, помимо прочего, возвращается список транзакций в блоке (result.transactions[]), для каждой транзакции
//описаны адрес отправителя, адрес получателя и сумма (result.transactions[].from, result.transactions[].to, result.transactions[].value).
//Есть метод, который возвращает номер последнего блока в сети:
//https://api.etherscan.io/api?module=proxy&action=eth_blockNumberНапиши программу, которая выдаст адрес, баланс которого изменился больше
//остальных (по абсолютной величине) за последние 100 блоков.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/time/rate"
)

type BlockNumber struct {
	Jsonrpc string `json:"-"`
	ID      int    `json:"-"`
	Result  string `json:"result"`
}

type Transaction struct {
	Jsonrpc string `json:"-"`
	ID      int    `json:"id"`
	Result  Result `json:"result"`
}
type Result struct {
	Difficulty       string         `json:"-"`
	ExtraData        string         `json:"-"`
	GasLimit         string         `json:"-"`
	GasUsed          string         `json:"-"`
	Hash             string         `json:"-"`
	LogsBloom        string         `json:"-"`
	Miner            string         `json:"-"`
	MixHash          string         `json:"-"`
	Nonce            string         `json:"-"`
	Number           string         `json:"-"`
	ParentHash       string         `json:"-"`
	ReceiptsRoot     string         `json:"-"`
	Sha3Uncles       string         `json:"-"`
	Size             string         `json:"-"`
	StateRoot        string         `json:"-"`
	Timestamp        string         `json:"-"`
	TotalDifficulty  string         `json:"-"`
	Transactions     []Transactions `json:"transactions"`
	TransactionsRoot string         `json:"-"`
	Uncles           []interface{}  `json:"-"`
}
type Transactions struct {
	BlockHash        string `json:"-"`
	BlockNumber      string `json:"-"`
	From             string `json:"from"`
	Gas              string `json:"-"`
	GasPrice         string `json:"-"`
	Hash             string `json:"-"`
	Input            string `json:"-"`
	Nonce            string `json:"-"`
	To               string `json:"to"`
	TransactionIndex string `json:"-"`
	Value            string `json:"value"`
	Type             string `json:"-`
	V                string `json:"-"`
	R                string `json:"-"`
	S                string `json:"-"`
}

type RLHTTPClient struct {
	client      *http.Client
	Ratelimiter *rate.Limiter
}

func (c *RLHTTPClient) Do(req *http.Request) (*http.Response, error) {
	ctx := context.Background()
	err := c.Ratelimiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewClient(rl *rate.Limiter) *RLHTTPClient {
	c := &RLHTTPClient{
		client:      http.DefaultClient,
		Ratelimiter: rl,
	}
	return c
}

func ConvertInt(val string, base, toBase int) (string, error) {
	i, err := strconv.ParseInt(val, base, 64)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(i, toBase), nil
}

func CreateUrl(chislo int) string {
	chislo16, _ := ConvertInt(strconv.Itoa(chislo), 10, 16)
	URL := "https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&0x" + chislo16 + "=0x10d4f&boolean=true"
	fmt.Println(URL)
	return URL
}

func (t *Transaction) summing() uint {
	var sl []int
	var sum int
	for i := range t.Result.Transactions {
		value, err := ConvertInt(t.Result.Transactions[i].Value[2:], 16, 10)
		if err != nil {
			continue
		} else {
			chislo, err := strconv.Atoi(value)
			if err != nil {
				panic(err)
			}
			sl = append(sl, chislo)
		}

	}
	for i := range sl {
		sum = sum + sl[i]
	}
	return uint(sum)
}

func main() {
	mapa := make(map[int]uint) // я не знаю что какой конкретно адрес должен быть в ключе, по этому пока ID
	BN := new(BlockNumber)
	rl := rate.NewLimiter(rate.Every(5*time.Second), 1)
	API_LAST_TRANSACTION := "https://api.etherscan.io/api?module=proxy&action=eth_blockNumber"
	c := NewClient(rl)
	req, _ := http.NewRequest("GET", API_LAST_TRANSACTION, nil)
	resp, _ := c.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &BN)
	chislo10, err := ConvertInt(BN.Result[2:], 16, 10)
	if err != nil {
		panic(err)
	}
	chisloInt, err := strconv.Atoi(chislo10)
	if err != nil {
		panic(err)
	}
	for i := 0; i != 5; i++ {
		T := new(Transaction)
		var sum uint
		resp, err := c.Do(req)
		reqb, _ := http.NewRequest("GET", CreateUrl(chisloInt-i), nil)
		respb, err := c.Do(reqb)
		bodyPerson, err := ioutil.ReadAll(respb.Body)
		json.Unmarshal(bodyPerson, &T)
		sum = T.summing()
		mapa[T.ID] = sum
		if err != nil {
			panic(err)
		}
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println(resp.StatusCode)
			return
		}
		if resp.StatusCode == 429 {
			fmt.Printf("Rate limit reached after %d requests", i)
			return
		}
	}
	var max uint
	var otvet int
	for key, value := range mapa {
		if value > max {
			otvet = key
		} else {
			continue
		}
		fmt.Println(otvet)
	}
}
