package postgres

import (
	"reflect"
	"testing"

	pb "github.com/project1/post-service/genproto"
)

func TestPostRepo_Create(t *testing.T) {
	tests := []struct {
		name  string
		input *pb.Post

		want    *pb.Post
		wantErr bool
	}{
		{
			name: "success case",
			input: &pb.Post{
				Name:        "test",
				Description: "dfs",
				Medias:      nil,
			},
			want: &pb.Post{
				Name:        "test",
				Description: "dfs",
				Medias:      nil,
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := repo.CreatePost(tc.input)
			if err != nil {
				t.Fatalf("%s: expected: %v, got: %v", tc.name, tc.wantErr, err)
			}

			got.Id = ""
			if !reflect.DeepEqual(&tc.want, got) {
				t.Fatalf("%s: expected: %v, got: %v", tc.name, tc.want, got)
			}
		})
	}

}
