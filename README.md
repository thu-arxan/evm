# 文档

evm实现了以太坊黄皮书中设计的虚拟机，可用于运行solidity编写的智能合约。

## 1.项目结构

```shell
|
|- abi          //实现了外部调用智能合约的格式转换工具
|- core         //实现了一些接口
|- crypto       //密码学相关函数实现
|- db           //数据库实现
|- errors       //错误码定义
|- example      //示例，可参考example/README.md
|- gas          //汇编代码消耗的gas定义
|- precompile   //本地合约，golang实现
|- rlp          //编解码算法
|- tests        //测试
|- util         //公共函数
|- cache.go     //缓存，加速数据库操作
|- context.go   //evm运行上下文
|- evm.go       //汇编实现
|- interface.go //接口定义
|- opcodes.go   //汇编表
|- memory.go    //evm存储实现
|- stack.go			//evm存储实现
```

## 2. 要实现的几类接口

### 2.1. Account

```golang
type Account interface {
// Getter of account address / code / balance
// Setter of account code / balance
GetAddress() Address
GetBalance() uint64
AddBalance(balance uint64) error
SubBalance(balance uint64) error
GetCode() []byte
SetCode(code []byte)
// GetCodeHash return the hash of account code, please return [32]byte, and // // return [32]byte{0, ..., 0} if code is empty
GetCodeHash() []byte
GetNonce() uint64
SetNonce(nonce uint64)
// Suicide will suicide an account
Suicide()
HasSuicide() bool
}
```

要注意的有以下几点。

- Nonce: Nonce需要是自增的（有些案例中交易为了随机性会有一个交易Nonce，这两者不一定要等价，虽然以太坊中是等价的）。
- GetCodeHash：允许用户自己实现code的哈希函数，如果用户返回nil则会调用Keccak256函数，这一点和以太坊保持一致。
- 对Balance的相关操作要注意溢出的错误处理。

### 2.2. Address

```golang
type Address interface {
    Bytes() []byte
}
```

用户需要为自己所实现的地址定义Bytes接口以转换为bytes供EVM所使用。

- 如果序列化长度为32，则使用时不进行处理。
- 如果序列化长度小于32，则在左边补0至32位。
- 如果序列化长度大于32，则忽视左侧的部分并缩短至32位处理。

但是，需要注意的是，***在EVM中如果地址长度超过20位，运行时将可能并不能得到预期的结果，所以请不要使用有效信息超过20位的地址***。

下面会详细阐释一下原因。

虽然在以太坊中，栈、内存等都是32位的机器，看起来能够支撑32位以内的地址，但是以下面的代码为例。

```js
function info() public view returns (address, uint) {
    return (msg.sender, balance);
}
```

其生成的汇编代码并不会老老实实的将传进来的sender地址返回，有可能会通过一个PUSH20指令将地址截断放到栈中并使用(可能是为了gas消耗角度考虑)，这样一来地址的有效信息就被截断了，所以会导致信息丢失。

### 2.3. DB & WriteBatch

底层数据库存储每个账户及其拥有的kv storage。当DB作为交易执行的cache使用时，交易完成后需要把更新及删除的账户状态写入底层数据库，这个过程使用batch执行来加速。

#### 2.3.1. DB

```golang
// Exist return if the account exist
// Note: if account is suicided, return true
Exist(address Address) bool
// GetStorage return a default account if unexist
GetAccount(address Address) Account
// Note: GetStorage return nil if key is not exist
GetStorage(address Address, key []byte) (value []byte)
// if db is used as cache, updated and removed account need to be synced to 
// database by writeBatch once execution finished
NewWriteBatch() WriteBatch
```

如果一个账户在交易执行过程中执行了selfdestruct汇编指令，会被标记为suicided，这说明该账户曾经存在过，在底层数据库中有相应的kv存储。由于Exist通常（ethereum里）在创建新账户的时候被调用，由于底层数据库和cache中还没有清除该账户的信息，创建新账户不需要过多操作，可以不用收取多余的gas。

#### 2.3.2. WriteBatch

```golang
SetStorage(address Address, key []byte, value []byte)
// Note: db should delete all storages if an account suicide
UpdateAccount(account Account) error
AddLog(log *Log)
```

### 2.4. Blockchain

```golang
GetBlockHash(num uint64) []byte
// CreateAddress will be called by CREATE Opcode
CreateAddress(caller Address, nonce uint64) Address
// Create2Address will be called by CREATE2 Opcode
Create2Address(caller Address, salt, code []byte) Address
// Note: NewAccount will create a default account in Blockchain service, but please do not append the account into db here
NewAccount(address Address) Account
// BytesToAddress provide a way convert bytes(normally [32]byte) to Address
BytesToAddress(bytes []byte) Address
```

- GetBlockHash：返回高度为num的区块哈希，注意num < block height 且 > blockheight - 257。
- CreateAddress：用户自定义的创建地址函数（对应CREATE指令），如果不想实现可以直接返回nil，EVM执行时会采取与以太坊相同的方式处理。
- Create2Address：用户自定义的创建地址函数（对应CREATE2指令），如果不想实现可以直接返回nil，EVM执行时会采取与以太坊相同的方式处理。
- NewAccount：根据一个地址返回默认的账户（请不要在DB里面也插入该账户，需要的时候EVM会调用DB的相关函数去插入）。
- BytesToAddress：将byte数组(长度一般为32位)解析为用户定义的Address。