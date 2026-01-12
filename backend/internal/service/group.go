package service

import "time"

// GroupModelRate represents per-model rate configuration within a group
type GroupModelRate struct {
	ID             int64
	GroupID        int64
	Model          string
	RateMultiplier float64
	CardPrice      *float64 // 次卡模式单次请求价格(USD)，nil表示不支持次卡
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Group struct {
	ID             int64
	Name           string
	Description    string
	Platform       string
	RateMultiplier float64
	IsExclusive    bool
	Status         string
	Hydrated       bool // indicates the group was loaded from a trusted repository source

	SubscriptionType    string
	DailyLimitUSD       *float64
	WeeklyLimitUSD      *float64
	MonthlyLimitUSD     *float64
	DefaultValidityDays int

	// 图片生成计费配置（antigravity 和 gemini 平台使用）
	ImagePrice1K *float64
	ImagePrice2K *float64
	ImagePrice4K *float64

	// Claude Code 客户端限制
	ClaudeCodeOnly  bool
	FallbackGroupID *int64

	// 计费模式
	BillingMode      string   // balance/subscription/card
	DefaultCardPrice *float64 // 次卡模式默认单次价格(USD)

	CreatedAt time.Time
	UpdatedAt time.Time

	AccountGroups []AccountGroup
	AccountCount  int64

	// 模型专属费率配置
	ModelRates []GroupModelRate
}

func (g *Group) IsActive() bool {
	return g.Status == StatusActive
}

func (g *Group) IsSubscriptionType() bool {
	return g.SubscriptionType == SubscriptionTypeSubscription
}

func (g *Group) IsFreeSubscription() bool {
	return g.IsSubscriptionType() && g.RateMultiplier == 0
}

// IsCardBillingMode 判断分组是否为次卡计费模式
func (g *Group) IsCardBillingMode() bool {
	return g.BillingMode == BillingModeCard
}

// IsSubscriptionBillingMode 判断分组是否为订阅计费模式
func (g *Group) IsSubscriptionBillingMode() bool {
	return g.BillingMode == BillingModeSubscription
}

// IsBalanceBillingMode 判断分组是否为余额计费模式
func (g *Group) IsBalanceBillingMode() bool {
	return g.BillingMode == "" || g.BillingMode == BillingModeBalance
}

func (g *Group) HasDailyLimit() bool {
	return g.DailyLimitUSD != nil && *g.DailyLimitUSD > 0
}

func (g *Group) HasWeeklyLimit() bool {
	return g.WeeklyLimitUSD != nil && *g.WeeklyLimitUSD > 0
}

func (g *Group) HasMonthlyLimit() bool {
	return g.MonthlyLimitUSD != nil && *g.MonthlyLimitUSD > 0
}

// GetImagePrice 根据 image_size 返回对应的图片生成价格
// 如果分组未配置价格，返回 nil（调用方应使用默认值）
func (g *Group) GetImagePrice(imageSize string) *float64 {
	switch imageSize {
	case "1K":
		return g.ImagePrice1K
	case "2K":
		return g.ImagePrice2K
	case "4K":
		return g.ImagePrice4K
	default:
		// 未知尺寸默认按 2K 计费
		return g.ImagePrice2K
	}
}

// GetCardPrice 获取指定模型的次卡价格
// 优先级：模型级价格 > 分组默认价格
// 如果都未配置，返回 nil
func (g *Group) GetCardPrice(model string) *float64 {
	if g == nil {
		return nil
	}
	// 1. 先查模型级别的次卡价格
	for _, mr := range g.ModelRates {
		if mr.Model == model && mr.CardPrice != nil {
			return mr.CardPrice
		}
	}
	// 2. 没有则使用分组默认次卡价格
	return g.DefaultCardPrice
}

// IsGroupContextValid reports whether a group from context has the fields required for routing decisions.
func IsGroupContextValid(group *Group) bool {
	if group == nil {
		return false
	}
	if group.ID <= 0 {
		return false
	}
	if !group.Hydrated {
		return false
	}
	if group.Platform == "" || group.Status == "" {
		return false
	}
	return true
}
