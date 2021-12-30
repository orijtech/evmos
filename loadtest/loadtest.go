package loadtest

var calls = []string{
	`{"jsonrpc":"2.0","method":"eth_call","params":[{"from":"0x3b7252d007059ffc82d16d022da3cbf9992d2f70", "to":"0xddd64b4712f7c8f1ace3c145c950339eddaf221d", "gas":"0x5208", "gasPrice":"0x55ae82600", "value":"0x16345785d8a0000", "data": "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"}, "0x0"],"id":1}`,
	`{"jsonrpc":"2.0","id":42,"method":"web3_clientVersion","params":[]}`,
	`{"jsonrpc":"2.0","id":42,"method":"web3_sha3","params":["0xdeadbeeeeeef"]}`,
	`{"jsonrpc":"2.0","method":"miner_start","params":["0x1"],"id":1}`,
	`{"jsonrpc":"2.0","method":"debug_traceBlockByNumber","params":["0xe", {"tracer": "{data: [], fault: function(log) {}, step: function(log) { if(log.op.toString() == \"CALL\") this.data.push(log.stack.peek(0)); }, result: function() { return this.data; }}"}],"id":1}`,
	`{"jsonrpc":"2.0","method":"txpool_content","params":[],"id":1}`,
	`{"jsonrpc":"2.0","method":"txpool_inspect","params":[],"id":1}`,
	`{"jsonrpc":"2.0","method":"txpool_status","params":[],"id":1}`,
}
