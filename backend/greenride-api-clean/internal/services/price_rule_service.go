package services

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"
	"sync"
	"time"

	"greenride/internal/models"
	"greenride/internal/protocol"
	"greenride/internal/utils"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Supported price rule categories
var SupportedRuleCategories = []string{
	protocol.PriceRuleCategoryBasePricing,
	protocol.PriceRuleCategorySurgePricing,
	protocol.PriceRuleCategoryDiscount,
	protocol.PriceRuleCategoryPromotion,
	protocol.PriceRuleCategoryUserPromotion,
	protocol.PriceRuleCategorySpecialOffer,
	protocol.PriceRuleCategoryDistanceFare,
	protocol.PriceRuleCategoryTimeFare,
	protocol.PriceRuleCategoryServiceFee,
}
var (
	OncePriceRuleCategories = []string{
		protocol.PriceRuleCategoryBasePricing,
		protocol.PriceRuleCategorySurgePricing,
		protocol.PriceRuleCategoryDiscount,
		protocol.PriceRuleCategoryPromotion,
		protocol.PriceRuleCategoryUserPromotion,
		protocol.PriceRuleCategoryDistanceFare,
		protocol.PriceRuleCategorySpecialOffer,
		protocol.PriceRuleCategoryTimeFare,
		protocol.PriceRuleCategoryServiceFee,
	}
)

// IsSupportedRuleCategory checks if the rule category is supported
func IsSupportedRuleCategory(category string) bool {
	return slices.Contains(SupportedRuleCategories, category)
}

// GetSupportedRuleCategories returns the list of supported rule categories
func GetSupportedRuleCategories() []string {
	return append([]string{}, SupportedRuleCategories...) // Return a copy to prevent external modification
}

// PriceRuleService manages price rules
type PriceRuleService struct {
	db *gorm.DB
}

var (
	priceRuleServiceInstance *PriceRuleService
	priceRuleServiceOnce     sync.Once
)

// GetPriceRuleService returns the price rule service singleton
func GetPriceRuleService() *PriceRuleService {
	priceRuleServiceOnce.Do(func() {
		SetupPriceRuleService()
	})
	return priceRuleServiceInstance
}

// SetupPriceRuleService initializes the price rule service
func SetupPriceRuleService() {
	priceRuleServiceInstance = &PriceRuleService{
		db: models.GetDB(),
	}
}

// convertProtocolToModelsVehicleFilters 将protocol的VehicleFilter转换为models的VehicleFilter
func convertProtocolToModelsVehicleFilters(protocolFilters []*protocol.VehicleFilter) []*models.VehicleFilter {
	if protocolFilters == nil {
		return []*models.VehicleFilter{}
	}

	modelsFilters := make([]*models.VehicleFilter, len(protocolFilters))
	for i, pf := range protocolFilters {
		modelsFilters[i] = &models.VehicleFilter{
			Category: pf.Category,
			Level:    pf.Level,
		}
	}
	return modelsFilters
}

// 映射函数已删除，现在直接使用VehicleClass和ServiceLevel字段

// CreatePriceRuleRequest represents request for creating a price rule
type CreatePriceRuleRequest struct {
	UserID          string           `json:"user_id"` // Creator user ID
	RuleName        string           `json:"rule_name" binding:"required"`
	Description     string           `json:"description"`
	Category        string           `json:"category" binding:"required"`  // base_pricing, surge_pricing, discount, promotion, special_offer
	RuleType        string           `json:"rule_type" binding:"required"` // percentage, fixed_amount, multiplier, tiered, custom
	DiscountAmount  *float64         `json:"discount_amount"`
	DiscountPercent *float64         `json:"discount_percent"`
	SurgeMultiplier *float64         `json:"surge_multiplier"`
	BaseRate        *float64         `json:"base_rate"`       // Base price
	PerKmRate       *float64         `json:"per_km_rate"`     // Per kilometer price
	PerMinuteRate   *float64         `json:"per_minute_rate"` // Per minute price
	MinimumFare     *float64         `json:"minimum_fare"`
	MaximumFare     *float64         `json:"maximum_fare"`
	Priority        int              `json:"priority"`
	Status          string           `json:"status"` // draft, active, paused, expired
	StartDate       *int64           `json:"start_date"`
	EndDate         *int64           `json:"end_date"`
	VehicleFilters  []*VehicleFilter `json:"vehicle_filters"`
	ServiceAreas    []string         `json:"service_areas"`
	UserCategories  []string         `json:"user_categories"`
	ApplicableRides []string         `json:"applicable_rides"`
	MaxUsagePerUser *int             `json:"max_usage_per_user"`
	MaxUsageTotal   *int             `json:"max_usage_total"`
	PromoCode       string           `json:"promo_code"`
	RequiresCode    int              `json:"requires_code"` // 0:No 1:Yes
	Metadata        map[string]any   `json:"metadata"`
}

type VehicleFilter struct {
	Category string `json:"category"`
	Level    string `json:"level"`
}

// EnvironmentContext represents environmental factors affecting pricing
type EnvironmentContext struct {
	WeatherCondition string
	TrafficLevel     string
	DemandLevel      string
	SupplyLevel      string
	SurgeMultiplier  float64
	EventFactors     map[string]float64
}

// IsRuleApplied unified rule applicability judgment
func (s *PriceRuleService) IsRuleApplied(ctx *PriceContext) (bool, string) {
	rule := ctx.Rule
	req := ctx.Request
	// 1. Time validity check
	now := time.Now().UnixMilli()
	if rule.StartedAt != nil && now < *rule.StartedAt {
		return false, "Rule not effective yet"
	}
	if rule.EndedAt != nil && now > *rule.EndedAt {
		return false, "Rule has expired"
	}

	// 2. Vehicle type restriction check
	filters := rule.GetVehicleFilters()
	if len(filters) > 0 {
		if !rule.IsApplicableToVehicle(req.VehicleCategory, req.VehicleLevel) {
			return false, fmt.Sprintf("Not applicable to vehicle type: %s/%s", req.VehicleCategory, req.VehicleLevel)
		}
	}

	// 3. Service area restriction check
	serviceAreas := rule.ServiceAreas
	if len(serviceAreas) > 0 && !slices.Contains(serviceAreas, req.ServiceArea) {
		return false, fmt.Sprintf("Not applicable to service area: %v", req.ServiceArea)
	}

	// 4. User category restriction check
	userCategories := rule.UserCategories
	if len(userCategories) > 0 && !slices.Contains(userCategories, req.UserCategory) {
		return false, fmt.Sprintf("Not applicable to user category: %v", req.UserCategory)
	}

	// 5. Status check
	if rule.GetStatus() != protocol.StatusActive {
		return false, "Rule not activated"
	}

	// 6. Special checks based on rule category
	category := rule.GetCategory()
	switch category {
	case protocol.PriceRuleCategorySurgePricing:
		// Surge pricing rule needs to check surge conditions
		if ctx.Env.SurgeMultiplier <= 1.0 {
			return false, "No surge pricing needed currently"
		}

	case protocol.PriceRuleCategoryPromotion:
		// Regular promotion rules need to check promo codes
		if rule.RequiresCode != nil && *rule.RequiresCode == 1 {
			if len(req.PromoCodes) == 0 {
				return false, "Promo code required"
			}

			if rule.PromoCode != nil {
				matched := slices.Contains(req.PromoCodes, *rule.PromoCode)
				if !matched {
					return false, "Invalid promo code"
				}
			}
		}

	case protocol.PriceRuleCategoryDiscount:
		// Discount rules need to check usage limit
		if rule.MaxUsagePerUser != nil {
			// TODO: Query user usage count
			// usedCount := s.getUserRuleUsageCount(req.UserID, rule.RuleID)
			// if usedCount >= *rule.MaxUsagePerUser {
			//     return false, fmt.Sprintf("Reached personal usage limit: %d times", *rule.MaxUsagePerUser)
			// }
		}

		if rule.MaxUsageTotal != nil {
			// TODO: Query total usage count
			// totalUsedCount := s.getTotalRuleUsageCount(rule.RuleID)
			// if totalUsedCount >= *rule.MaxUsageTotal {
			//     return false, fmt.Sprintf("Reached total usage limit: %d times", *rule.MaxUsageTotal)
			// }
		}
	case protocol.PriceRuleCategoryUserPromotion:
	}

	return true, ""
}

// CalculateRule unified rule calculation entry point - includes applicability judgment and price calculation
func (s *PriceRuleService) CalculateRule(ctx *PriceContext) (result *protocol.PriceRuleResult) {
	rule := ctx.Rule
	result = &protocol.PriceRuleResult{
		Applied: false,
		Reason:  "Rule is empty",
	}
	// Check if this type has already been applied (by checking if key exists in map)
	if _, exists := ctx.RuleResults[rule.GetCategory()]; exists {
		return
	}
	if rule == nil {
		return
	}
	// Basic information
	result = &protocol.PriceRuleResult{
		RuleID:      rule.RuleID,
		RuleName:    rule.GetRuleName(),
		Category:    rule.GetCategory(),
		DisplayName: rule.GetDisplayName(),
		//Description: rule.GetDescription(),
	}

	// Check if rule is applicable
	applied, reason := s.IsRuleApplied(ctx)
	if !applied {
		result.Applied = false
		result.Reason = reason
		return result
	}
	defer func() {
		if result.Applied {
			log.Printf("Rule applied: %v (%v), calculation result: %v%.2f, description: %v", rule.RuleName, rule.GetCategory(), ctx.Request.Currency, result.Amount, result.Description)
			//Only apply once rules
			if slices.Contains(OncePriceRuleCategories, rule.GetCategory()) {
				ctx.RuleResults[rule.GetCategory()] = result
			}
		} else {
			log.Printf("Rule not applied: %v (%v), reason: %v", rule.RuleName, rule.GetCategory(), result.Reason)
		}
	}()

	// If applicable, perform specific calculation
	category := rule.GetCategory()
	switch category {
	case protocol.PriceRuleCategoryBasePricing:
		result = s.CalculateBasePricing(ctx)
		ctx.Snapshot.SetBaseFare(ctx.Snapshot.GetBaseFare().Add(decimal.NewFromFloat(result.Amount)))
		return
	case protocol.PriceRuleCategorySurgePricing:
		result = s.CalculateSurgePricing(ctx)
		ctx.Snapshot.SetSurgeFare(ctx.Snapshot.GetSurgeFare().Add(decimal.NewFromFloat(result.Amount)))
		return
	case protocol.PriceRuleCategoryDiscount:
		result = s.CalculateDiscount(ctx)
		newAmount := ctx.Snapshot.GetDiscountAmount().Add(decimal.NewFromFloat(result.Amount))
		ctx.Snapshot.SetDiscountAmount(newAmount)
		return
	case protocol.PriceRuleCategoryPromotion:
		result = s.CalculatePromotion(ctx)
		newAmount := ctx.Snapshot.GetPromoDiscount().Add(decimal.NewFromFloat(result.Amount))
		ctx.Snapshot.SetPromoDiscount(newAmount)
		return
	case protocol.PriceRuleCategoryUserPromotion:
		result = s.CalculateUserPromotion(ctx)
		newAmount := ctx.Snapshot.GetUserPromoDiscount().Add(decimal.NewFromFloat(result.Amount))
		ctx.Snapshot.SetUserPromoDiscount(newAmount)
		ctx.UserPromotionIDs = append(ctx.UserPromotionIDs, rule.RuleID)
		return
	case protocol.PriceRuleCategorySpecialOffer:
		result = s.CalculateSpecialOffer(ctx)
		newAmount := ctx.Snapshot.GetPromoDiscount().Add(decimal.NewFromFloat(result.Amount))
		ctx.Snapshot.SetPromoDiscount(newAmount)
		return
	case protocol.PriceRuleCategoryDistanceFare:
		result = s.CalculateDistanceFare(ctx)
		ctx.Snapshot.SetDistanceFare(ctx.Snapshot.GetDistanceFare().Add(decimal.NewFromFloat(result.Amount)))
		return
	case protocol.PriceRuleCategoryTimeFare:
		result = s.CalculateTimeFare(ctx)
		ctx.Snapshot.SetTimeFare(ctx.Snapshot.GetTimeFare().Add(decimal.NewFromFloat(result.Amount)))
		return
	case protocol.PriceRuleCategoryServiceFee:
		result = s.CalculateServiceFee(ctx)
		ctx.Snapshot.SetServiceFee(ctx.Snapshot.GetServiceFee().Add(decimal.NewFromFloat(result.Amount)))
		return
	default:
		result.Applied = false
		result.Reason = fmt.Sprintf("Unsupported rule category: %v (supported types: %v)", category, SupportedRuleCategories)
		return result
	}
}

