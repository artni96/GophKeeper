package user

//func LoginCallback(ctx context.Context, app *clientapp.App, loginEntity user.LoginRequest) error {
//	req := &pb.LoginRequest{}
//	req.SetUsername(loginEntity.Login)
//	req.SetPassword(loginEntity.Password)
//
//	token, err := app.UserClient.Login(ctx, req)
//	if err != nil {
//		return err
//	}
//	app.State.Token = token.GetToken()
//	return nil
//}
//
//func RegisterCallback(ctx context.Context, app *clientapp.App, loginEntity user.RegisterRequest) error {
//	req := &pb.UserCreateRequest{}
//	req.SetUsername(loginEntity.Login)
//	req.SetPassword(loginEntity.Password)
//
//	_, err := app.UserClient.CreateUser(ctx, req)
//	if err != nil {
//		return err
//	}
//	return nil
//}
