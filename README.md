# docment


## architecture design

通过mvc进行管理
routers：router request
services：business process 
dao：sql process
model：entity 

## function

### User

#### log in & register


注册通过发放code方式进行注册

#### user

基本属性

### article

文章被分类所管理
状态：草稿/发布/待解决/已解决/已关闭
分类：场景设计，项目答疑，qa板块，讨论板块

根据不同的分类也会有对应的状态，qa板块额外有待解决/已解决/已关闭


### comment

TODO: 中台

### system

发放邀请码code

