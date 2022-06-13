package postgres

import (
	"reflect"
	"testing"

	pb "github.com/project1/user-service/genproto"
)

func TestUserRepo_Create(t *testing.T) {
	tests := []struct {
		name    string
		input   *pb.User
		want    *pb.User
		wantErr bool
	}{
		{
			name: "success case",
			input: &pb.User{
				FirstName: "aaaaa",
				LastName:  "ccccc",
				Email:        "saaaaaa",
				Bio:          "assdasss",
				Address:      nil,
				Status:       "pppppppppp",
				CreatedAt: "",
				PhoneNumbers: []string{"7777777777", "8888888888"},
			},
			want: &pb.User{
				FirstName: "aaaaa",
				LastName:  "ccccc",
				Email:        "saaaaaa",
				Bio:          "assdasss",
				Address:      nil,
				Status:       "pppppppppp",
				CreatedAt: "",
				PhoneNumbers: []string{"7777777777", "8888888888"},
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := repo.CreateUser(tc.input)
			if err != nil {
				t.Fatalf("%s: expected: %v, got: %v", tc.name, tc.wantErr, err)
			}

			got.Id = ""
			got.Posts = nil
			got.Address = nil
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("%s: expected: %v, got: %v", tc.name, tc.want, got)
			}
		})
	}
}