// CalculateBasePricing calculates base pricing
func (s *PriceRuleService) CalculateBasePricing(ctx *PriceContext) (result *protocol.PriceRuleResult) {
	result = &protocol.PriceRuleResult{
		RuleID:      ctx.Rule.RuleID,
		RuleName:    ctx.Rule.GetRuleName(),
		Category:    ctx.Rule.GetCategory(),
		DisplayName: ctx.Rule.GetDisplayName(),
		//Description: ctx.Rule.GetDescription(),
		Applied: true,
	}
	rule := ctx.Rule
	req := ctx.Request
	amount := 0.0
	description := ""
	switch rule.GetRuleType() {
	case protocol.PriceRuleTypeFixedAmount:
		// Fixed amount pricing - supports segmented pricing mode
		if rule.BaseRate != nil && rule.PerKmRate != nil && rule.MinDistance != nil {
			// Segmented pricing mode: base fare includes basic kilometers
			basePrice := *rule.BaseRate     // Base fare
			includedKm := *rule.MinDistance // Included basic kilometers
			perKmRate := *rule.PerKmRate    // Per kilometer fee for excess distance
			distance := req.EstimatedDistance

			if distance <= includedKm {
				// Distance within base fare coverage
				amount = basePrice
				description = fmt.Sprintf("Base fare: %v%.0f (includes %.1fkm)", req.Currency, basePrice, includedKm)
			} else {
				// Distance exceeds coverage, need to calculate excess charge
				excessDistance := distance - includedKm
				excessAmount := excessDistance * perKmRate
				amount = basePrice + excessAmount
				description = fmt.Sprintf("Base fare: %v%.0f (includes %.1fkm) + Excess charge: %v%.0f/km × %.1fkm = %v%.0f",
					req.Currency, basePrice, includedKm,
					req.Currency, perKmRate, excessDistance,
					req.Currency, amount)
			}
		} else if _v := rule.GetBaseRate(); _v > 0 {
			// Traditional fixed amount pricing
			amount = _v
			description = fmt.Sprintf("Fixed base price: %v%.2f", req.Currency, amount)
		} else {
			result.Applied = false
			result.Reason = "Fixed amount pricing rule missing base price configuration"
			return result
		}

	case protocol.PriceRuleTypePercentage:
		// Percentage pricing (based on standard distance + time rate)
		if rule.DiscountPercent != nil {
			percent := *rule.DiscountPercent
			// Calculate standard price base: base rate + distance rate + time rate
			baseRate := rule.GetBaseRate()
			perKmRate := 0.0
			perMinuteRate := 0.0
			if rule.PerKmRate != nil {
				perKmRate = *rule.PerKmRate
			}
			if rule.PerMinuteRate != nil {
				perMinuteRate = *rule.PerMinuteRate
			}
			standardPrice := baseRate + (req.EstimatedDistance * perKmRate) + (float64(req.EstimatedDuration) * perMinuteRate)
			amount = standardPrice * (percent / 100.0)
			description = fmt.Sprintf("Percentage pricing: %.1f%% x (base%.2f + distance%.2f + time%.2f) = %v%.2f",
				percent, baseRate, req.EstimatedDistance*perKmRate, float64(req.EstimatedDuration)*perMinuteRate, req.Currency, amount)
		} else {
			result.Applied = false
			result.Reason = "Percentage pricing rule missing percentage configuration"
			return result
		}

	case protocol.PriceRuleTypeMultiplier:
		// Multiplier pricing (based on standard distance + time rate)
		if rule.SurgeMultiplier != nil {
			multiplier := *rule.SurgeMultiplier
			// Calculate standard price base: base rate + distance rate + time rate
			baseRate := rule.GetBaseRate()
			perKmRate := 0.0
			perMinuteRate := 0.0
			if rule.PerKmRate != nil {
				perKmRate = *rule.PerKmRate
			}
			if rule.PerMinuteRate != nil {
				perMinuteRate = *rule.PerMinuteRate
			}
			standardPrice := baseRate + (req.EstimatedDistance * perKmRate) + (float64(req.EstimatedDuration) * perMinuteRate)
			amount = standardPrice * multiplier
			description = fmt.Sprintf("Multiplier pricing: %.2fx x (base%.2f + distance%.2f + time%.2f) = %v%.2f",
				multiplier, baseRate, req.EstimatedDistance*perKmRate, float64(req.EstimatedDuration)*perMinuteRate, req.Currency, amount)
		} else {
			result.Applied = false
			result.Reason = "Multiplier pricing rule missing multiplier configuration"
			return result
		}

	case protocol.PriceRuleTypeTiered:
		// Tiered pricing
		if rule.TieredRules != nil && len(rule.TieredRules.Tiers) > 0 {
			// Select appropriate tier based on distance
			for _, tier := range rule.TieredRules.Tiers {
				if tier == nil {
					continue
				}

				// Check distance conditions
				if *tier.MinDuration > 0 && req.EstimatedDistance < *tier.MinDuration {
					continue
				}
				if tier.MaxDuration != nil && req.EstimatedDistance > *tier.MaxDuration {
					continue
				}

				// Apply tier rate
				if tier.Rate > 0 {
					// Base rate + tier rate calculated by distance
					baseRate := rule.GetBaseRate()
					distanceAmount := req.EstimatedDistance * tier.Rate
					amount = baseRate + distanceAmount
					description = fmt.Sprintf("Tiered pricing: Base fare%.2f + Tier rate%v%.2f/km x %.1fkm = %v%.2f",
						baseRate, req.Currency, tier.Rate, req.EstimatedDistance, req.Currency, amount)
				}
				break // Exit after finding the first matching tier
			}

			if amount == 0.0 && description == "" {
				result.Applied = false
				result.Reason = "No matching tier rule found"
				return result
			}
		} else {
			result.Applied = false
			result.Reason = "Tiered pricing rule missing tier configuration"
			return result
		}

	case protocol.PriceRuleTypeCustom:
		// Custom pricing: comprehensive calculation based on distance, time and dynamic factors
		baseRate := rule.GetBaseRate()
		perKmRate := 0.0
		perMinuteRate := 0.0

		if rule.PerKmRate != nil {
			perKmRate = *rule.PerKmRate
		}
		if rule.PerMinuteRate != nil {
			perMinuteRate = *rule.PerMinuteRate
		}

		// Base calculation
		amount = baseRate + (req.EstimatedDistance * perKmRate) + (float64(req.EstimatedDuration) * perMinuteRate)

		// Apply dynamic factors
		dynamicMultiplier := 1.0
		factors := []string{}

		if rule.DemandFactor != nil && *rule.DemandFactor != 1.0 {
			dynamicMultiplier *= *rule.DemandFactor
			factors = append(factors, fmt.Sprintf("demand:%.2fx", *rule.DemandFactor))
		}

		if rule.SupplyFactor != nil && *rule.SupplyFactor != 1.0 {
			dynamicMultiplier *= *rule.SupplyFactor
			factors = append(factors, fmt.Sprintf("supply:%.2fx", *rule.SupplyFactor))
		}

		if rule.WeatherFactor != nil && *rule.WeatherFactor != 1.0 {
			dynamicMultiplier *= *rule.WeatherFactor
			factors = append(factors, fmt.Sprintf("weather:%.2fx", *rule.WeatherFactor))
		}

		amount *= dynamicMultiplier

		if len(factors) > 0 {
			description = fmt.Sprintf("Custom pricing: Base cost%.2f x Dynamic factor%.2fx (%v) = %v%.2f",
				amount/dynamicMultiplier, dynamicMultiplier, strings.Join(factors, ", "), req.Currency, amount)
		} else {
			description = fmt.Sprintf("Custom pricing: Base fare%.2f + Distance fee%.2f + Time fee%.2f = %v%.2f",
				baseRate, req.EstimatedDistance*perKmRate, float64(req.EstimatedDuration)*perMinuteRate, req.Currency, amount)
		}

	default:
		// Default distance + time calculation method
		baseRate := rule.GetBaseRate()
		perKmRate := 0.0
		perMinuteRate := 0.0

		if rule.PerKmRate != nil {
			perKmRate = *rule.PerKmRate
		}
		if rule.PerMinuteRate != nil {
			perMinuteRate = *rule.PerMinuteRate
		}

		amount = baseRate + (req.EstimatedDistance * perKmRate) + (float64(req.EstimatedDuration) * perMinuteRate)
		description = fmt.Sprintf("Base pricing: Base fare%.2f + Distance fee%.2f + Time fee%.2f = %v%.2f",
			baseRate, req.EstimatedDistance*perKmRate, float64(req.EstimatedDuration)*perMinuteRate, req.Currency, amount)
	}

	// Apply minimum/maximum fare limits
	if _v := rule.GetMinimumFare(); amount < _v {
		amount = _v
		description += fmt.Sprintf(" (Apply minimum fare: %v%.2f)", req.Currency, amount)
	}

	if _v := rule.GetMaximumFare(); amount > _v {
		amount = _v
		description += fmt.Sprintf(" (Apply maximum fare: %v%.2f)", req.Currency, amount)
	}
	result.Amount = amount
	result.Description = description
	return result
}

