package pb_handler

import (
	"context"
	"fmt"
	"time"

	pb "github.com/demkowo/users/internal/generated"
	model "github.com/demkowo/users/internal/models"
	service "github.com/demkowo/users/internal/services"
	"github.com/demkowo/utils/resp"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UsersServer struct {
	Service service.Users
	pb.UsersServer
}

func (h *UsersServer) Add(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	log.Trace("Add user via gRPC")

	var clubs []model.Club
	for _, clubName := range req.Clubs {
		clubs = append(clubs, model.Club{
			ID:   uuid.New(),
			Name: clubName,
		})
	}

	user := &model.User{
		ID:       uuid.New(),
		Nickname: req.Nickname,
		Img:      req.Img,
		Country:  req.Country,
		City:     req.City,
		Clubs:    clubs,
		Created:  time.Now(),
		Updated:  time.Now(),
	}

	if err := h.Service.Add(ctx, user); err != nil {
		log.Errorf("Failed to add user: %v", err)
		return nil, toGRPCError(err)
	}

	return &pb.AddUserResponse{
		User: toProtoUser(user),
	}, nil
}

func (h *UsersServer) Delete(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	log.Trace("Delete user via gRPC")

	if err := h.Service.Delete(ctx, req.GetUserId()); err != nil {
		log.Errorf("Failed to delete user: %v", err)
		return nil, toGRPCError(err)
	}

	return &pb.DeleteUserResponse{}, nil
}

func (h *UsersServer) Find(ctx context.Context, req *pb.FindUsersRequest) (*pb.FindUsersResponse, error) {
	log.Trace("Find all users via gRPC")

	usersFound, err := h.Service.Find(ctx)
	if err != nil {
		log.Errorf("Failed to find users: %v", err)
		return nil, toGRPCError(err)
	}

	return &pb.FindUsersResponse{
		Users: toProtoUsers(usersFound),
	}, nil
}

func (h *UsersServer) GetAvatarByNickname(ctx context.Context, req *pb.GetAvatarByNicknameRequest) (*pb.GetAvatarByNicknameResponse, error) {
	log.Trace("GetAvatarByNickname via gRPC")

	avatar, err := h.Service.GetAvatarByNickname(ctx, req.GetNickname())
	if err != nil {
		log.Errorf("Failed to get avatar: %v", err)
		return nil, toGRPCError(err)
	}

	return &pb.GetAvatarByNicknameResponse{Avatar: avatar}, nil
}

func (h *UsersServer) GetById(ctx context.Context, req *pb.GetByIdRequest) (*pb.GetByIdResponse, error) {
	log.Trace("GetById via gRPC")

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		log.Errorf("Invalid user ID: %v", err)
		return nil, err
	}

	user, e := h.Service.GetByID(ctx, userID)
	if e != nil {
		log.Errorf("Failed to get user by ID: %v", e)
		return nil, toGRPCError(e)
	}

	return &pb.GetByIdResponse{
		User: toProtoUser(user),
	}, nil
}

func (h *UsersServer) List(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	log.Trace("List users via gRPC")

	limit := req.GetLimit()
	offset := req.GetOffset()

	found, e := h.Service.List(ctx, limit, offset)
	if e != nil {
		log.Errorf("Failed to list users: %v", e)
		return nil, toGRPCError(e)
	}

	return &pb.ListUsersResponse{
		Users: toProtoUsers(found),
	}, nil
}

func (h *UsersServer) Update(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	log.Trace("Update user via gRPC")

	fmt.Println(req.City)
	fmt.Println(req.Clubs)
	fmt.Println(req.Country)
	fmt.Println(req.UserId)

	uid, err := uuid.Parse(req.GetUserId())
	if err != nil {
		log.Errorf("Invalid user ID: %v", err)
		return nil, err
	}

	var clubs []model.Club
	for _, name := range req.Clubs {
		clubs = append(clubs, model.Club{
			ID:   uuid.New(),
			Name: name,
		})
	}

	user := &model.User{
		ID:      uid,
		Country: req.GetCountry(),
		City:    req.GetCity(),
		Clubs:   clubs,
		Updated: time.Now(),
	}

	if e := h.Service.Update(ctx, user); e != nil {
		log.Errorf("Failed to update user: %v", e)
		return nil, toGRPCError(e)
	}

	return &pb.UpdateUserResponse{}, nil
}

func (h *UsersServer) UpdateImg(ctx context.Context, req *pb.UpdateImgRequest) (*pb.UpdateImgResponse, error) {
	log.Trace("Update user image via gRPC")

	uid, err := uuid.Parse(req.GetUserId())
	if err != nil {
		log.Errorf("Invalid user ID: %v", err)
		return nil, err
	}

	if e := h.Service.UpdateImg(ctx, uid, req.GetImg()); e != nil {
		log.Errorf("Failed to update user image: %v", e)
		return nil, toGRPCError(e)
	}

	return &pb.UpdateImgResponse{}, nil
}

func toProtoUser(u *model.User) *pb.User {
	if u == nil {
		return nil
	}
	protoClubs := make([]*pb.Club, 0, len(u.Clubs))
	for _, c := range u.Clubs {
		protoClubs = append(protoClubs, &pb.Club{
			Id:   c.ID.String(),
			Name: c.Name,
		})
	}

	return &pb.User{
		Id:       u.ID.String(),
		Nickname: u.Nickname,
		Img:      u.Img,
		Country:  u.Country,
		City:     u.City,
		Clubs:    protoClubs,
		Created:  timestamppb.New(u.Created),
		Updated:  timestamppb.New(u.Updated),
		Deleted:  u.Deleted,
	}
}

func toProtoUsers(us []model.User) []*pb.User {
	res := make([]*pb.User, 0, len(us))
	for _, u := range us {
		user := u
		res = append(res, toProtoUser(&user))
	}
	return res
}

func toGRPCError(err *resp.Err) error {
	if err == nil {
		return nil
	}

	return status.Error(codes.Code(err.Code), err.Error)
}
