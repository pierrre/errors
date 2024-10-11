package errors

func CheckGlobalInit(err error, report func(error)) {
	checkGlobalInit(err, report)
}
