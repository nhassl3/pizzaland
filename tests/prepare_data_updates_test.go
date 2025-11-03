package tests

import (
	"testing"

	pizzalndv1 "github.com/nhassl3/pizzaland/api/generated/go/pizzaland"
	"github.com/nhassl3/pizzaland/tests/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestPrepareDataUpdates(t *testing.T) {
	ctx, st := suite.NewSuite(t)

	ListResp, err := st.PizzaClient.Update(ctx, &pizzalndv1.UpdateRequest{
		Identifier: &pizzalndv1.UpdateRequest_Id{
			Id: 1,
		},
		Name:  wrapperspb.String("Клевая"),
		Price: wrapperspb.Float(1999.99),
		TypeDough: []pizzalndv1.TypeDough{
			pizzalndv1.TypeDough_THIN_DOUGH,
		},
	})
	require.NoError(t, err)
	assert.True(t, ListResp.GetSuccess())
}
