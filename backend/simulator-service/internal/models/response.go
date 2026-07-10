package models

import (
	"github.com/NekNB/CyberNavigate/swagger/gen/simulator"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type StepMetaResponse struct {
	Id            *string             `json:"id,omitempty"`
	Actions       *[]simulator.Action `json:"actions,omitempty"`
	MaxTrust      *int                `json:"maxTrust,omitempty"`
	MinTrust      *int                `json:"minTrust,omitempty"`
	PreviosAnswer *openapi_types.UUID `json:"previosAnswer,omitempty"`
	PreviousStep  *openapi_types.UUID `json:"previousStep,omitempty"`
	ScenarioId    *openapi_types.UUID `json:"scenarioId,omitempty"`
}
