package util

// TCF struct
type TCF struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

// Exception interface
type Exception interface{}

// Throw Exception
func Throw(up Exception) {
	panic(up)
}

// Do block
func (tcf TCF) Do() {
	if tcf.Finally != nil {

		defer tcf.Finally()
	}
	if tcf.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				tcf.Catch(r)
			}
		}()
	}
	tcf.Try()
}

////=========== Use code template ===========////
// TCF{
// 	Try: func() {
// 		fmt.Println("I tried")
// 		Throw("Oh,...sh...")
// 	},
// 	Catch: func(e Exception) {
// 		fmt.Printf("Caught %v\n", e)
// 	},
// 	Finally: func() {
// 		fmt.Println("Finally...")
// 	},
// }.Do()
