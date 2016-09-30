# glog level说明

0 -- 关闭所有V info日志
1 -- 返回error，但是不适合输出error或warning日志
2 -- info log，只在某些关键分支上输出
3 -- trace log，跟踪日志，通过跟踪日志可以比较完整地回溯一个请求的调用流程
4 -- debug log，输出协议日志的级别，以及用于调试阶段的日志

通常开到level2或者level3，在出错需要调试时，开到level4
