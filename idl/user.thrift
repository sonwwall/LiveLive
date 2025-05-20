namespace go livelive.user

include "base.thrift"

struct RegisterReq{
    1:string username
    2:string password
    3:string email
    4:string mobile
}

struct RegisterResp{

    255:base.BaseResp baseResp
}

struct LoginReq{
    1:string username
    2:string password
}

struct LoginResp{

    255:base.BaseResp baseResp
}

struct UserInfoReq{
    1:string username

}

struct UserInfoResp{
    1:string username
    2:string email
    3:string mobile

    255:base.BaseResp baseResp
}

service UserService{
    RegisterResp Register(1:RegisterReq req)
    LoginResp Login(1:LoginReq req)
    UserInfoResp UserInfo(1:UserInfoReq req)
}