// CalculateSurgePricing calculates surge pricing
func (s *PriceRuleService) CalculateSurgePricing(ctx *PriceContext) (result *protocol.PriceRuleResult) {
	result = &protocol.PriceRuleResult{
		RuleID:      ctx.Rule.RuleID,
		RuleName:    ctx.Rule.GetRuleName(),
		Category:    ctx.Rule.GetCategory(),
		DisplayName: ctx.Rule.GetDisplayName(),
		//Description: ctx.Rule.GetDescription(),
		Applied: true,
	}
	rule := ctx.Rule
	req := ctx.Request

	amount := 0.0
	description := ""
	switch rule.GetRuleType() {
	case protocol.PriceRuleTypeMultiplier:
		if rule.SurgeMultiplier != nil {
			multiplier := *rule.SurgeMultiplier
			// Apply surge pricing based on current calculated base fare
			currentBaseFare := ctx.Snapshot.GetBaseFare()
			if currentBaseFare.Equal(decimal.Zero) {
				// If no base fare calculated yet, use BasePrice from request
				currentBaseFare = decimal.NewFromFloat(req.BasePrice)
			}
			// Calculate only the additional amount
			multiplierDec := decimal.NewFromFloat(multiplier).Sub(decimal.NewFromFloat(1.0))
			amountDec := currentBaseFare.Mul(multiplierDec)
			amount, _ = amountDec.Float64()
			currentBaseFareFloat, _ := currentBaseFare.Float64()
			description = fmt.Sprintf("Surge multiplier: %.2fx, Base price: %v%.2f, Additional fee: %v%.2f",
				multiplier, req.Currency, currentBaseFareFloat, req.Currency, amount)
		} else {
			result.Applied = false
			result.Reason = "Multiplier surge pricing rule missing multiplier configuration"
			return result
		}

	case protocol.PriceRuleTypeFixedAmount:
		if rule.DiscountAmount != nil {
			amount = *rule.DiscountAmount
			description = fmt.Sprintf("Fixed surge: %v%.2f", req.Currency, amount)
		} else {
			result.Applied = false
			result.Reason = "Fixed surge pricing rule missing amount configuration"
			return result
		}

	case protocol.PriceRuleTypePercentage:
		if rule.DiscountPercent != nil {
			percent := *rule.DiscountPercent
			// Apply percentage surge based on current calculated base fare
			currentBaseFare := ctx.Snapshot.GetBaseFare()
			if currentBaseFare.Equal(decimal.Zero) {
				// If no base fare calculated yet, use BasePrice from request
				currentBaseFare = decimal.NewFromFloat(req.BasePrice)
			}
			// Percentage surge
			percentDec := decimal.NewFromFloat(percent).Div(decimal.NewFromFloat(100.0))
			amountDec := currentBaseFare.Mul(percentDec)
			amount, _ = amountDec.Float64()
			currentBaseFareFloat, _ := currentBaseFare.Float64()
			description = fmt.Sprintf("Percentage surge: %.1f%%, Base price: %v%.2f, Additional fee: %v%.2f",
				percent, req.Currency, currentBaseFareFloat, req.Currency, amount)
		} else {
			result.Applied = false
			result.Reason = "Percentage surge pricing rule missing percentage configuration"
			return result
		}

	case protocol.PriceRuleTypeTiered:
		if rule.TieredRules != nil && len(rule.TieredRules.Tiers) > 0 {
			// Select appropriate tier based on distance
			for _, tier := range rule.TieredRules.Tiers {
				if tier == nil {
					continue
				}

				// Check distance conditions
				if *tier.MinDuration > 0 && req.EstimatedDistance < *tier.MinDuration {
					continue
				}
				if tier.MaxDuration != nil && req.EstimatedDistance > *tier.MaxDuration {
					continue
				}

				// Apply tier rate
				if tier.Rate > 0 {
					// Tier rate calculated by distance
					amount = req.EstimatedDistance * tier.Rate
					description = fmt.Sprintf("Tiered rate: %v%.2f/km (distance:%.1fkm), Additional fee: %v%.2f",
						req.Currency, tier.Rate, req.EstimatedDistance, req.Currency, amount)
				}
				break // Exit after finding the first matching tier
			}

			if amount == 0.0 && description == "" {
				result.Applied = false
				result.Reason = "No matching tier rule found"
				return result
			}
		} else {
			result.Applied = false
			result.Reason = "Tiered surge pricing rule missing tier configuration"
			return result
		}

	case protocol.PriceRuleTypeCustom:
		// Custom rule: dynamic calculation based on environmental factors
		if ctx.Env != nil {
			multiplier := 1.0
			factors := []string{}

			// Apply demand factor
			if rule.DemandFactor != nil && *rule.DemandFactor > 1.0 {
				multiplier *= *rule.DemandFactor
				factors = append(factors, fmt.Sprintf("demand:%.2fx", *rule.DemandFactor))
			}

			// Apply supply factor
			if rule.SupplyFactor != nil && *rule.SupplyFactor > 1.0 {
				multiplier *= *rule.SupplyFactor
				factors = append(factors, fmt.Sprintf("supply:%.2fx", *rule.SupplyFactor))
			}

			// Apply weather factor
			if rule.WeatherFactor != nil && *rule.WeatherFactor > 1.0 {
				multiplier *= *rule.WeatherFactor
				factors = append(factors, fmt.Sprintf("weather:%.2fx", *rule.WeatherFactor))
			}

			// Apply event factor
			if rule.EventFactor != nil && *rule.EventFactor > 1.0 {
				multiplier *= *rule.EventFactor
				factors = append(factors, fmt.Sprintf("event:%.2fx", *rule.EventFactor))
			}

			// Apply event factors from environmental context
			if len(ctx.Env.EventFactors) > 0 {
				for event, factor := range ctx.Env.EventFactors {
					if factor > 1.0 {
						multiplier *= factor
						factors = append(factors, fmt.Sprintf("%v:%.2fx", event, factor))
					}
				}
			}

			if multiplier > 1.0 {
				amount = req.BasePrice * (multiplier - 1.0)
				description = fmt.Sprintf("Custom surge: %.2fx (%v), Additional fee: %v%.2f",
					multiplier, strings.Join(factors, ", "), req.Currency, amount)
			} else {
				result.Applied = false
				result.Reason = "Custom rule did not generate surge effect"
				return result
			}
		} else {
			result.Applied = false
			result.Reason = "Custom surge pricing rule missing environmental context"
			return result
		}

	default:
		result.Applied = false
		result.Reason = fmt.Sprintf("Surge pricing does not support rule type: %v", rule.GetRuleType())
		return result
	}
	result.Amount = amount
	result.Description = description
	return result
}

// CalculateDiscount calculates discount amount
func (s *PriceRuleService) CalculateDiscount(ctx *PriceContext) (result *protocol.PriceRuleResult) {
	result = &protocol.PriceRuleResult{
		RuleID:      ctx.Rule.RuleID,
		RuleName:    ctx.Rule.GetRuleName(),
		Category:    ctx.Rule.GetCategory(),
		DisplayName: ctx.Rule.GetDisplayName(),
		//Description: ctx.Rule.GetDescription(),
		Applied: true,
	}
	defer func() {
		if result.Applied && result.Amount > 0 {
			result.Amount = -1 * result.Amount
		}
	}()

	rule := ctx.Rule
	req := ctx.Request
	amount := 0.0
	description := ""

	defer func() {
		if result.Applied {
			log.Printf("Base pricing rule applied: %v, calculation result: %v%.2f, description: %v", rule.RuleName, req.Currency, result.Amount, description)
		}
	}()

	switch rule.GetRuleType() {
	case protocol.PriceRuleTypePercentage:
		if rule.DiscountPercent != nil {
			percent := *rule.DiscountPercent
			// Apply discount based on current calculated price (base price + surge)
			currentTotalFare := ctx.Snapshot.GetBaseFare().Add(ctx.Snapshot.GetSurgeFare())
			if utils.IsDecimalZero(currentTotalFare) {
				// If no price calculated yet, use BasePrice from request
				currentTotalFare = decimal.NewFromFloat(req.BasePrice)
			}
			// Calculate percentage
			percentDecimal := decimal.NewFromFloat(percent).Div(decimal.NewFromFloat(100.0))
			amountDecimal := currentTotalFare.Mul(percentDecimal)
			amount = utils.DecimalToFloat64(amountDecimal)
			currentTotalFareFloat := utils.DecimalToFloat64(currentTotalFare)

			description = fmt.Sprintf("Percentage discount: %.1f%%, Based on amount: %v%.2f, Discount amount: %v%.2f",
				percent, req.Currency, currentTotalFareFloat, req.Currency, amount)
		} else {
			result.Applied = false
			result.Reason = "Percentage discount rule missing percentage configuration"
			return result
		}

	case protocol.PriceRuleTypeFixedAmount:
		if rule.DiscountAmount != nil {
			amount = *rule.DiscountAmount
			description = fmt.Sprintf("Fixed discount: %v%.2f", req.Currency, amount)
		} else {
			result.Applied = false
			result.Reason = "Fixed discount rule missing discount amount configuration"
			return result
		}

	case protocol.PriceRuleTypeMultiplier:
		if rule.SurgeMultiplier != nil {
			multiplier := *rule.SurgeMultiplier
			// Apply discount based on current calculated price
			currentTotalFare := ctx.Snapshot.GetBaseFare().Add(ctx.Snapshot.GetSurgeFare())
			if utils.IsDecimalZero(currentTotalFare) {
				// If no price calculated yet, use BasePrice from request
				currentTotalFare = decimal.NewFromFloat(req.BasePrice)
			}

			// Calculate discount amount
			multiplierDec := decimal.NewFromFloat(1.0 - multiplier)
			amountDec := currentTotalFare.Mul(multiplierDec)
			amount = utils.DecimalToFloat64(amountDec)

			currentTotalFareFloat := utils.DecimalToFloat64(currentTotalFare)
			description = fmt.Sprintf("Discount multiplier: %.2fx, Based on amount: %v%.2f, Discount amount: %v%.2f",
				multiplier, req.Currency, currentTotalFareFloat, req.Currency, amount)
		} else {
			result.Applied = false
			result.Reason = "Multiplier discount rule missing multiplier configuration"
			return result
		}

	case protocol.PriceRuleTypeTiered:
		if rule.TieredRules != nil && len(rule.TieredRules.Tiers) > 0 {
			// Select appropriate tier discount based on distance
			for _, tier := range rule.TieredRules.Tiers {
				if tier == nil {
					continue
				}

				// Check distance conditions
				if *tier.MinDuration > 0 && req.EstimatedDistance < *tier.MinDuration {
					continue
				}
				if tier.MaxDuration != nil && req.EstimatedDistance > *tier.MaxDuration {
					continue
				}

				// Apply tier discount rate
				if tier.Rate > 0 {
					// Tier discount rate should be a value between 0-1
					discountRate := tier.Rate
					if discountRate > 1.0 {
						discountRate = discountRate / 100.0 // Convert percentage
					}
					// Apply discount based on current calculated price
					currentTotalFare := ctx.Snapshot.GetBaseFare().Add(ctx.Snapshot.GetSurgeFare())
					if utils.IsDecimalZero(currentTotalFare) {
						// If no price calculated yet, use BasePrice from request
						currentTotalFare = decimal.NewFromFloat(req.BasePrice)
					}

					discountRateDecimal := decimal.NewFromFloat(discountRate)
					amountDecimal := currentTotalFare.Mul(discountRateDecimal)
					amount = utils.DecimalToFloat64(amountDecimal)
					currentTotalFareFloat := utils.DecimalToFloat64(currentTotalFare)

					description = fmt.Sprintf("Tiered discount: %.1f%% (distance:%.1fkm), Based on amount: %v%.2f, Discount amount: %v%.2f",
						discountRate*100, req.EstimatedDistance, req.Currency, currentTotalFareFloat, req.Currency, amount)
				}
				break // Exit after finding the first matching tier
			}

			if amount == 0.0 && description == "" {
				result.Applied = false
				result.Reason = "No matching tier discount rule found"
				return result
			}
		} else {
			result.Applied = false
			result.Reason = "Tiered discount rule missing tier configuration"
			return result
		}

	case protocol.PriceRuleTypeCustom:
		// Custom discount: compound discount based on multiple factors
		if req.BasePrice > 0 {
			discountRate := 0.0
			factors := []string{}

			// Base discount percentage
			if rule.DiscountPercent != nil {
				discountRate = *rule.DiscountPercent / 100.0
				factors = append(factors, fmt.Sprintf("base:%.1f%%", *rule.DiscountPercent))
			}

			// Distance-related discount
			if rule.PerKmRate != nil && req.EstimatedDistance > 0 {
				kmDiscount := *rule.PerKmRate * req.EstimatedDistance
				amount += kmDiscount
				factors = append(factors, fmt.Sprintf("distance:%v%.2f", req.Currency, kmDiscount))
			}

			// Base percentage discount
			if discountRate > 0 {
				baseDiscountAmount := req.BasePrice * discountRate
				amount += baseDiscountAmount
			}

			if len(factors) > 0 {
				description = fmt.Sprintf("Custom discount (%v), Total discount: %v%.2f",
					strings.Join(factors, " + "), req.Currency, amount)
			} else {
				result.Applied = false
				result.Reason = "Custom discount rule has no valid discount parameters configured"
				return result
			}
		} else {
			result.Applied = false
			result.Reason = "Custom discount rule missing base price"
			return result
		}

	default:
		result.Applied = false
		result.Reason = fmt.Sprintf("Discount does not support rule type: %v", rule.GetRuleType())
		return result
	}

	// Apply maximum discount limit
	if rule.MaxDiscount != nil && amount > *rule.MaxDiscount {
		amount = *rule.MaxDiscount
		description += fmt.Sprintf(" (Apply maximum discount limit: %v%.2f)", req.Currency, amount)
	}

	// Ensure discount does not exceed current total price
	currentTotalForLimit := ctx.Snapshot.GetBaseFare().Add(ctx.Snapshot.GetSurgeFare())
	if utils.IsDecimalZero(currentTotalForLimit) {
		currentTotalForLimit = decimal.NewFromFloat(req.BasePrice)
	}
	currentTotalForLimitFloat := utils.DecimalToFloat64(currentTotalForLimit)
	if currentTotalForLimitFloat > 0 && amount > currentTotalForLimitFloat {
		amount = currentTotalForLimitFloat
		description += " (Discount cannot exceed current total price)"
	}

	result.Amount = amount
	result.Description = description
	return result
}

