## Fabric 2.2 司法存证智能合约分析
### 智能合约骨架
任何一个只能合约都由以下三部分组成
1) 合约对象`SimpleChaincode`，该数据结构内容字段为空
 ```go
type SimpleChaincode struct {

}
```
2) 合约对象必须实现两个方法
 ```go
// 合约安装成功后调用一次，疑惑再也不被调用
Init(stub ChaincodeStubInterface) pb.Response
// 每次向合约发起请求，会调用一次
Invoke(stub ChaincodeStubInterface) pb.Reponse
```
3) 合约对象必须被系统创建、加载
```go
func main() {
err := shim.Start(new(SimpleChaincode))
if err != nil {
    fmt.Printf("Error starting Simple chaincode: %s", err)
    }
}
```
其中第一部分、第三部分合约长得一样，第二部门不同，我们来详细分析
### Init 方法
Init方法系统仅在初始化阶段调用一次，一般用于系统初始化，该合约不需要系统初始化内容，所以实现方法为空。
```go
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Reponse {
    return shim.Success(nil)
}
```
### Invoke 方法
Invoke方法是业务逻辑的核心，后台系统发送字段要与智能合约处理字段保持一致。
1）客户端生成交易的时候需要指明调用方法，以及方法参数传递给智能合约，`function, args:=stub.GetFunctionAndParameters()`从交易报文中，提取本次交易需要调用的方法名字function，以及交易携带的参数。
2）系统根据不同的funciton进行内容处理
`upload_hash`：用于上传司法存证的hash值以及相关参数
`query_hash`：根据客户端传入的hash值，判断该hash值是否已经上链
`query_user_hash`：一个用户可能上传很多司法存证信息，也就是有很多hash值，该方法查询某个用户的全部
`upload_setting`：上传用户信息，包括姓名、邮箱、电话等
`query_setting`：根据用户姓名，查询用户信息
`query_all_setting`：查询区块链上全部用户信息
后续几个方法是用于设置发送邮件的模板的，这里不进行详细描述。
```go
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	switch function {
	case "upload_hash":
		return t.upload_hash(stub, args)
	case "query_hash":
		return t.query_hash(stub, args)
	case "query_user_hash":
		return t.query_user_hash(stub, args)
	case "upload_setting":
		return t.upload_setting(stub, args)
	case "query_setting":
		return t.query_setting(stub, args)
	case "query_all_setting":
		return t.query_all_setting(stub, args)
	case "upload_temp":
		return t.upload_temp(stub, args)
	case "upload_selected_temp":
		return t.upload_selected_temp(stub, args)
	case "update_temp":
		return t.update_temp(stub, args)
	case "delete_temp":
		return t.delete_temp(stub, args)
	case "query_temp":
		return t.query_temp(stub, args)
	case "query_selected_temp":
		return t.query_selected_temp(stub, args)
	}

	return shim.Error("Invalid invoke function Name. " + function)
}
```
3) 分析upload_hash函数
upload_hash函数用于上传司法存证数据信息，上传证书的数据结构为：`HashInfo`，包含4个字段，Hash上传证据的Hash值，Name上传者姓名、Date上传日期、Description上传内容描述。
```go
type HashInfo struct {
	Hash string `json:"hash,omitempty"`
	Name string `json:"name,omitempty"`
	Date string `json:"date,omitempty"`
	Description string `json:"description,omitempty"`
}
```
upload_hash 函数逻辑分析，看注释
```go
func (t *SimpleChaincode) upload_hash(stub shim.ChaincodeStubInterface, args []string) pb.Response {
//解析客户端传递过来的存证数据属性
	info := HashInfo{}
	argsBytes := []byte(args[0])

	err := json.Unmarshal(argsBytes, &info)
	if err != nil {
		return shim.Error(err.Error())
	}
// 调用GetState函数，查看合约数据库是否已经存在该Hash，也就是证据是否已经上传过
	value, err := stub.GetState("Hash#0#" + info.Hash)
	if err != nil {
		return shim.Error(err.Error())
	}
// 如果已经上传过，报错hash已经存在
	if value != nil {
		return shim.Error("hash existed")
	}
// 如果没有上传过，则以前缀"Hash#0#" + 证据Hash作为索引，证据属性作为值，存储到智能合约数据库。
	err = stub.PutState("Hash#0#" + info.Hash, argsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

// 同时以该用户信息 + hash作为联合索引，证据属性作为值，存储到智能合约数据库。
	err = stub.PutState("User#" + info.Name + "#0#" + info.Hash, argsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
```
4) 分析query_hash函数
   args[0]是客户端传来的，表示带查询的hash值
```go
func (t *SimpleChaincode) query_hash(stub shim.ChaincodeStubInterface, args []string) pb.Response {
// 从合约数据库中查看，是否存在这样的hash
	var hashInfo *HashInfo
	result, err := stub.GetState("Hash#0#" + args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
// 如果存在这样的hash值，把从数据库中返回的结果序列化后 返回给客户端；如果不存在也直接把空数据返回给客户端。
	if result != nil {
		err = json.Unmarshal(result, &hashInfo)
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	m,_ := json.Marshal(hashInfo)

	return shim.Success(m)
}
```
5) 分析query_user_hash函数
    该函数比较特别，需要利用到范围查询功能。`"User#" + info.Name + "#0#" + info.Hash`,智能合约数据库支持范围查询，是根据前缀匹配进行，根据尾部的ASCII范围进行选择，因此我们的查询方式为`beigin: "User#" + info.Name + "#0"`，`end: "User#" + info.Name + "#F"`, 匹配所有以`beigin: "User#" + info.Name + "#"`为前缀，下一个字符是从0到F的所有Key。
```go
func (t *SimpleChaincode) query_user_hash(stub shim.ChaincodeStubInterface, args []string) pb.Response {
//GetStateByRange 范围查询
	result, err := stub.GetStateByRange("User#" + args[0] + "#0", "User#" + args[0] + "#F")
	if err != nil {
		return shim.Error(err.Error())
	}
	var hashInfoList []HashInfo

	if result != nil {
		if result != nil {
			defer func() {
				result.Close()
			}()
			//通过迭代器拿到所有满足要求的数据
			for result.HasNext() {
				record, err := result.Next()
				if err != nil {
					return shim.Error(err.Error())
				}

				g := HashInfo{}
				err = json.Unmarshal(record.Value, &g)
				if err != nil {
					return shim.Error(err.Error())
				}

				hashInfoList = append(hashInfoList, g)
			}
		}
	}

// 查询到的结果返回区块链
	m,_ := json.Marshal(hashInfoList)

	return shim.Success(m)
}

其他方法分析思路类似
```