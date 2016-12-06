# go-lightrpc
简单的封装了 net/http ，提供 jsonrpc 快速服务，协议描述如下：

### 协议说明


**【后台服务入口】 提供统一标准的 输入/输出 协议为客户端提供服务，协议说明如下**

* 服务URL
* 测试系统: http://123.178.27.74/pet-hub/request
* 正式系统: http://www.52pet.net/pet-hub/request
* 输入参数
* 参数名称：body

**例如：**
<pre><code>{
  "service":"service.uri.pet_sso",
  "method":"login",
  "channel":"1",
  "sn":"全局唯一的UUID"
  "params":{"username":"cc","password":"123"}
}
</code></pre>

**参数格式：**

<p>以用户登录的请求参数为例，必填项如下

<pre><code>{
    "service":"service.uri.pet_sso",
    "method":"login",
    "channel":"1",
    "sn":"全局唯一的UUID"
    "params":{
        "username":"cc",
        "password":"123"
    }
}</pre></code>

**<p>参数说明（只包含了必填项）：**

* service: 业务模块的注册名，下文会给出业务模块的注册表；
* method: 具体的业务方法；
* channel: 请求来源渠道
* sn: <font color='red' >请求流水号，要求全局唯一，建议使用 UUID</font>
* params: 供业务方法使用的参数，具体参数内容，应参考业务模块注册表，其中有部分模块可以不需要此参数；

**<p>返回值说明：**

统一格式的返回，其中 success 标识请求是否成功，返回 true 则 entity 为 object，object 格式由业务模块定义
返回 false 则表示请求异常，其中 entity 为 exception 信息

<pre><code>{
    "sn":"请求流水号"
    "success": true / false, 
    "entity": object / exception
}</pre></code>
