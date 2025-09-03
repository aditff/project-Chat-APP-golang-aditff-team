// handler/user_handler.go
package handler

import (
	"context"
	"strings"
	"time"

	pb "project-Chat-APP-golang-aditff-user-service/proto"
	"project-Chat-APP-golang-aditff-user-service/service"
)

type UserHandler struct {
	Service *service.UserService
	pb.UnimplementedUserServiceServer
}

func (h *UserHandler) GetAllUsers(ctx context.Context, _ *pb.Empty) (*pb.UserList, error) {
	users, err := h.Service.GetAllUsers(ctx)
	if err != nil { return nil, err }

	out := make([]*pb.User, 0, len(users))
	for _, u := range users {
		var ls string
		if u.LastSeen != nil { ls = u.LastSeen.UTC().Format(time.RFC3339) }
		out = append(out, &pb.User{
			Id: u.ID, Name: u.Name, Email: u.Email, Online: u.Online, LastSeen: ls,
		})
	}
	return &pb.UserList{Users: out}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	u, err := h.Service.GetUser(ctx, req.Id)
	if err != nil { return nil, err }
	var ls string
	if u.LastSeen != nil { ls = u.LastSeen.UTC().Format(time.RFC3339) }
	return &pb.User{
		Id: u.ID, Name: u.Name, Email: u.Email, Online: u.Online, LastSeen: ls,
	}, nil
}

func (h *UserHandler) UpdateStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.Empty, error) {
	if err := h.Service.UpdateStatus(ctx, req.Id, req.Online); err != nil { return nil, err }
	return &pb.Empty{}, nil
}

// Optional: stream presence via Redis Pub/Sub
func (h *UserHandler) StreamPresence(_ *pb.Empty, stream pb.UserService_StreamPresenceServer) error {
	ctx := stream.Context()
	sub := h.Service.Redis.Subscribe(ctx, "presence")
	defer sub.Close()

	ch := sub.Channel()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-ch:
			// payload: "id,val,timestamp"
			var id, val, at string
			parts := strings.SplitN(msg.Payload, ",", 3)
			if len(parts) == 3 { id, val, at = parts[0], parts[1], parts[2] }
			online := val == "1"
			if err := stream.Send(&pb.PresenceEvent{UserId: id, Online: online, At: at}); err != nil {
				return err
			}
		}
	}
}
