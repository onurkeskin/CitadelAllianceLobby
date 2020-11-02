package service

import (
	"context"
	"net/http"

	apiModel "keon.com/CitadelAllianceLobbyServer/user-service-api/model"
	pb "keon.com/CitadelAllianceLobbyServer/user-service-api/proto/user"
)

//---- User Request API v0 ----

// A ListUsersResponse parameter model.
//
// Used as a response for seeing users.
//
// swagger:response listUsersResponse_v0
type ListUsersResponse_v0 struct {
	// in: body
	Users   apiModel.Users `json:"users"`
	Message string         `json:"message,omitempty"`
	Success bool           `json:"success"`
}

// A ErrorResponse parameter model.
//
// Used as a response for errors.
//
// swagger:response errorResponse_v0
type ErrorResponse_v0 struct {
	// in: body
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

func (resource *Resource) HandleListUsers_v0(w http.ResponseWriter, req *http.Request) {
	// // filter & pagination params
	// field := req.FormValue("field")
	// query := req.FormValue("q")
	// lastID := req.FormValue("last_id")
	// perPageStr := req.FormValue("per_page")
	// sort := req.FormValue("sort")

	// perPage, err := strconv.Atoi(perPageStr)
	// if err != nil {
	// 	perPage = 20
	// }

	// u := repo.FilterUsers(field, query, lastID, perPage, sort)
	// users := *u.(*Users)
	// if len(users) > 0 {
	// 	lastID = users[len(users)-1].ID.Hex()
	// }

	ctx := context.Background()
	conn, err := resource.GrpcHelper.getGRPCClient()
	defer conn.Close()
	if err != nil {
		resource.Render(w, req, http.StatusBadRequest, err.Error())
	}

	client := pb.NewUserServiceClient(conn)
	pbUsers, err := client.GetAllUsers(ctx, &pb.GetAllRequest{})
	if err != nil {
		resource.Render(w, req, http.StatusBadRequest, err.Error())
		return
	}

	users := apiModel.FromProtobufModels(pbUsers.Users)
	resource.Render(w, req, http.StatusOK, ListUsersResponse_v0{
		Users:   users,
		Message: "User list retrieved",
		Success: true,
	})
}
