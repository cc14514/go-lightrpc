# go-lightrpc

### 协议说明


**提供统一标准的 输入/输出 协议为客户端提供服务，协议说明如下**

* 输入参数
* 参数名称：body

**例如：**
<pre><code>
http://localhost:8080/api/?body={"service":"sso","method":"login","sn":"UUID""params":{"username":"cc","password":"123"}}
</code></pre>

**参数格式：**

<p>以用户登录的请求参数为例，必填项如下

<pre><code>{
    "service":"sso",
    "method":"login",
    "sn":"全局唯一的UUID"
    "params":{
        "username":"cc",
        "password":"123"
    }
}</pre></code>

**<p>参数说明（只包含了必填项）：**

* service: 业务模块的注册名，下文会给出业务模块的注册表；
* method: 具体的业务方法；
* sn: <font color='red' >请求流水号，要求全局唯一，建议使用 UUID</font>
* params: 业务方法对应的参数；

**<p>返回值说明：**

统一格式的返回，其中 success 标识请求是否成功，返回 true 则 entity 为 object，object 格式由业务模块定义
返回 false 则表示请求异常，其中 entity 为 exception 信息

<p>成功：
<pre><code>{
    "sn":"请求流水号"
    "success": true , 
    "entity": object 
}</pre></code>
<p>失败：
<pre><code>{
    "sn":"请求流水号"
    "success": false , 
    "entity": {
        "errCode":"错误码",
        "reason":"原因描述"
    } 
}</pre></code>

### 使用样例：
<a href="https://github.com/cc14514/go-lightrpc-example">https://github.com/cc14514/go-lightrpc-example</a>
