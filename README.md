# gohessian

gohessian implements  Hessian   Protocol by golang

修改了包结构，便于直接使用 go get github.com/shanggl/gohessian 之后直接使用

增加直接使用Hessian序列化、反序列化方式

POJO branch 采用模拟java的POJO geter、setter方法来操作struct对象，去掉structs 包的依赖