// CalculatePromotion calculates promotion discount
func (s *PriceRuleService) CalculatePromotion(ctx *PriceContext) (result *protocol.PriceRuleResult) {
	result = &protocol.PriceRuleResult{
		RuleID:      ctx.Rule.RuleID,
		RuleName:    ctx.Rule.GetRuleName(),
		Category:    ctx.Rule.GetCategory(),
		DisplayName: ctx.Rule.GetDisplayName(),
		//Description: ctx.Rule.GetDescription(),
		Applied: true,
	}
	defer func() {
		if result.Applied && result.Amount > 0 {
			result.Amount = -1 * result.Amount
		}
	}()
	rule := ctx.Rule
	req := ctx.Request
	amount := 0.0
	description := ""

	switch rule.GetRuleType() {
	case protocol.PriceRuleTypePercentage:
		if rule.DiscountPercent != nil {
			percent := *rule.DiscountPercent
			// Apply promotional discount based on current calculated price
			currentTotalFare := ctx.Snapshot.GetBaseFare().Add(ctx.Snapshot.GetSurgeFare())
			if utils.IsDecimalZero(currentTotalFare) {
				currentTotalFare = decimal.NewFromFloat(req.BasePrice)
			}

			// Calculate percentage
			percentDecimal := decimal.NewFromFloat(percent).Div(decimal.NewFromFloat(100.0))
			amountDecimal := currentTotalFare.Mul(percentDecimal)
			amount = utils.DecimalToFloat64(amountDecimal)
			currentTotalFareFloat := utils.DecimalToFloat64(currentTotalFare)

			description = fmt.Sprintf("Promotional discount: %.1f%%, Based on amount: %v%.2f, Discount amount: %v%.2f",
				percent, req.Currency, currentTotalFareFloat, req.Currency, amount)
		} else {
			result.Applied = false
			result.Reason = "Percentage promotion rule missing percentage configuration"
			return result
		}

	case protocol.PriceRuleTypeFixedAmount:
		if rule.DiscountAmount != nil {
			amount = *rule.DiscountAmount
			description = fmt.Sprintf("Fixed promotional discount: %v%.2f", req.Currency, amount)
		} else {
			result.Applied = false
			result.Reason = "Fixed promotion rule missing discount amount configuration"
			return result
		}

	case protocol.PriceRuleTypeMultiplier:
		if rule.SurgeMultiplier != nil {
			multiplier := *rule.SurgeMultiplier
			if multiplier >= 1.0 {
				result.Applied = false
				result.Reason = "Promotion multiplier must be less than 1.0"
				return result
			}
			// Apply promotional discount based on current calculated price
			currentTotalFare := ctx.Snapshot.GetBaseFare().Add(ctx.Snapshot.GetSurgeFare())
			if utils.IsDecimalZero(currentTotalFare) {
				currentTotalFare = decimal.NewFromFloat(req.BasePrice)
			}

			// Calculate discount amount
			multiplierDec := decimal.NewFromFloat(1.0 - multiplier)
			amountDec := currentTotalFare.Mul(multiplierDec)
			amount = utils.DecimalToFloat64(amountDec)
			currentTotalFareFloat := utils.DecimalToFloat64(currentTotalFare)

			description = fmt.Sprintf("Promotion multiplier discount: %.2fx, Based on amount: %v%.2f, Discount amount: %v%.2f",
				multiplier, req.Currency, currentTotalFareFloat, req.Currency, amount)
		} else {
			result.Applied = false
			result.Reason = "Multiplier promotion rule missing multiplier configuration"
			return result
		}

	case protocol.PriceRuleTypeTiered:
		// Promotion tiered rules, may be based on order amount, distance, etc.
		if rule.TieredRules != nil && len(rule.TieredRules.Tiers) > 0 {
			// Apply tiered promotion based on current calculated price
			currentTotalFare := ctx.Snapshot.GetBaseFare().Add(ctx.Snapshot.GetSurgeFare())
			if utils.IsDecimalZero(currentTotalFare) {
				currentTotalFare = decimal.NewFromFloat(req.BasePrice)
			}
			for _, tier := range rule.TieredRules.Tiers {
				if tier == nil {
					continue
				}

				// Distance-based tiered promotion
				if *tier.MinDuration > 0 && req.EstimatedDistance < *tier.MinDuration {
					continue
				}
				if tier.MaxDuration != nil && req.EstimatedDistance > *tier.MaxDuration {
					continue
				}

				// Apply tiered promotion rate
				if tier.Rate > 0 {
					promotionRate := tier.Rate
					if promotionRate > 1.0 {
						promotionRate = promotionRate / 100.0 // Convert percentage
					}

					promotionRateDecimal := decimal.NewFromFloat(promotionRate)
					amountDecimal := currentTotalFare.Mul(promotionRateDecimal)
					amount = utils.DecimalToFloat64(amountDecimal)
					currentTotalFareFloat := utils.DecimalToFloat64(currentTotalFare)

					description = fmt.Sprintf("Tiered promotion: %.1f%% (distance:%.1fkm), Based on amount: %v%.2f, Discount amount: %v%.2f",
						promotionRate*100, req.EstimatedDistance, req.Currency, currentTotalFareFloat, req.Currency, amount)
				}
				break
			}

			if amount == 0.0 && description == "" {
				result.Applied = false
				result.Reason = "No matching tiered promotion rule found"
				return result
			}
		} else {
			result.Applied = false
			result.Reason = "Tiered promotion rule missing tier configuration"
			return result
		}

	case protocol.PriceRuleTypeCustom:
		// Custom promotion: combination of multiple discount methods
		// Apply custom promotion based on current calculated price
		currentTotalFare := ctx.Snapshot.GetBaseFare().Add(ctx.Snapshot.GetSurgeFare())
		if utils.IsDecimalZero(currentTotalFare) {
			currentTotalFare = decimal.NewFromFloat(req.BasePrice)
		}
		currentTotalFareFloat := utils.DecimalToFloat64(currentTotalFare)
		if currentTotalFareFloat > 0 {
			basePromotion := 0.0
			factors := []string{}

			// Base promotion percentage
			if rule.DiscountPercent != nil {
				percentDecimal := decimal.NewFromFloat(*rule.DiscountPercent).Div(decimal.NewFromFloat(100.0))
				basePromotionDecimal := currentTotalFare.Mul(percentDecimal)
				basePromotion = utils.DecimalToFloat64(basePromotionDecimal)
				factors = append(factors, fmt.Sprintf("percentage:%.1f%%", *rule.DiscountPercent))
			}

			// Fixed promotion amount
			if rule.DiscountAmount != nil {
				basePromotion += *rule.DiscountAmount
				factors = append(factors, fmt.Sprintf("fixed:%v%.2f", req.Currency, *rule.DiscountAmount))
			}

			amount = basePromotion

			if len(factors) > 0 {
				description = fmt.Sprintf("Combined promotion (%v), Total discount: %v%.2f",
					strings.Join(factors, " + "), req.Currency, amount)
			} else {
				result.Applied = false
				result.Reason = "Custom promotion rule has no valid discount parameters configured"
				return result
			}
		} else {
			result.Applied = false
			result.Reason = "Custom promotion rule missing base price"
			return result
		}

	default:
		result.Applied = false
		result.Reason = fmt.Sprintf("Unsupported promotion rule type: %v", rule.GetRuleType())
		return result
	}

	// Apply maximum discount limit
	if rule.MaxDiscount != nil && amount > *rule.MaxDiscount {
		amount = *rule.MaxDiscount
		description += fmt.Sprintf(" (Applied max discount limit: %v%.2f)", req.Currency, amount)
	}

	// Ensure discount does not exceed base price
	if req.BasePrice > 0 && amount > req.BasePrice {
		amount = req.BasePrice
		description += " (Discount cannot exceed base price)"
	}

	result.Amount = amount
	result.Description = description
	return result
}

