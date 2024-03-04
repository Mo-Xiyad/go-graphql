package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.44

import (
	"context"
	"fmt"
	"server/graph"
	"server/pkg/model"
)

// ID is the resolver for the id field.
func (r *companyResolver) ID(ctx context.Context, obj *model.Company) (string, error) {
	return fmt.Sprint(obj.ID), nil
}

// Company returns graph.CompanyResolver implementation.
func (r *Resolver) Company() graph.CompanyResolver { return &companyResolver{r} }

type companyResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *companyResolver) Email(ctx context.Context, obj *model.Company) (string, error) {
	panic(fmt.Errorf("not implemented: Email - email"))
}
