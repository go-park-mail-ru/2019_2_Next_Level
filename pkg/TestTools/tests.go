package TestTools

type Params interface{}

type TestStruct struct {
	Input []Params
	Expected []Params
	MockParams []Params
}

type TestStructMap struct {
	Input map[string]Params
	Expected map[string]Params
	MockParams map[string]Params
}

func NewTestStructMap(input map[string]Params, expected map[string]Params, mockParams map[string]Params) *TestStructMap {
	return &TestStructMap{Input: input, Expected: expected, MockParams: mockParams}
}

func NewTestStruct(input []Params, expected []Params, mockParams []Params) *TestStruct {
	return &TestStruct{Input: input, Expected: expected, MockParams: mockParams}
}



type ExpectedFunc func(params ...Params)

func RunTesting (tests []TestStruct, f ExpectedFunc, testFunc func(TestStruct)) {
	for _, test := range tests{
		f(test.MockParams...)
		testFunc(test)
	}
}

func RunTestingMapped (tests []TestStructMap, f func(map[string]Params), testFunc func(TestStructMap)) {
	for _, test := range tests{
		f(test.MockParams)
		testFunc(test)
	}
}