// CalculateUserPromotion User coupon calculation
func (s *PriceRuleService) CalculateUserPromotion(ctx *PriceContext) (result *protocol.PriceRuleResult) {
	result = &protocol.PriceRuleResult{
		Applied:     true,
		RuleID:      ctx.Rule.RuleID,
		RuleName:    ctx.Rule.GetRuleName(),
		Category:    ctx.Rule.GetCategory(),
		Amount:      0,
		DisplayName: ctx.Rule.GetDisplayName(),
		//Description: "User coupon calculation",
	}

	defer func() {
		if result.Applied && result.Amount > 0 {
			result.Amount = -1 * result.Amount
		}
	}()

	rule := ctx.Rule
	req := ctx.Request
	amount := 0.0
	description := ""

	// Get current calculated original price (for percentage calculation)
	currentOriginalFare := ctx.Snapshot.GetBaseFare().Add(ctx.Snapshot.GetSurgeFare()).Add(ctx.Snapshot.GetDistanceFare()).Add(ctx.Snapshot.GetTimeFare()).Add(ctx.Snapshot.GetServiceFee())
	if utils.IsDecimalZero(currentOriginalFare) {
		currentOriginalFare = decimal.NewFromFloat(req.BasePrice)
	}

	// Check minimum order amount requirement
	if currentOriginalFare.LessThan(decimal.NewFromFloat(rule.GetMinOrderAmount())) {
		result.Applied = false
		currentOriginalFareFloat := utils.DecimalToFloat64(currentOriginalFare)
		result.Reason = fmt.Sprintf("Order amount %v%.2f does not meet minimum amount requirement %v%.2f", req.Currency, currentOriginalFareFloat, req.Currency, *rule.MinOrderAmount)
		return result
	}

	switch rule.GetRuleType() {
	case protocol.PriceRuleTypePercentage:
		if _v := rule.GetDiscountPercent(); _v > 0 {
			percentDecimal := decimal.NewFromFloat(_v).Div(decimal.NewFromFloat(100.0))
			amountDecimal := currentOriginalFare.Mul(percentDecimal)
			amount = utils.DecimalToFloat64(amountDecimal)
			currentOriginalFareFloat := utils.DecimalToFloat64(currentOriginalFare)
			description = fmt.Sprintf("User coupon: %.1f%% discount, based on amount: %v%.2f", _v, req.Currency, currentOriginalFareFloat)
		} else {
			result.Applied = false
			result.Reason = "Percentage user coupon missing discount percentage configuration"
			return result
		}

	case protocol.PriceRuleTypeFixedAmount:
		if _v := rule.GetDiscountAmount(); _v > 0 {
			amount = _v
			description = fmt.Sprintf("User coupon: fixed discount %v%.2f", req.Currency, _v)
		} else {
			result.Applied = false
			result.Reason = "Fixed amount user coupon missing discount amount configuration"
			return result
		}

	default:
		result.Applied = false
		result.Reason = fmt.Sprintf("User coupon unsupported rule type: %v (only supports percentage and fixed_amount)", rule.GetRuleType())
		return result
	}

	// Apply maximum discount limit
	if _v := rule.GetMaxDiscount(); amount > _v {
		amount = _v
		description += fmt.Sprintf(" (Applied max discount limit: %v%.2f)", req.Currency, amount)
	}

	// Ensure discount does not exceed current original price
	currentOriginalFareFloat := utils.DecimalToFloat64(currentOriginalFare)
	if currentOriginalFareFloat > 0 && amount > currentOriginalFareFloat {
		amount = currentOriginalFareFloat
		description += " (Discount cannot exceed current total price)"
	}

	// Ensure discount is positive
	if amount <= 0 {
		result.Applied = false
		result.Reason = "Calculated discount amount is invalid"
		return result
	}

	result.Amount = amount
	result.Description = description
	log.Printf("Applied user coupon: %s, discount amount: %v%.2f", description, req.Currency, amount)
	return result
}

// CalculateSpecialOffer Special offer calculation
func (s *PriceRuleService) CalculateSpecialOffer(ctx *PriceContext) (result *protocol.PriceRuleResult) {
	result = &protocol.PriceRuleResult{
		RuleID:      ctx.Rule.RuleID,
		RuleName:    ctx.Rule.GetRuleName(),
		Category:    ctx.Rule.GetCategory(),
		DisplayName: ctx.Rule.GetDisplayName(),
		//Description: ctx.Rule.GetDescription(),
		Applied: true,
	}
	defer func() {
		if result.Applied && result.Amount > 0 {
			result.Amount = -1 * result.Amount
		}
	}()

	rule := ctx.Rule
	req := ctx.Request
	amount := 0.0
	description := ""
	switch rule.GetRuleType() {
	case protocol.PriceRuleTypePercentage:
		if rule.DiscountPercent != nil {
			percent := *rule.DiscountPercent
			// 基于当前已计算的价格进行特殊优惠
			currentTotalFare := ctx.Snapshot.GetBaseFare().Add(ctx.Snapshot.GetSurgeFare())
			if utils.IsDecimalZero(currentTotalFare) {
				currentTotalFare = decimal.NewFromFloat(req.BasePrice)
			}
			percentDecimal := decimal.NewFromFloat(percent).Div(decimal.NewFromFloat(100.0))
			amountDecimal := currentTotalFare.Mul(percentDecimal)
			amount = utils.DecimalToFloat64(amountDecimal)
			currentTotalFareFloat := utils.DecimalToFloat64(currentTotalFare)
			description = fmt.Sprintf("Special offer: %.1f%%, based on amount: %v%.2f, discount amount: %v%.2f",
				percent, req.Currency, currentTotalFareFloat, req.Currency, amount)
		} else {
			result.Applied = false
			result.Reason = "Percentage special offer rule missing percentage configuration"
			return result
		}

	case protocol.PriceRuleTypeFixedAmount:
		if rule.DiscountAmount != nil {
			amount = *rule.DiscountAmount
			description = fmt.Sprintf("Fixed special offer: %v%.2f", req.Currency, amount)
		} else {
			result.Applied = false
			result.Reason = "Fixed special offer rule missing discount amount configuration"
			return result
		}

	case protocol.PriceRuleTypeMultiplier:
		if rule.SurgeMultiplier != nil {
			multiplier := *rule.SurgeMultiplier
			if multiplier >= 1.0 {
				result.Applied = false
				result.Reason = "Special offer multiplier must be less than 1.0"
				return result
			}
			// 基于当前已计算的价格进行特殊优惠
			currentTotalFare := ctx.Snapshot.GetBaseFare().Add(ctx.Snapshot.GetSurgeFare())
			if utils.IsDecimalZero(currentTotalFare) {
				currentTotalFare = decimal.NewFromFloat(req.BasePrice)
			}
			multiplierDec := decimal.NewFromFloat(1.0 - multiplier)
			amountDec := currentTotalFare.Mul(multiplierDec)
			amount = utils.DecimalToFloat64(amountDec)
			currentTotalFareFloat := utils.DecimalToFloat64(currentTotalFare)
			description = fmt.Sprintf("Special offer multiplier: %.2fx, based on amount: %v%.2f, discount amount: %v%.2f",
				multiplier, req.Currency, currentTotalFareFloat, req.Currency, amount)
		} else {
			result.Applied = false
			result.Reason = "Multiplier special offer rule missing multiplier configuration"
			return result
		}

	case protocol.PriceRuleTypeTiered:
		// Special offer tiered rules
		if rule.TieredRules != nil && len(rule.TieredRules.Tiers) > 0 {
			// Based on current calculated price for tiered special offer
			currentTotalFare := ctx.Snapshot.GetBaseFare().Add(ctx.Snapshot.GetSurgeFare())
			if utils.IsDecimalZero(currentTotalFare) {
				currentTotalFare = decimal.NewFromFloat(req.BasePrice)
			}
			for _, tier := range rule.TieredRules.Tiers {
				if tier == nil {
					continue
				}

				// Check distance condition
				if *tier.MinDuration > 0 && req.EstimatedDistance < *tier.MinDuration {
					continue
				}
				if tier.MaxDuration != nil && req.EstimatedDistance > *tier.MaxDuration {
					continue
				}

				// Apply tiered special offer rate
				if tier.Rate > 0 {
					offerRate := tier.Rate
					if offerRate > 1.0 {
						offerRate = offerRate / 100.0 // Convert percentage
					}
					offerRateDecimal := decimal.NewFromFloat(offerRate)
					amountDecimal := currentTotalFare.Mul(offerRateDecimal)
					amount = utils.DecimalToFloat64(amountDecimal)
					currentTotalFareFloat := utils.DecimalToFloat64(currentTotalFare)
					description = fmt.Sprintf("Tiered special offer: %.1f%% (distance:%.1fkm), based on amount: %v%.2f, discount amount: %v%.2f",
						offerRate*100, req.EstimatedDistance, req.Currency, currentTotalFareFloat, req.Currency, amount)
				}
				break
			}

			if amount == 0.0 && description == "" {
				result.Applied = false
				result.Reason = "No matching tiered special offer rule found"
				return result
			}
		} else {
			result.Applied = false
			result.Reason = "Tiered special offer rule missing tier configuration"
			return result
		}

	case protocol.PriceRuleTypeCustom:
		// Custom special offer: may contain complex calculation logic
		// Based on current calculated price for custom special offer
		currentTotalFare := ctx.Snapshot.GetBaseFare().Add(ctx.Snapshot.GetSurgeFare())
		if utils.IsDecimalZero(currentTotalFare) {
			currentTotalFare = decimal.NewFromFloat(req.BasePrice)
		}
		currentTotalFareFloat := utils.DecimalToFloat64(currentTotalFare)
		if currentTotalFareFloat > 0 {
			baseOffer := 0.0
			factors := []string{}

			// Base special offer
			if rule.DiscountPercent != nil {
				percentDecimal := decimal.NewFromFloat(*rule.DiscountPercent).Div(decimal.NewFromFloat(100.0))
				baseOfferDecimal := currentTotalFare.Mul(percentDecimal)
				baseOffer = utils.DecimalToFloat64(baseOfferDecimal)
				factors = append(factors, fmt.Sprintf("Percentage:%.1f%%", *rule.DiscountPercent))
			}

			// Fixed special offer
			if rule.DiscountAmount != nil {
				baseOffer += *rule.DiscountAmount
				factors = append(factors, fmt.Sprintf("Fixed:%v%.2f", req.Currency, *rule.DiscountAmount))
			}

			// Distance-related special offer
			if rule.PerKmRate != nil && req.EstimatedDistance > 0 {
				distanceOffer := *rule.PerKmRate * req.EstimatedDistance
				baseOffer += distanceOffer
				factors = append(factors, fmt.Sprintf("Distance:%v%.2f", req.Currency, distanceOffer))
			}

			amount = baseOffer

			if len(factors) > 0 {
				description = fmt.Sprintf("Custom special offer (%v), total discount: %v%.2f",
					strings.Join(factors, " + "), req.Currency, amount)
			} else {
				result.Applied = false
				result.Reason = "Custom special offer rule has no valid discount parameters configured"
				return result
			}
		} else {
			result.Applied = false
			result.Reason = "Custom special offer rule missing base price"
			return result
		}

	default:
		result.Applied = false
		result.Reason = fmt.Sprintf("Special offer unsupported rule type: %v", rule.GetRuleType())
		return result
	}

	// Apply maximum discount limit
	if rule.MaxDiscount != nil && amount > *rule.MaxDiscount {
		amount = *rule.MaxDiscount
		description += fmt.Sprintf(" (Applied max discount limit: %v%.2f)", req.Currency, amount)
	}

	// Ensure discount does not exceed base price
	if req.BasePrice > 0 && amount > req.BasePrice {
		amount = req.BasePrice
		description += " (Discount cannot exceed base price)"
	}

	result.Amount = amount
	result.Description = description
	return result
}

