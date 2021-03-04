package Middle

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

//错误封装
func CheckError(err error,str string){
	if err != nil {
		panic(fmt.Sprintf("err: %s",str))
	}
}

//Gin当中的异常捕获
func ErrorRecover() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func(){
			if err := recover() ; err != nil {
				context.AbortWithStatusJSON(400,gin.H{"error":err})
			}
		}()
		context.Next()
	}
}