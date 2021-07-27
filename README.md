# work
У etherscan.io есть API, позволяющий получить информацию о транзакциях в блоке в сети ethereum, по номеру блока:
https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag=0x10d4f&boolean=true (tag - номер блока в 16 системе)
В этом методе, помимо прочего, возвращается список транзакций в блоке (result.transactions[]), для каждой транзакции описаны адрес отправителя, адрес получателя и сумма (result.transactions[].from, result.transactions[].to, result.transactions[].value).Есть метод, который возвращает номер последнего блока в сети: 
https://api.etherscan.io/api?module=proxy&action=eth_blockNumberНапиши программу, которая выдаст адрес, баланс которого изменился больше остальных (по абсолютной величине) за последние 100 блоков.