// CalculateDistanceFare Distance fare calculation
func (s *PriceRuleService) CalculateDistanceFare(ctx *PriceContext) *protocol.PriceRuleResult {
	rule := ctx.Rule
	req := ctx.Request
	result := &protocol.PriceRuleResult{
		Applied:     true,
		RuleID:      rule.RuleID,
		RuleName:    rule.GetRuleName(),
		Category:    rule.GetCategory(),
		Amount:      0,
		DisplayName: rule.GetDisplayName(),
		//Description: "Distance fare calculation",
	}

	distance := req.EstimatedDistance
	if distance <= 0 {
		result.Applied = false
		result.Reason = "Distance is 0, cannot calculate distance fare"
		return result
	}

	amount := 0.0
	description := "Distance fare"

	switch rule.GetRuleType() {
	case protocol.PriceRuleTypeFixedAmount:
		// Fixed amount: fixed cost per kilometer - supports base included kilometers
		if rule.PerKmRate != nil {
			if rule.MinDistance != nil && *rule.MinDistance > 0 {
				// Has base included kilometers, only charge for excess portion
				includedKm := *rule.MinDistance
				if distance > includedKm {
					excessDistance := distance - includedKm
					amount = *rule.PerKmRate * excessDistance
					description = fmt.Sprintf("Distance fare: %v%.2f/km × %.1fkm (exceeds %.1fkm included range)",
						req.Currency, *rule.PerKmRate, excessDistance, includedKm)
				} else {
					amount = 0.0
					description = fmt.Sprintf("Distance fare: %.1fkm within included range (%.1fkm), no additional distance fare", distance, includedKm)
				}
			} else {
				// Traditional mode: full distance billing
				amount = *rule.PerKmRate * distance
				description = fmt.Sprintf("Distance fare: %v%.2f/km × %.1fkm", req.Currency, *rule.PerKmRate, distance)
			}
		}

	case protocol.PriceRuleTypePercentage:
		// 百分比：基于基础价格的百分比
		if rule.DiscountPercent != nil {
			basePrice := ctx.Snapshot.GetBaseFare()
			if utils.IsDecimalZero(basePrice) {
				basePrice = decimal.NewFromFloat(req.BasePrice)
			}
			basePriceFloat := utils.DecimalToFloat64(basePrice)
			if basePriceFloat > 0 {
				amount = basePriceFloat * (*rule.DiscountPercent / 100.0)
				description = fmt.Sprintf("Distance fare: %.1f%% × %v%.2f", *rule.DiscountPercent, req.Currency, basePriceFloat)
			}
		}

	case protocol.PriceRuleTypeTiered:
		// Tiered billing: different rates for different distance segments
		if rule.PerKmRate != nil {
			// Basic implementation: simple tiered billing
			if distance <= 5 {
				amount = *rule.PerKmRate * distance
			} else if distance <= 15 {
				amount = *rule.PerKmRate*5 + (*rule.PerKmRate*0.8)*(distance-5)
			} else {
				amount = *rule.PerKmRate*5 + (*rule.PerKmRate*0.8)*10 + (*rule.PerKmRate*0.6)*(distance-15)
			}
			description = fmt.Sprintf("Distance fare (tiered): %.1fkm", distance)
		}

	case protocol.PriceRuleTypeCustom:
		// Custom distance fare: may contain complex calculation logic
		if rule.PerKmRate != nil && rule.DiscountPercent != nil {
			// Combined calculation: base rate + distance adjustment
			baseRate := *rule.PerKmRate
			adjustment := 1.0 + (*rule.DiscountPercent / 100.0)
			amount = baseRate * distance * adjustment
			description = fmt.Sprintf("Distance fare (custom): %v%.2f/km × %.1fkm × %.2f", req.Currency, baseRate, distance, adjustment)
		} else if rule.PerKmRate != nil {
			amount = *rule.PerKmRate * distance
			description = fmt.Sprintf("Distance fare (custom): %v%.2f/km × %.1fkm", req.Currency, *rule.PerKmRate, distance)
		}

	default:
		result.Applied = false
		result.Reason = fmt.Sprintf("Distance fare unsupported rule type: %v", rule.GetRuleType())
		return result
	}

	result.Amount = amount
	result.Description = description
	return result
}

// CalculateTimeFare Time fare calculation
func (s *PriceRuleService) CalculateTimeFare(ctx *PriceContext) *protocol.PriceRuleResult {
	rule := ctx.Rule
	req := ctx.Request
	result := &protocol.PriceRuleResult{
		Applied:  true,
		RuleID:   rule.RuleID,
		RuleName: rule.GetRuleName(),
		Category: rule.GetCategory(),
		Amount:   0,
		//Description: "Time fare calculation",
	}

	duration := float64(req.EstimatedDuration) // minutes
	if duration <= 0 {
		result.Applied = false
		result.Reason = "Duration is 0, cannot calculate time fare"
		return result
	}

	amount := 0.0
	description := "Time fare"

	switch rule.GetRuleType() {
	case protocol.PriceRuleTypeFixedAmount:
		// Fixed amount: fixed cost per minute
		if rule.PerMinuteRate != nil {
			amount = *rule.PerMinuteRate * duration
			description = fmt.Sprintf("Time fare: %v%.2f/min × %.0fmin", req.Currency, *rule.PerMinuteRate, duration)
		}

	case protocol.PriceRuleTypePercentage:
		// Percentage: based on base price percentage
		if rule.DiscountPercent != nil {
			basePrice := ctx.Snapshot.GetBaseFare()
			if utils.IsDecimalZero(basePrice) {
				basePrice = decimal.NewFromFloat(req.BasePrice)
			}
			basePriceFloat := utils.DecimalToFloat64(basePrice)
			if basePriceFloat > 0 {
				amount = basePriceFloat * (*rule.DiscountPercent / 100.0)
				description = fmt.Sprintf("Time fare: %.1f%% × %v%.2f", *rule.DiscountPercent, req.Currency, basePriceFloat)
			}
		}

	case protocol.PriceRuleTypeTiered:
		// Tiered billing: using tiered rule configuration from database
		if rule.TieredRules != nil && len(rule.TieredRules.Tiers) > 0 {
			durationSeconds := duration * 60 // Convert to seconds

			for _, tier := range rule.TieredRules.Tiers {
				// Check if current duration is within this tier range
				if durationSeconds >= *tier.MinDuration {
					if tier.MaxDuration == nil || durationSeconds <= *tier.MaxDuration {
						// Within this tier range, calculate fee for this segment
						if rule.TieredRules.Unit == "minute" {
							// Billing by minute
							segmentDuration := duration
							if tier.MaxDuration != nil {
								segmentDuration = (*tier.MaxDuration - *tier.MinDuration) / 60
							}
							amount += tier.Rate * segmentDuration
						}
					} else if tier.MaxDuration != nil {
						// Exceeds current tier limit, only calculate fee up to the limit
						segmentDuration := (*tier.MaxDuration - *tier.MinDuration) / 60
						amount += tier.Rate * segmentDuration
					}
				}
			}
			description = fmt.Sprintf("Time fare (tiered): %.0fmin, unit: %s", duration, rule.TieredRules.Unit)
		} else if rule.PerMinuteRate != nil {
			// If no tier configuration, use simple per-minute billing
			amount = *rule.PerMinuteRate * duration
			description = fmt.Sprintf("Time fare (simple): %v%.2f/min × %.0fmin", req.Currency, *rule.PerMinuteRate, duration)
		}

	case protocol.PriceRuleTypeCustom:
		// Custom time fare: may contain complex calculation logic
		if rule.PerMinuteRate != nil && rule.DiscountPercent != nil {
			// Combined calculation: base rate + duration adjustment
			baseRate := *rule.PerMinuteRate
			adjustment := 1.0 + (*rule.DiscountPercent / 100.0)
			amount = baseRate * duration * adjustment
			description = fmt.Sprintf("Time fare (custom): %v%.2f/min × %.0fmin × %.2f", req.Currency, baseRate, duration, adjustment)
		} else if rule.PerMinuteRate != nil {
			amount = *rule.PerMinuteRate * duration
			description = fmt.Sprintf("Time fare (custom): %v%.2f/min × %.0fmin", req.Currency, *rule.PerMinuteRate, duration)
		}

	default:
		result.Applied = false
		result.Reason = fmt.Sprintf("Time fare unsupported rule type: %v", rule.GetRuleType())
		return result
	}

	result.Amount = amount
	result.Description = description
	return result
}

// CalculateServiceFee Service fee calculation
func (s *PriceRuleService) CalculateServiceFee(ctx *PriceContext) *protocol.PriceRuleResult {
	rule := ctx.Rule
	req := ctx.Request
	result := &protocol.PriceRuleResult{
		Applied:  true,
		RuleID:   rule.RuleID,
		RuleName: rule.GetRuleName(),
		Category: rule.GetCategory(),
		Amount:   0,
		//Description: "Service fee calculation",
	}

	amount := 0.0
	description := "Service fee"

	switch rule.GetRuleType() {
	case protocol.PriceRuleTypeFixedAmount:
		// Fixed amount: fixed service fee
		if rule.DiscountAmount != nil {
			amount = *rule.DiscountAmount
			description = fmt.Sprintf("Service fee (fixed): %v%.2f", req.Currency, amount)
		}

	case protocol.PriceRuleTypePercentage:
		// 百分比：基于基础价格的百分比
		if rule.DiscountPercent != nil {
			basePrice := ctx.Snapshot.GetBaseFare().Add(ctx.Snapshot.GetDistanceFare()).Add(ctx.Snapshot.GetTimeFare())
			if utils.IsDecimalZero(basePrice) {
				basePrice = decimal.NewFromFloat(req.BasePrice)
			}
			basePriceFloat := utils.DecimalToFloat64(basePrice)
			if basePriceFloat > 0 {
				amount = basePriceFloat * (*rule.DiscountPercent / 100.0)
				description = fmt.Sprintf("Service fee: %.1f%% × %v%.2f", *rule.DiscountPercent, req.Currency, basePriceFloat)
			}
		}

	case protocol.PriceRuleTypeMultiplier:
		// Multiplier: based on other fees multiplier
		if rule.SurgeMultiplier != nil {
			basePrice := ctx.Snapshot.GetBaseFare().Add(ctx.Snapshot.GetDistanceFare()).Add(ctx.Snapshot.GetTimeFare())
			if utils.IsDecimalZero(basePrice) {
				basePrice = decimal.NewFromFloat(req.BasePrice)
			}
			basePriceFloat := utils.DecimalToFloat64(basePrice)
			if basePriceFloat > 0 {
				amount = basePriceFloat * *rule.SurgeMultiplier
				description = fmt.Sprintf("Service fee (multiplier): %v%.2f × %.2f", req.Currency, basePriceFloat, *rule.SurgeMultiplier)
			}
		}

	case protocol.PriceRuleTypeTiered:
		// Tiered billing: different service fees for different price segments
		if rule.DiscountAmount != nil && rule.DiscountPercent != nil {
			basePrice := ctx.Snapshot.GetBaseFare().Add(ctx.Snapshot.GetDistanceFare()).Add(ctx.Snapshot.GetTimeFare())
			if utils.IsDecimalZero(basePrice) {
				basePrice = decimal.NewFromFloat(req.BasePrice)
			}
			basePriceFloat := utils.DecimalToFloat64(basePrice)

			if basePriceFloat <= 100 {
				amount = *rule.DiscountAmount // Low price order fixed service fee
			} else {
				amount = basePriceFloat * (*rule.DiscountPercent / 100.0) // High price order by percentage
			}
			description = fmt.Sprintf("Service fee (tiered): based on %v%.2f", req.Currency, basePriceFloat)
		}

	case protocol.PriceRuleTypeCustom:
		// Custom service fee: may contain complex calculation logic
		basePrice := ctx.Snapshot.GetBaseFare().Add(ctx.Snapshot.GetDistanceFare()).Add(ctx.Snapshot.GetTimeFare())
		if utils.IsDecimalZero(basePrice) {
			basePrice = decimal.NewFromFloat(req.BasePrice)
		}
		basePriceFloat := utils.DecimalToFloat64(basePrice)

		if rule.DiscountAmount != nil {
			amount += *rule.DiscountAmount // Fixed part
		}

		if rule.DiscountPercent != nil && basePriceFloat > 0 {
			amount += basePriceFloat * (*rule.DiscountPercent / 100.0) // Percentage part
		}

		description = fmt.Sprintf("Service fee (custom): based on %v%.2f", req.Currency, basePriceFloat)

	default:
		result.Applied = false
		result.Reason = fmt.Sprintf("Service fee unsupported rule type: %v", rule.GetRuleType())
		return result
	}

	result.Amount = amount
	result.Description = description
	return result
}

