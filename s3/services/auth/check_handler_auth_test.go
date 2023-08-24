package auth

//func TestV2CheckRequestAuthType(t *testing.T) {
//	var aSys AuthSys
//	aSys.Init()
//	req := testsign.MustNewSignedV2Request("GET", "http://127.0.0.1:9000", 0, nil, t)
//	_, _, err := aSys.CheckRequestAuthTypeCredential(context.Background(), req, s3action.ListAllMyBucketsAction, "test", "testobject")
//	fmt.Println(responses.GetAPIError(err))
//}
//func TestV4CheckRequestAuthType(t *testing.T) {
//	var aSys AuthSys
//	aSys.Init()
//	req := testsign.MustNewSignedV4Request("GET", "http://127.0.0.1:9000", 0, nil, "test", "test", "s3", t)
//	_, _, err := aSys.CheckRequestAuthTypeCredential(context.Background(), req, s3action.ListAllMyBucketsAction, "test", "testobject")
//	fmt.Println(responses.GetAPIError(err))
//}
