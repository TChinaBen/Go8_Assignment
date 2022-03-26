package error

type WrapError struct {
	Code string
	Msg  string
	Err error
}
// 实现了Error中string方法，就实现了该接口
func (e *WrapError)Error() string{
	return e.Msg
}
func (e *WrapError)Unwrap()error{
	return e.Err
}