// CreatePriceRule 创建价格规则
func (s *PriceRuleService) CreatePriceRule(req *CreatePriceRuleRequest) (*models.PriceRule, protocol.ErrorCode) {
	// 1. 业务验证
	// 检查规则名称重复
	var existingRule models.PriceRule
	if err := models.GetDB().Where("rule_name = ? AND status != ?", req.RuleName, protocol.StatusDeleted).First(&existingRule).Error; err == nil {
		log.Printf("Rule name already exists: %v", req.RuleName)
		return nil, protocol.InvalidParams // 使用现有错误码
	}

	// 检查促销码唯一性（如果提供了促销码）
	if req.PromoCode != "" {
		if err := models.GetDB().Where("promo_code = ? AND status != ?", req.PromoCode, protocol.StatusDeleted).First(&existingRule).Error; err == nil {
			log.Printf("Promo code already exists: %v", req.PromoCode)
			return nil, protocol.InvalidParams // 使用现有错误码
		}
	}

	// 验证时间范围合理性
	if req.StartDate != nil && req.EndDate != nil && *req.StartDate >= *req.EndDate {
		log.Printf("Invalid date range: start=%d, end=%d", *req.StartDate, *req.EndDate)
		return nil, protocol.InvalidParams
	}

	// 验证折扣百分比范围
	if req.DiscountPercent != nil && (*req.DiscountPercent < 0 || *req.DiscountPercent > 100) {
		log.Printf("Invalid discount percent: %f", *req.DiscountPercent)
		return nil, protocol.InvalidParams
	}

	// 验证折扣金额不能为负数
	if req.DiscountAmount != nil && *req.DiscountAmount < 0 {
		log.Printf("Invalid discount amount: %f", *req.DiscountAmount)
		return nil, protocol.InvalidParams
	}

	// 验证加价倍数范围
	if req.SurgeMultiplier != nil && (*req.SurgeMultiplier < 0.1 || *req.SurgeMultiplier > 10.0) {
		log.Printf("Invalid surge multiplier: %f", *req.SurgeMultiplier)
		return nil, protocol.InvalidParams
	}

	// 验证最低费用不能大于最高费用
	if req.MinimumFare != nil && req.MaximumFare != nil && *req.MinimumFare > *req.MaximumFare {
		log.Printf("Minimum fare cannot be greater than maximum fare: min=%f, max=%f", *req.MinimumFare, *req.MaximumFare)
		return nil, protocol.InvalidParams
	}

	// 2. 创建规则对象
	// 生成规则ID
	ruleID := utils.GeneratePriceRuleID()

	rule := &models.PriceRule{
		RuleID: ruleID,
		Salt:   utils.GenerateSalt(),
		PriceRuleValues: &models.PriceRuleValues{
			RuleName:        &req.RuleName,
			Description:     &req.Description,
			Category:        &req.Category,
			RuleType:        &req.RuleType,
			DiscountAmount:  req.DiscountAmount,
			DiscountPercent: req.DiscountPercent,
			SurgeMultiplier: req.SurgeMultiplier,
			BaseRate:        req.BaseRate,
			PerKmRate:       req.PerKmRate,
			PerMinuteRate:   req.PerMinuteRate,
			MinimumFare:     req.MinimumFare,
			MaximumFare:     req.MaximumFare,
			Priority:        &req.Priority,
			Status:          &req.Status,
			StartedAt:       req.StartDate,
			EndedAt:         req.EndDate,
			ServiceAreas:    req.ServiceAreas,
			UserCategories:  req.UserCategories,
			ApplicableRides: req.ApplicableRides,
			MaxUsagePerUser: req.MaxUsagePerUser,
			MaxUsageTotal:   req.MaxUsageTotal,
			PromoCode:       &req.PromoCode,
			RequiresCode:    &req.RequiresCode,
			CreatedBy:       &req.UserID,
			Metadata:        req.Metadata,
		},
		CreatedAt: time.Now().UnixMilli(),
	}
	for _, item := range req.VehicleFilters {
		if item.Category == "" && item.Level == "" {
			continue
		}
		rule.VehicleFilters = append(rule.VehicleFilters, &models.VehicleFilter{
			Category: item.Category,
			Level:    item.Level,
		})
	}

	// 设置默认值
	if rule.Status == nil || *rule.Status == "" {
		defaultStatus := protocol.StatusDraft
		rule.Status = &defaultStatus
	}

	if rule.Priority == nil {
		defaultPriority := 1
		rule.Priority = &defaultPriority
	}

	// 保存到数据库
	if err := models.GetDB().Create(rule).Error; err != nil {
		log.Printf("Failed to create price rule: %v", err)
		return nil, protocol.DatabaseError
	}

	return rule, protocol.Success
}

// GetPriceRuleByID 根据ID获取价格规则
func (s *PriceRuleService) GetPriceRuleByID(ruleID string) *models.PriceRule {
	if ruleID == "" {
		return nil
	}

	var rule models.PriceRule
	if err := models.GetDB().Where("rule_id = ?", ruleID).First(&rule).Error; err != nil {
		return nil
	}

	return &rule
}

// SearchPriceRule 获取价格规则列表
func (s *PriceRuleService) SearchPriceRule(page, pageSize int, category, status string) ([]models.PriceRule, int64, protocol.ErrorCode) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	query := models.GetDB().Model(&models.PriceRule{})

	// 添加过滤条件
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Printf("Failed to count price rules: %v", err)
		return nil, 0, protocol.DatabaseError
	}

	// 获取列表数据
	var rules []models.PriceRule
	if err := query.Offset(offset).Limit(pageSize).Order("priority ASC, created_at DESC").Find(&rules).Error; err != nil {
		log.Printf("Failed to get price rules: %v", err)
		return nil, 0, protocol.DatabaseError
	}

	return rules, total, protocol.Success
}

// UpdatePriceRule 更新价格规则
func (s *PriceRuleService) UpdatePriceRule(req *protocol.UpdatePriceRuleRequest) (*models.PriceRule, protocol.ErrorCode) {

	// 查找规则
	rule := s.GetPriceRuleByID(req.RuleID)
	if rule == nil {
		return nil, protocol.PriceIDNotFound
	}
	values := models.PriceRuleValues{}
	// 更新字段
	if req.RuleName != nil {
		values.RuleName = req.RuleName
	}
	if req.Description != nil {
		values.Description = req.Description
	}
	if req.Category != nil {
		values.Category = req.Category
	}
	if req.RuleType != nil {
		values.RuleType = req.RuleType
	}
	if req.DiscountAmount != nil {
		values.DiscountAmount = req.DiscountAmount
	}
	if req.DiscountPercent != nil {
		values.DiscountPercent = req.DiscountPercent
	}
	if req.SurgeMultiplier != nil {
		values.SurgeMultiplier = req.SurgeMultiplier
	}
	if req.MinimumFare != nil {
		values.MinimumFare = req.MinimumFare
	}
	if req.MaximumFare != nil {
		values.MaximumFare = req.MaximumFare
	}
	if req.Priority != nil {
		values.Priority = req.Priority
	}
	if req.Status != nil {
		values.Status = req.Status
	}
	if req.StartDate != nil {
		values.StartedAt = req.StartDate
	}
	if req.EndDate != nil {
		values.EndedAt = req.EndDate
	}
	if req.VehicleFilters != nil {
		values.VehicleFilters = convertProtocolToModelsVehicleFilters(req.VehicleFilters)
	}
	if req.ServiceAreas != nil {
		values.ServiceAreas = req.ServiceAreas
	}
	if req.UserCategories != nil {
		values.UserCategories = req.UserCategories
	}
	if req.ApplicableRides != nil {
		values.ApplicableRides = req.ApplicableRides
	}
	if req.MaxUsagePerUser != nil {
		values.MaxUsagePerUser = req.MaxUsagePerUser
	}
	if req.MaxUsageTotal != nil {
		values.MaxUsageTotal = req.MaxUsageTotal
	}
	if req.PromoCode != nil {
		values.PromoCode = req.PromoCode
	}
	if req.RequiresCode != nil {
		values.RequiresCode = req.RequiresCode
	}
	if req.Metadata != nil {
		values.Metadata = req.Metadata
	}

	// 设置更新时间
	values.UpdatedAt = time.Now().UnixMilli()

	// 保存更新 - 重要：必须指定WHERE条件避免更新所有记录
	if err := models.GetDB().Model(&models.PriceRule{}).
		Where("rule_id = ?", req.RuleID).
		UpdateColumns(values).Error; err != nil {
		log.Printf("Failed to update price rule: %v", err)
		return nil, protocol.DatabaseError
	}

	// 重新获取更新后的规则
	updatedRule := s.GetPriceRuleByID(req.RuleID)
	if updatedRule == nil {
		return nil, protocol.PriceIDNotFound
	}

	return updatedRule, protocol.Success
}

// DeletePriceRule 删除价格规则（软删除）
func (s *PriceRuleService) DeletePriceRule(ruleID string) protocol.ErrorCode {
	if ruleID == "" {
		return protocol.InvalidParams
	}

	// 检查规则是否存在
	rule := s.GetPriceRuleByID(ruleID)
	if rule == nil {
		return protocol.PriceIDNotFound
	}

	// 检查规则是否正在使用中（活跃状态的规则不能删除）
	if rule.GetStatus() == protocol.StatusActive {
		log.Printf("Cannot delete active price rule: %v", ruleID)
		return protocol.InvalidParams // 可以考虑添加更具体的错误码
	}

	// TODO: 这里可以添加更复杂的使用检查逻辑
	// 比如检查是否有正在进行的订单使用了这个规则
	// var orderCount int64
	// if err := models.GetDB().Model(&models.Order{}).Where("price_rule_id = ? AND status IN (?)", ruleID, []string{"pending", "confirmed", "in_progress"}).Count(&orderCount).Error; err == nil && orderCount > 0 {
	//     log.Printf("Price rule is being used by %d orders: %v", orderCount, ruleID)
	//     return protocol.InvalidParams
	// }

	// 执行软删除（将状态改为deleted）
	result := models.GetDB().Model(&models.PriceRule{}).
		Where("rule_id = ?", ruleID).
		Updates(map[string]interface{}{
			"status":     protocol.StatusDeleted,
			"updated_at": time.Now().UnixMilli(),
		})

	if result.Error != nil {
		log.Printf("Failed to delete price rule: %v", result.Error)
		return protocol.DatabaseError
	}

	if result.RowsAffected == 0 {
		return protocol.PriceIDNotFound
	}

	log.Printf("Price rule soft deleted: %v", ruleID)
	return protocol.Success
}

// UpdatePriceRuleStatus 更新价格规则状态
func (s *PriceRuleService) UpdatePriceRuleStatus(ruleID, status string) protocol.ErrorCode {
	if ruleID == "" {
		return protocol.InvalidParams
	}

	// 验证状态值
	validStatuses := []string{
		protocol.StatusDraft,
		protocol.StatusActive,
		protocol.StatusPaused,
		protocol.StatusExpired,
		protocol.StatusDeleted,
	}

	valid := slices.Contains(validStatuses, status)

	if !valid {
		log.Printf("Invalid status provided: %v", status)
		return protocol.InvalidParams
	}

	result := models.GetDB().Model(&models.PriceRule{}).
		Where("rule_id = ?", ruleID).
		Update("status", status)

	if result.Error != nil {
		log.Printf("Failed to update price rule status: %v", result.Error)
		return protocol.DatabaseError
	}

	if result.RowsAffected == 0 {
		return protocol.UserNotFound
	}

	return protocol.Success
}

