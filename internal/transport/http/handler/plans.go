package handler

import (
	"context"
	"math/big"
	"smartDriver/internal/db"
	"smartDriver/pkg/log"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

type createPlanIn struct {
	Body struct {
		Name          string  `json:"name" maxLength:"50" doc:"Plan name"`
		Cost          float64 `json:"cost" doc:"Plan cost"`
		EmployeeLimit int32   `json:"employee_limit" doc:"Plan employee limit"`
	}
}

type planOut struct {
	Body struct {
		ID            int64   `json:"id" doc:"Plan ID"`
		Name          string  `json:"name" maxLength:"50" doc:"Plan name"`
		Cost          float64 `json:"cost" doc:"Plan cost"`
		EmployeeLimit int32   `json:"employee_limit" doc:"Plan employee limit"`
	}
}

func CreatePlan(ctx context.Context, in *createPlanIn) (*planOut, error) {
	plan, err := db.Pool.CreatePlan(ctx, db.CreatePlanParams{
		Name: in.Body.Name,
		// FIXME: this shit write in db without decimal part LOL
		Cost: pgtype.Numeric{
			Int:   big.NewInt(int64(in.Body.Cost)),
			Exp:   0,
			NaN:   false,
			Valid: true,
		},
		EmployeeLimit: in.Body.EmployeeLimit,
	})
	if err != nil {
		log.SugaredLogger.Errorf("failed to create plan: %v", err)
		return nil, huma.Error500InternalServerError("failed to create plan", err)
	}

	cost, err := plan.Cost.Float64Value()
	if err != nil {
		log.SugaredLogger.Errorf("failed to get cost float64 value: %v", err)
		return nil, huma.Error500InternalServerError("failed to get cost float64 value", err)
	}

	var resp planOut
	resp.Body.Name = plan.Name
	resp.Body.Cost = cost.Float64
	resp.Body.EmployeeLimit = plan.EmployeeLimit

	return &resp, nil
}

func GetPlan(ctx context.Context, in *idPathIn) (*planOut, error) {
	plan, err := db.Pool.GetPlan(ctx, in.ID)
	if err != nil {
		log.SugaredLogger.Errorf("failed to get plan: %v", err)
		return nil, huma.Error500InternalServerError("failed to get plan", err)
	}

	cost, err := plan.Cost.Float64Value()
	if err != nil {
		log.SugaredLogger.Errorf("failed to get cost float64 value: %v", err)
		return nil, huma.Error500InternalServerError("failed to get cost float64 value", err)
	}

	var resp planOut
	resp.Body.Name = plan.Name
	resp.Body.Cost = cost.Float64
	resp.Body.EmployeeLimit = plan.EmployeeLimit

	return &resp, nil
}

type updatePlanIn struct {
	ID   int64 `path:"id" json:"id" doc:"Plan ID"`
	Body struct {
		Name          string  `json:"name" maxLength:"50" doc:"Plan name"`
		Cost          float64 `json:"cost" doc:"Plan cost"`
		EmployeeLimit int32   `json:"employee_limit" doc:"Plan employee limit"`
	}
}

func UpdatePlan(ctx context.Context, in *updatePlanIn) (*planOut, error) {
	plan, err := db.Pool.UpdatePlan(ctx, db.UpdatePlanParams{
		ID:   in.ID,
		Name: in.Body.Name,
		Cost: pgtype.Numeric{
			Int:   big.NewInt(int64(in.Body.Cost)),
			Exp:   0,
			NaN:   false,
			Valid: true,
		},
		EmployeeLimit: in.Body.EmployeeLimit,
	})
	if err != nil {
		log.SugaredLogger.Errorf("failed to update plan: %v", err)
		return nil, huma.Error500InternalServerError("failed to update plan", err)
	}

	cost, err := plan.Cost.Float64Value()
	if err != nil {
		log.SugaredLogger.Errorf("failed to get cost float64 value: %v", err)
		return nil, huma.Error500InternalServerError("failed to get cost float64 value", err)
	}

	var resp planOut
	resp.Body.Name = plan.Name
	resp.Body.Cost = cost.Float64
	resp.Body.EmployeeLimit = plan.EmployeeLimit

	return &resp, nil
}

func DeletePlan(ctx context.Context, in *idPathIn) (*successOut, error) {
	if err := db.Pool.DeletePlan(ctx, in.ID); err != nil {
		log.SugaredLogger.Errorf("failed to delete plan: %v", err)
		return nil, huma.Error500InternalServerError("failed to delete plan", err)
	}

	var resp successOut
	resp.Body.Success = true

	return &resp, nil
}

type listPlansOut struct {
	Body struct {
		Plans []db.Plan
	}
}

func ListPlans(ctx context.Context, in *listIn) (*listPlansOut, error) {
	plans, err := db.Pool.ListPlans(ctx)
	if err != nil {
		log.SugaredLogger.Errorf("failed to list plans: %v", err)
		return nil, huma.Error500InternalServerError("failed to list plans", err)
	}

	var resp listPlansOut
	resp.Body.Plans = plans

	return &resp, nil
}
