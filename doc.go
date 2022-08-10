// sentinel-golang document：https://sentinelguard.io/zh-cn/docs/golang/basic-api-usage.html
// example：https://github.com/alibaba/sentinel-golang/tree/master/example
//
// log document：https://github.com/alibaba/sentinel-golang/wiki/%E5%AE%9E%E6%97%B6%E7%9B%91%E6%8E%A7#%E7%A7%92%E7%BA%A7%E7%9B%91%E6%8E%A7%E6%97%A5%E5%BF%97
// 1	1529573107000	时间戳
// 2	2018-06-21 17:25:07	日期
// 3	foo-service	资源名称
// 4	10	这一秒通过的资源请求个数 (pass)
// 5	3601	这一秒资源被拦截的个数 (block)
// 6	10	这一秒完成调用的资源个数 (complete)，包括正常结束和异常结束的情况
// 7	0	这一秒资源的异常个数 (error)
// 8	27	资源平均响应时间（ms）
package main