type PriceContext struct {
	Request          *protocol.EstimateRequest
	Env              *EnvironmentContext
	Snapshot         *models.PriceSnapshot // 使用快照实体代替 OrderPrice
	Rule             *models.PriceRule
	RuleResults      map[string]*protocol.PriceRuleResult
	UserPromotionIDs []string  // 记录使用的用户优惠券ID列表
	StartTime        time.Time // 计算开始时间
}

// CreatePriceContext 创建价格计算上下文，包含初始快照
func CreatePriceContext(req *protocol.EstimateRequest) *PriceContext {
	nowtime := utils.TimeNowMilli()
	// 向后兼容性处理
	if req.VehicleCategory == "" {
		req.VehicleCategory = protocol.VehicleCategorySedan // 默认小车
	}
	if req.VehicleLevel == "" {
		req.VehicleLevel = protocol.VehicleLevelEconomy // 默认经济型
	}
	if req.Currency == "" {
		req.Currency = protocol.CurrencyRWF // 默认卢旺达法郎
	}
	if req.SnapshotDuration <= 0 {
		req.SnapshotDuration = int64(30 * time.Minute) // 默认快照有效期30分钟
	}
	if req.RequestedAt == 0 {
		req.RequestedAt = nowtime
	}
	// 时间相关
	if req.ScheduledAt == 0 {
		req.ScheduledAt = req.RequestedAt
	}

	// 构建环境上下文
	env := &EnvironmentContext{
		SurgeMultiplier: 1.5, // 设为1.5以启用涌潮定价
	}

	// 创建价格快照
	snapshot := models.NewPriceSnapshot(req.UserID)
	snapshot.SetExpiresAt(time.Now().Add(time.Duration(req.SnapshotDuration) * time.Minute).UnixMilli()).
		SetCurrency(req.Currency).
		SetDistance(req.EstimatedDistance).
		SetDuration(req.EstimatedDuration).
		SetOrderType(req.OrderType).
		SetVehicleCategory(req.VehicleCategory).
		SetVehicleLevel(req.VehicleLevel).
		SetScheduledAt(req.ScheduledAt)

	// 将详细信息存储到 metadata 中（不重复存储快照结构中已有专门字段的数据）
	metadata := protocol.MapData{}

	// 位置坐标信息
	metadata.Set("pickup_latitude", req.PickupLatitude)
	metadata.Set("pickup_longitude", req.PickupLongitude)
	metadata.Set("dropoff_latitude", req.DropoffLatitude)
	metadata.Set("dropoff_longitude", req.DropoffLongitude)
	metadata.Set("pickup_address", req.PickupAddress)
	metadata.Set("dropoff_address", req.DropoffAddress)
	metadata.Set("pickup_landmark", req.PickupLandmark)
	metadata.Set("dropoff_landmark", req.DropoffLandmark)
	metadata.Set("vehicle_category", req.VehicleCategory)
	metadata.Set("vehicle_level", req.VehicleLevel)

	// 基础价格（内部计算用）
	if req.BasePrice > 0 {
		metadata.Set("base_price", req.BasePrice)
	}

	metadata.Set("scheduled_at", req.ScheduledAt)

	// 上下文信息
	metadata.Set("passenger_count", req.PassengerCount)
	if req.SessionID != "" {
		metadata.Set("session_id", req.SessionID)
	}
	if req.VehicleCategory != "" {
		metadata.Set("vehicle_category", req.VehicleCategory)
	}
	if req.VehicleLevel != "" {
		metadata.Set("vehicle_level", req.VehicleLevel)
	}
	if req.ServiceArea != "" {
		metadata.Set("service_area", req.ServiceArea)
	}
	if req.PaymentMethod != "" {
		metadata.Set("payment_method", req.PaymentMethod)
	}
	if len(req.PromoCodes) > 0 {
		metadata.Set("promo_codes", req.PromoCodes)
		snapshot.SetPromoCodes(req.PromoCodes)
	}
	if req.UserCategory != "" {
		metadata.Set("user_category", req.UserCategory)
	}

	// 客户端信息
	if req.Platform != "" {
		metadata.Set("platform", req.Platform)
	}
	if req.AppVersion != "" {
		metadata.Set("app_version", req.AppVersion)
	}
	if req.UserAgent != "" {
		metadata.Set("user_agent", req.UserAgent)
	}
	if req.RequestIP != "" {
		metadata.Set("request_ip", req.RequestIP)
	}

	// 配置选项
	if req.SnapshotDuration > 0 {
		metadata.Set("snapshot_duration", req.SnapshotDuration)
	}
	snapshot.SetMetadata(metadata)

	// 创建并返回价格计算上下文
	return &PriceContext{
		Request:     req,
		Env:         env,
		Snapshot:    snapshot,
		RuleResults: map[string]*protocol.PriceRuleResult{},
		StartTime:   time.Now(),
	}
}

// EstimatePrice 估算价格并返回快照实体
func (s *PriceRuleService) EstimatePrice(req *protocol.EstimateRequest) *models.PriceSnapshot {
	// 设置默认币种
	if req.Currency == "" {
		req.Currency = protocol.CurrencyRWF
	}

	// 创建价格计算上下文（包含初始快照）
	ctx := CreatePriceContext(req)
	snapshot := ctx.Snapshot

	// 获取活跃的价格规则
	rules := models.GetActivePriceRules()
	baseRules := []*models.PriceRule{}
	otherRules := []*models.PriceRule{}
	for _, rule := range rules {
		if rule.GetCategory() == protocol.PriceRuleCategoryBasePricing {
			baseRules = append(baseRules, rule)
			continue
		}
		otherRules = append(otherRules, rule)
	}

	// 2. 优先处理用户优惠券
	userRules := []*models.PriceRule{}
	if len(req.PromoCodes) > 0 {
		userPromotions := models.GetAvailablePromotionsByUser(req.UserID, req.PromoCodes)
		for _, userPromo := range userPromotions {
			// 检查用户券是否可用且有效
			if userPromo.IsAvailable() {
				// 将用户优惠券适配为 PriceRule 实体
				adaptedRule := userPromo.ToPriceRule()
				if adaptedRule != nil {
					userRules = append(userRules, adaptedRule)
				}
			}
		}
	}

	// 1. 先计算基础定价规则
	for _, rule := range baseRules {
		ctx.Rule = rule
		s.CalculateRule(ctx)
	}
	// 3. 然后计算其他类型的系统规则（加价、折扣、促销等）
	for _, rule := range otherRules {
		category := rule.GetCategory()
		// 只处理支持的类型
		if !IsSupportedRuleCategory(category) {
			continue
		}
		ctx.Rule = rule
		s.CalculateRule(ctx)
	}

	for _, rule := range userRules {
		// 使用现有的价格规则计算逻辑
		ctx.Rule = rule
		s.CalculateRule(ctx)
	}
	// 完成快照的最终计算和设置
	FinalizeSnapshot(ctx)

	return snapshot
}

// FinalizeSnapshot 完成快照的最终计算和设置
func FinalizeSnapshot(ctx *PriceContext) {
	snapshot := ctx.Snapshot

	// 确保 PriceSnapshotValues 不为 nil
	if snapshot.PriceSnapshotValues == nil {
		snapshot.PriceSnapshotValues = &models.PriceSnapshotValues{}
	}

	// 计算原始价格 - 包含所有费用项，优惠前的价格
	originalFareBeforeDiscount := snapshot.GetBaseFare().Add(snapshot.GetSurgeFare()).Add(snapshot.GetDistanceFare()).Add(snapshot.GetTimeFare()).Add(snapshot.GetServiceFee())

	// 优惠后折扣价格 = 原始价格 - 折扣 - 促销优惠 - 用户优惠券折扣
	discountedFareAfterPromotions := originalFareBeforeDiscount.Add(snapshot.GetDiscountAmount()).
		Add(snapshot.GetPromoDiscount()).
		Add(snapshot.GetUserPromoDiscount())
	if discountedFareAfterPromotions.LessThan(decimal.Zero) {
		discountedFareAfterPromotions = decimal.Zero
	}

	// 将所有费用字段统一格式化为两位小数并设置到快照
	snapshot.SetBaseFare(snapshot.GetBaseFare().Round(2))
	snapshot.SetDistanceFare(snapshot.GetDistanceFare().Round(2))
	snapshot.SetTimeFare(snapshot.GetTimeFare().Round(2))
	snapshot.SetServiceFee(snapshot.GetServiceFee().Round(2))
	snapshot.SetSurgeFare(snapshot.GetSurgeFare().Round(2))
	snapshot.SetOriginalFare(originalFareBeforeDiscount.Round(2))      // 优惠前原始价格
	snapshot.SetDiscountedFare(discountedFareAfterPromotions.Round(2)) // 优惠后折扣价格

	// 将decimal金额处理为两位小数
	snapshot.SetDiscountAmount(snapshot.GetDiscountAmount().Round(2))
	snapshot.SetPromoDiscount(snapshot.GetPromoDiscount().Round(2))
	snapshot.SetUserPromoDiscount(snapshot.GetUserPromoDiscount().Round(2))

	// 设置统计信息
	snapshot.SetRulesEvaluated(len(ctx.RuleResults))
	snapshot.SetRulesApplied(len(ctx.RuleResults))

	// 设置计算元信息
	snapshot.SetCalculationTime(time.Since(ctx.StartTime).Milliseconds())
	snapshot.SetEngineVersion("v1.0.0")

	// 设置价格明细分解
	var breakdowns []*protocol.PriceRuleResult
	for _, ruleResult := range ctx.RuleResults {
		// 确保货币信息正确
		ruleResult.Ccy = snapshot.GetCurrency()
		// 转换 RuleResult 为 protocol.PriceRuleResult，并确保金额保留2位小数
		priceRuleResult := &protocol.PriceRuleResult{
			RuleID:      ruleResult.RuleID,
			RuleName:    ruleResult.RuleName,
			Category:    ruleResult.Category,
			DisplayName: ruleResult.DisplayName,
			Description: ruleResult.Description,

			Amount:  utils.RoundToTwoDecimal(ruleResult.Amount), // 四舍五入到2位小数
			Applied: ruleResult.Applied,
			Ccy:     ruleResult.Ccy,
			Reason:  ruleResult.Reason,
		}
		breakdowns = append(breakdowns, priceRuleResult)
	}
	snapshot.SetBreakdowns(breakdowns).
		SetUserPromotionIDs(ctx.UserPromotionIDs)

	// 记录使用的用户优惠券ID列表
}

// SavePriceSnapshot 保存价格快照到数据库
func (s *PriceRuleService) SavePriceSnapshot(snapshot *models.PriceSnapshot) error {
	if snapshot == nil {
		return errors.New("snapshot cannot be nil")
	}
	return models.GetDB().Create(snapshot).Error
}

// ValidateAndLockPriceID 验证并锁定价格ID
func (s *PriceRuleService) ValidateAndLockPriceID(priceID string) (*models.PriceSnapshot, protocol.ErrorCode) {
	if priceID == "" {
		return nil, protocol.MissingParams
	}

	snapshot := models.GetPriceSnapshotByID(priceID)
	if snapshot == nil {
		return nil, protocol.PriceIDNotFound
	}

	// 检查是否过期
	if snapshot.IsExpired() {
		return nil, protocol.PriceIDExpired
	}

	// 检查是否已被使用
	if snapshot.GetOrderID() != "" {
		return nil, protocol.PriceIDAlreadyUsed
	}

	return snapshot, protocol.Success
}
