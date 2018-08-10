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
// util.TCF{
// 	Try: func() {
// 		log.Println("I tried")
// 		util.Throw("Oh,...sh...")
// 	},
// 	Catch: func(e util.Exception) {
// 		log.Printf("Caught %v\n", e)
// 	},
// 	Finally: func() {
// 		log.Println("Finally...")
// 	},
// }.Do()
