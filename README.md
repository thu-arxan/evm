# README

// 先用中文写好了。

## 1. Example

在项目的example中有一个简单样例用来描述如何使用EVM虚拟机及如何调用。

## 2. 要实现的几类接口

### 2.1. Account

```golang
type Account interface {
GetAddress() Address
GetBalance() uint64
AddBalance(balance uint64) error
SubBalance(balance uint64) error
GetCode() []byte
SetCode(code []byte)
// GetCodeHash return the hash of account code, please return [32]byte, and return [32]byte{0, ..., 0} if code is empty
GetCodeHash() []byte
GetNonce() uint64
SetNonce(nonce uint64)
}
```

要注意的有以下几点。

- Nonce: Nonce需要是自增的（有些案例中交易为了随机性会有一个交易Nonce，这两者不一定要等价，虽然以太坊中是等价的）。
- GetCodeHash：允许用户自己实现code的哈希函数，如果用户返回nil则会调用Keccak256函数，这一点和以太坊保持一致。

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

二者暂定接口如下,关于二者的定义，还需要进一步明确一下边界情况。

#### 2.3.1. DB

```golang
// Exist return if the account exist
// Note: if account is suicided, return true
Exist(address Address) bool
// GetStorage return a default account if unexist
GetAccount(address Address) Account
// Note: GetStorage return nil if key is not exist
GetStorage(address Address, key core.Word256) (value []byte)
```

#### 2.3.2. WriteBatch

```golang
SetStorage(address Address, key core.Word256, value []byte)
UpdateAccount(account Account) error
// Remove the account at address
RemoveAccount(address Address) error
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

- GetBlockHash：返回高度为num的区块哈希。
- CreateAddress：用户自定义的创建地址函数（对应CREATE指令），如果不想实现可以直接返回nil，EVM执行时会采取与以太坊相同的方式处理。
- Create2Address：用户自定义的创建地址函数（对应CREATE2指令），如果不想实现可以直接返回nil，EVM执行时会采取与以太坊相同的方式处理。
- NewAccount：根据一个地址返回默认的账户（请不要在DB里面也插入该账户，需要的时候EVM会调用DB的相关函数去插入）。
- BytesToAddress：将byte数组(长度一般为32位)解析为用户定义的Address。

## 3. 如何调用

参考example中的样例，文档待更新。
